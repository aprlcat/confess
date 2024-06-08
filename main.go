package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gorm.io/gorm"
)

type Config struct {
	port         uint   // port to run http server on
	staticPath   string // static files path
	databasePath string // sqlite database
	ntfyUrl      string // notification sending url

	BehindReverseProxy bool   // If running behind reverse proxy, enable this
	TrustedProxy       string // Trusted proxy for reverse proxy
}

type Application struct {
	router *gin.Engine
	db     *gorm.DB

	ws *melody.Melody

	Config
}

func (app *Application) ParseConfig() {
	flag.UintVar(&app.port, "port", 3000, "port")
	flag.StringVar(&app.staticPath, "static", "static", "static files path")
	flag.StringVar(&app.databasePath, "database", "confession.db", "database path")
	flag.StringVar(&app.ntfyUrl, "ntfy", "", "ntfy url")
	flag.BoolVar(&app.BehindReverseProxy, "reverse-proxy", false, "behind reverse proxy")
	flag.StringVar(&app.TrustedProxy, "trusted-proxy", "", "trusted proxy for reverse proxy")
	flag.Parse()

	// Parse ntfy from ENV
	if app.ntfyUrl == "" {
		app.ntfyUrl = os.Getenv("NTFY_URL")
	}

	if !app.BehindReverseProxy {
		app.BehindReverseProxy = os.Getenv("BEHIND_REVERSE_PROXY") == "true"
	}

	if app.TrustedProxy == "" {
		app.TrustedProxy = os.Getenv("TRUSTED_PROXY")
	}
}

func main() {
	var app Application
	app.ParseConfig()

	app.SetupWebsocket()

	if err := app.SetupDatabase(); err != nil {
		log.Fatal("Database setup failed:", err)
	}

	app.SetupRouter()

	log.Fatal(app.router.Run(":" + fmt.Sprint(app.port)))
}
