package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (app *Application) SetupRouter() {
	app.router = gin.Default()
	app.router.ForwardedByClientIP = app.BehindReverseProxy
	app.router.SetTrustedProxies([]string{app.TrustedProxy})

	app.router.StaticFile("/", app.staticPath+"/html/index.html")
	app.router.Static("/static", app.staticPath)
	app.router.POST("/api/confess", app.Confess)
	app.router.Any("/ws", func(c *gin.Context) {
		app.ws.HandleRequest(c.Writer, c.Request)
	})
}

// Curl command for testing
// curl 'http://localhost:3000/api/confess' --data-raw 'confession=I am confessing something'
type ConfessInput struct {
	Confession string `form:"confession" json:"confession"`

	// Wether confession should show up on the feed
	Public bool `form:"public" default:"false" json:"public"`
}

const MaxBodySize = 1000

func (app *Application) Confess(c *gin.Context) {
	var input ConfessInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rawConfession := strings.TrimSpace(input.Confession)
	if rawConfession == "" {
		c.String(http.StatusBadRequest, "confession cannot be empty")
		return
	} else if len(rawConfession) > MaxBodySize {
		c.String(http.StatusBadRequest, "confession cannot be longer than 1000 characters")
		return
	}

	confession := Confession{Confession: rawConfession, IpAddress: c.ClientIP(), Public: input.Public}

	if err := app.db.Create(&confession).Error; err != nil {
		log.Println("failed to add new confession:", err)
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if err := app.sendNtfyNotification(confession); err != nil {
		log.Println("ntfy error:", err)
	}

	// Sends notification to all sessions
	if input.Public {
		bs, err := json.Marshal(ConfessionOut{
			Confession: rawConfession,
			Date:       confession.CreatedAt,
		})
		if err != nil {
			log.Println("failed to marshal json:", err)
		} else {
			if err := app.ws.Broadcast(bs); err != nil {
				log.Println("failed to send notification to websockets:", err)
			}
		}
	}

	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}
