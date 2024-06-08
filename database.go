package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database model
type Confession struct {
	gorm.Model

	Confession string
	IpAddress  string
	Public     bool // Wether confession should show up on the feed, true = shows up on feed!
}

func (app *Application) SetupDatabase() (err error) {
	app.db, err = gorm.Open(sqlite.Open(app.databasePath), &gorm.Config{})
	if err != nil {
		return
	}

	app.db.AutoMigrate(&Confession{})

	return
}
