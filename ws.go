package main

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// Removes closed sessions
func (app *Application) CleanWsSessions() {
	app.wsMutex.Lock()
	var cleanedSessions []*WsSession
	for _, s := range app.wsSessions {
		if !s.closed {
			cleanedSessions = append(cleanedSessions, s)
		}
	}
	app.wsSessions = cleanedSessions
	app.wsMutex.Unlock()
}

func (app *Application) WsUpgrader(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Confession data to show to user
type ConfessionOut struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
}

func (app *Application) ConfessionFeed(c *websocket.Conn) {
	// Fetch 5 recent confessions
	var confessions []Confession
	if err := app.db.Order("created_at desc").
		Where("public = true").Limit(5).Find(&confessions).Error; err != nil {
		log.Println("failed to fetch confessions:", err)
	} else {
		var out []ConfessionOut
		for _, confession := range confessions {
			out = append(out, ConfessionOut{
				Confession: confession.Confession,
				Date:       confession.CreatedAt,
			})
		}

		if err := c.WriteJSON(out); err != nil {
			log.Println("failed to send initial confessions:", err)
		}
	}

	session := WsSession{notify: make(chan ConfessionOut)}

	app.wsMutex.Lock()
	app.wsSessions = append(app.wsSessions, &session)
	app.wsMutex.Unlock()

	exit := make(chan bool)

	// exit worker
	go func() {
		for {
			// on exit sends websocket.ErrCloseSent
			if _, _, err := c.NextReader(); err != nil {
				exit <- true
				break
			}
		}
	}()

	for {
		select {
		case <-exit:
			session.close()
			return
		case cfn := <-session.notify:
			if err := c.WriteJSON(cfn); err != nil {
				log.Println("failed to send confession:", err)
			}
		}
	}
}
