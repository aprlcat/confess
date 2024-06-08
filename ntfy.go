package main

import (
	"net/http"
	"strings"
)

func (app *Application) sendNtfyNotification(confession Confession) (err error) {
	if app.ntfyUrl == "" {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, app.ntfyUrl, strings.NewReader(confession.Confession))
	if err != nil {
		return
	}

	req.Header.Set("Title", "New confession from "+confession.IpAddress)
	_, err = http.DefaultClient.Do(req)
	return
}
