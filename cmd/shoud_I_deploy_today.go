package cmd

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"vivcis/github.com/whatsapp-bot/helpers"
)

type ShouldDeploy struct {
	TimeZone     string `json:"time_zone"`
	Day          string `json:"day"`
	Message      string `json:"message"`
	ShouldDeploy bool   `json:"should_deploy"`
}

func ShouldIDeployTodayHandler(evt interface{}, c *whatsmeow.Client) {
	switch v := evt.(type) {
	case *events.Message:
		msg := strings.ToLower(v.Message.GetConversation())

		if msg == "should i deploy" || msg == "should i deploy today?" || msg == "should i deploy today" {
			err := ShouldIDeployToday(c, v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func ShouldIDeployToday(client *whatsmeow.Client, receiver types.JID) error {
	sd := new(ShouldDeploy)
	helpers.GetJson("https://shouldideploy.today/api?tz=America/Sao_Paulo", sd)
	err := helpers.SendMessage(
		sd.Message,
		client,
		receiver)

	return err
}
