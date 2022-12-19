package message

import (
	"context"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"vivcis/github.com/whatsapp-bot/cmd"
	"vivcis/github.com/whatsapp-bot/internal"
)

var (
	prefix = "!"
	self   = true
	owner  = "0905-xxxx-xxxx"
)

func Msg(client *whatsmeow.Client, msg *events.Message) {
	simp := internal.NewSimpleImpl(client, msg)
	from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	args := strings.Split(simp.GetCMD(), " ")
	command := strings.ToLower(args[0])
	//query := strings.Join(args[1:], " ")
	//isAdmin := simp.GetGroupAdmin(from, sender)
	//isGroup := msg.Info.IsGroup
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	extended := msg.Message.GetExtendedTextMessage()
	quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()
	if self && !isOwner {
		return
	}
	switch command {
	case prefix + "menu":
		simp.Reply(cmd.Menu(pushName, prefix))
	case prefix + "owner":
		simp.SendContact(from, owner, "cece")
	case prefix + "source":
		simp.Reply("Source Code : https://github.com/vivcis/whatsapp-bot")
	case prefix + "sticker":
		if quotedImage != nil {
			data, _ := client.Download(quotedImage)
			stc := simp.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		} else if msg.Message.GetImageMessage() != nil {
			data, _ := client.Download(msg.Message.GetImageMessage())
			stc := simp.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		}
	}
}
