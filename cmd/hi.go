package cmd

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"vivcis/github.com/whatsapp-bot/helpers"
)

func HiHandler(evt interface{}, c *whatsmeow.Client) {
	switch v := evt.(type) {
	case *events.Message:
		msg := strings.ToLower(v.Message.GetConversation())

		switch msg {
		case "hi bot", "hello bot", "hey bot", "oi bot", "olá bot", "hi", "hello":
			err := helpers.SendMessage("Hello 👋!", c, v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
