package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (app *Application) SetupRouter() {
	app.router = fiber.New(fiber.Config{
		EnableTrustedProxyCheck: app.BehindReverseProxy,
		TrustedProxies:          []string{app.TrustedProxy},
	})
	app.router.Static("/", app.staticPath)
	app.router.Post("/api/confess", app.Confess)

	app.router.Use("/ws", app.WsUpgrader)
	app.router.Get("/ws", websocket.New(app.ConfessionFeed))
}

// Curl command for testing
// curl 'http://localhost:3000/api/confess' --data-raw 'confession=I am confessing something'
type ConfessInput struct {
	Confession string `form:"confession" json:"confession"`

	// Wether confession should show up on the feed
	Public bool `form:"public" default:"false" json:"public"`
}

const MaxBodySize = 1000

func (app *Application) Confess(c *fiber.Ctx) error {
	input := new(ConfessInput)
	if err := c.BodyParser(input); err != nil {
		return err
	}

	rawConfession := strings.TrimSpace(input.Confession)
	if rawConfession == "" {
		c.Status(http.StatusBadRequest)
		return c.SendString("confession cannot be empty")
	} else if len(rawConfession) > MaxBodySize {
		c.Status(http.StatusBadRequest)
		return c.SendString("confession cannot be longer than 1000 characters")
	}

	confession := Confession{Confession: rawConfession, IpAddress: c.IP(), Public: input.Public}

	if err := app.db.Create(&confession).Error; err != nil {
		log.Println("failed to add new confession:", err)
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := app.sendNtfyNotification(confession); err != nil {
		log.Println("ntfy error:", err)
	}

	// Sends notification to all sessions
	if input.Public {
		app.CleanWsSessions()
		for _, s := range app.wsSessions {
			if !s.closed {
				s.notify <- ConfessionOut{
					Confession: rawConfession,
					Date:       confession.CreatedAt,
				}
			}
		}
	}

	return c.SendStatus(http.StatusOK)
}
