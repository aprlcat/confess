package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/olahol/melody"
)

// Confession data to show to user
type ConfessionOut struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
}

func (app *Application) SetupWebsocket() {
	app.ws = melody.New()
	app.ws.HandleConnect(app.HandleConnectWs)
}

func (app *Application) HandleConnectWs(s *melody.Session) {
	// Fetch 5 recent confessions in the last 24 hours
	var confessions []Confession
	if err := app.db.Order("created_at desc").Where("public = true").Where("created_at > ?", time.Now().Add(-24*time.Hour)).Limit(5).Find(&confessions).Error; err != nil {
		log.Println("failed to fetch confessions:", err)
		return
	}

	var out []ConfessionOut
	for _, confession := range confessions {
		out = append(out, ConfessionOut{
			Confession: confession.Confession,
			Date:       confession.CreatedAt,
		})
	}

	bs, err := json.Marshal(out)
	if err != nil {
		log.Println("failed to marshal json:", err)
	}

	if err := s.Write(bs); err != nil {
		log.Println("failed to send initial confessions:", err)
	}
}
