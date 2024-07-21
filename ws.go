package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/olahol/melody"
)

type WsEventType uint

const (
	InitialDataEventType WsEventType = iota
	ConfessionEventType
)

type WsEventCommons struct {
	Type WsEventType `json:"type"`
}

// Event used for sending the initial data on connect
type InitialDataEventInner struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
}
type InitialDataEvent struct {
	Confessions []InitialDataEventInner `json:"confessions"`

	WsEventCommons
}

func NewInitialDataEvent(Confessions []InitialDataEventInner) InitialDataEvent {
	return InitialDataEvent{
		Confessions: Confessions,
		WsEventCommons: WsEventCommons{
			Type: InitialDataEventType,
		},
	}
}

// Event sent when a new confession is encountered
type ConfessionEvent struct {
	Confession string    `json:"confession"`
	Date       time.Time `json:"date"`
	WsEventCommons
}

func NewConfessionEvent(Confession string, Date time.Time) ConfessionEvent {
	return ConfessionEvent{
		Confession: Confession,
		Date:       Date,
		WsEventCommons: WsEventCommons{
			Type: ConfessionEventType,
		},
	}
}

func (app *Application) SetupWebsocket() {
	app.ws = melody.New()
	app.ws.HandleConnect(app.HandleConnectWs)
	app.ws.HandleDisconnect(app.HandleDisconnectWs)
}

func (app *Application) HandleConnectWs(s *melody.Session) {
	// Fetch 5 recent confessions in the last 24 hours
	var confessions []Confession
	if err := app.db.Order("created_at desc").Where("public = true").Where("created_at > ?", time.Now().Add(-24*time.Hour)).Limit(5).Find(&confessions).Error; err != nil {
		log.Println("failed to fetch confessions:", err)
		return
	}

	var initialConfessions []InitialDataEventInner
	for _, confession := range confessions {
		initialConfessions = append(initialConfessions, InitialDataEventInner{Confession: confession.Confession, Date: confession.CreatedAt})
	}

	bs, err := json.Marshal(NewInitialDataEvent(initialConfessions))
	if err != nil {
		log.Println("failed to marshal json for initial event:", err)
		return
	}
	if err := s.Write(bs); err != nil {
		log.Println("failed to send initial info:", err)
	}
}

func (app *Application) HandleDisconnectWs(s *melody.Session) {
}
