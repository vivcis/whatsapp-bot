package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/probandula/figlet4go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	
	"vivcis/github.com/whatsapp-bot/message"
)

func main() {
	dbLog := waLog.Stdout("db", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:sessions.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	//If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead
	device, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Clent", "INFO", true)
	client := whatsmeow.NewClient(device, clientLog)
	eventHandler := registerHandler(client)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for events := range qrChan {
			if events.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(events.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("Scan the QR code above to log in")
			} else {
				fmt.Println("Login Success")
			}
		}
	} else {
		//Already logged in, just connect
		err = client.Connect()
		fmt.Println("Login Success")
		if err != nil {
			panic(err)
		}
	}
	//Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func init() {
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render("cece bot")
	// Set Browser
	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	store.DeviceProps.Os = proto.String("cece bot")
	// Print Banner
	fmt.Print(renderStr)
}

func registerHandler(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			go message.Msg(client, v)
			break
		}
	}
}
