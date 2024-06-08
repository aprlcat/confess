package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WsSession struct {
	notify chan ConfessionOut
	closed bool
}

func (ws *WsSession) close() {
	ws.closed = true
	close(ws.notify)
}

type Config struct {
	port         uint   // port to run http server on
	staticPath   string // static files path
	databasePath string // sqlite database
	ntfyUrl      string // notification sending url
}

type Application struct {
	router *fiber.App
	db     *gorm.DB

	wsSessions []*WsSession
	wsMutex    sync.Mutex

	Config
}

func (app *Application) ParseConfig() {
	flag.UintVar(&app.port, "port", 3000, "port")
	flag.StringVar(&app.staticPath, "static", "static", "static files path")
	flag.StringVar(&app.databasePath, "database", "confession.db", "database path")
	flag.StringVar(&app.ntfyUrl, "ntfy", "", "ntfy url")
	flag.Parse()

	// Parse ntfy from ENV
	if app.ntfyUrl == "" {
		app.ntfyUrl = os.Getenv("NTFY_URL")
	}
}

func main() {
	var app Application
	app.ParseConfig()

	if err := app.SetupDatabase(); err != nil {
		log.Fatal("Database setup failed:", err)
	}

	app.SetupRouter()

	log.Fatal(app.router.Listen(":" + fmt.Sprint(app.port)))
}
