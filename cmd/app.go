package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gorm.io/gorm"
)

type config struct {
	port         uint   // port to run http server on
	staticPath   string // static files path
	databasePath string // sqlite database
	ntfyUrl      string // notification sending url

	behindReverseProxy bool   // If running behind reverse proxy, enable this
	trustedProxy       string // Trusted proxy for reverse proxy
}

type Application struct {
	router *gin.Engine
	db     *gorm.DB

	ws *melody.Melody

	config
}

func (app *Application) parseConfig() {
	flag.UintVar(&app.port, "port", 3000, "port")
	flag.StringVar(&app.staticPath, "static", "static", "static files path")
	flag.StringVar(&app.databasePath, "database", "confession.db", "database path")
	flag.StringVar(&app.ntfyUrl, "ntfy", "", "ntfy url")
	flag.BoolVar(&app.behindReverseProxy, "reverse-proxy", false, "behind reverse proxy")
	flag.StringVar(&app.trustedProxy, "trusted-proxy", "", "trusted proxy for reverse proxy")
	flag.Parse()

	// Parse ntfy from ENV
	if app.ntfyUrl == "" {
		app.ntfyUrl = os.Getenv("NTFY_URL")
	}

	if !app.behindReverseProxy {
		app.behindReverseProxy = os.Getenv("BEHIND_REVERSE_PROXY") == "true"
	}

	if app.trustedProxy == "" {
		app.trustedProxy = os.Getenv("TRUSTED_PROXY")
	}
}

func NewApplication() (app Application) {
	app.parseConfig()

	app.setupWebsocket()

	if err := app.setupDatabase(); err != nil {
		log.Fatal("Database setup failed:", err)
	}

	app.setupRouter()

	return
}

func (app *Application) Run() {
	log.Fatal(app.router.Run(":" + fmt.Sprint(app.port)))
}
