package cmd

import (
	"encoding/json"
	"log"
	"time"

	"github.com/olahol/melody"
)

type wsEventType uint

const (
	initialDataEventType wsEventType = iota
	confessionEventType
)

type wsEventCommons struct {
	Type wsEventType `json:"type"`
}

// Event used for sending the initial data on connect
type initialDataEventInner struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
}
type initialDataEvent struct {
	Confessions []initialDataEventInner `json:"confessions"`

	wsEventCommons
}

func newInitialDataEvent(Confessions []initialDataEventInner) initialDataEvent {
	return initialDataEvent{
		Confessions: Confessions,
		wsEventCommons: wsEventCommons{
			Type: initialDataEventType,
		},
	}
}

// Event sent when a new confession is encountered
type confessionEvent struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
	wsEventCommons
}

func newConfessionEvent(Confession string, Date time.Time) confessionEvent {
	return confessionEvent{
		Confession: Confession,
		Date:       Date,
		wsEventCommons: wsEventCommons{
			Type: confessionEventType,
		},
	}
}

func (app *Application) setupWebsocket() {
	app.ws = melody.New()
	app.ws.HandleConnect(app.handleConnectWs)
	app.ws.HandleDisconnect(app.handleDisconnectWs)
}

func (app *Application) handleConnectWs(s *melody.Session) {
	// Fetch 5 recent confessions in the last 24 hours
	var confessions []confession
	if err := app.db.Order("created_at desc").Where("public = true").Where("created_at > ?", time.Now().Add(-24*time.Hour)).Limit(5).Find(&confessions).Error; err != nil {
		log.Println("failed to fetch confessions:", err)
		return
	}

	var initialConfessions []initialDataEventInner
	for _, confession := range confessions {
		initialConfessions = append(initialConfessions, initialDataEventInner{Confession: confession.Confession, Date: confession.CreatedAt})
	}

	bs, err := json.Marshal(newInitialDataEvent(initialConfessions))
	if err != nil {
		log.Println("failed to marshal json for initial event:", err)
		return
	}
	if err := s.Write(bs); err != nil {
		log.Println("failed to send initial info:", err)
	}
}

func (app *Application) handleDisconnectWs(s *melody.Session) {
}
