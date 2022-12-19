package internal

import (
	"context"
	"fmt"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"os/exec"
)

type SimpleImpl struct {
	VClient *whatsmeow.Client
	Msg     *events.Message
}

func NewSimpleImpl(client *whatsmeow.Client, msg *events.Message) *SimpleImpl {
	return &SimpleImpl{
		VClient: client,
		Msg:     msg,
	}
}

func (simp *SimpleImpl) Reply(teks string) {
	simp.VClient.SendMessage(context.Background(), simp.Msg.Info.Chat, "", &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	})
}

//SendHydratedBtn is a function to send hydrated button
func (simp *SimpleImpl) SendHydratedBtn(jid types.JID, teks string, footer string, btns []*waProto.HydratedTemplateButton) {
	simp.VClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		TemplateMessage: &waProto.TemplateMessage{
			HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
				HydratedContentText: proto.String(teks),
				HydratedFooterText:  proto.String(footer),
				HydratedButtons:     btns,
			},
		},
	})
}

func (simp *SimpleImpl) SendContact(jid types.JID, contact string, name string) {
	simp.VClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(name),
			Vcard:       proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", name, name, contact, contact)),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	})
}

func (simp *SimpleImpl) FetchGroupAdmin(jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := simp.VClient.GetGroupInfo(jid)
	if err != nil {
		return Admin, err
	}
	for _, participant := range resp.Participants {
		if participant.IsAdmin || participant.IsSuperAdmin {
			Admin = append(Admin, participant.JID.String())
		}
	}
	return Admin, nil
}

func (simp *SimpleImpl) GetGroupAdmin(jid types.JID, sender string) bool {
	if !simp.Msg.Info.IsGroup {
		return false
	}
	admin, err := simp.FetchGroupAdmin(jid)
	if err != nil {
		return false
	}
	for _, v := range admin {
		if v == sender {
			return true
		}
	}
	return false
}

func (simp *SimpleImpl) CreateStickerIMG(data []byte) *waProto.Message {
	RawPath := fmt.Sprintf("tmp/%s%s", simp.Msg.Info.ID, ".jpg")
	ConvertedPath := fmt.Sprintf("tmp/sticker/%s%s", simp.Msg.Info.ID, ".webp")
	err := os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save image: %v", err)
	}
	exc := exec.Command("cwebp", "-q", "80", RawPath, "-o", ConvertedPath)
	err = exc.Run()
	if err != nil {
		fmt.Println("Failed to Convert Image to WebP")
	}
	createExif := fmt.Sprintf("webpmux -set exif %s %s -o %s", "tmp/exif/raw.exif", ConvertedPath, ConvertedPath)
	cmd := exec.Command("bash", "-c", createExif)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to set webp metadata", err)
	}
	stc, err := os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}
	uploaded, err := simp.VClient.Upload(context.Background(), stc, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(stc)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	}
}

func (simp *SimpleImpl) GetCMD() string {
	extended := simp.Msg.Message.GetExtendedTextMessage().GetText()
	text := simp.Msg.Message.GetConversation()
	imageMatch := simp.Msg.Message.GetImageMessage().GetCaption()
	videoMatch := simp.Msg.Message.GetVideoMessage().GetCaption()
	tempBtnId := simp.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
	var command string
	if text != "" {
		command = text
	} else if imageMatch != "" {
		command = imageMatch
	} else if videoMatch != "" {
		command = videoMatch
	} else if extended != "" {
		command = extended
	} else if tempBtnId != "" {
		command = tempBtnId
	}
	return command
}
