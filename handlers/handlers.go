package handlers

import (
	"go.mau.fi/whatsmeow"
	"vivcis/github.com/whatsapp-bot/cmd"
)

func SetHandlers(c *whatsmeow.Client) {
	c.AddEventHandler(func(evt interface{}) {
		cmd.HiHandler(evt, c)
	})

	c.AddEventHandler(func(evt interface{}) {
		cmd.ShouldIDeployTodayHandler(evt, c)
	})
}
