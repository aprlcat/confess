package cmd

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database model
type confession struct {
	gorm.Model

	Confession string
	IpAddress  string
	Public     bool // Wether confession should show up on the feed, true = shows up on feed!
}

func (app *Application) setupDatabase() (err error) {
	app.db, err = gorm.Open(sqlite.Open(app.databasePath), &gorm.Config{})
	if err != nil {
		return
	}

	app.db.AutoMigrate(&confession{})

	return
}
