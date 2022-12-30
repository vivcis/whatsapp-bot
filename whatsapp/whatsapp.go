package whatsapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vivcis/github.com/whatsapp-bot/handlers"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func Connect() error {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:sessions.db?_foreign_keys=on", dbLog)
	if err != nil {
		return err
	}
	//If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead
	device, err := container.GetFirstDevice()
	if err != nil {
		return err
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	c := whatsmeow.NewClient(device, clientLog)
	handlers.SetHandlers(c)

	if c.Store.ID == nil {
		qrChan, _ := c.GetQRChannel(context.Background())
		err = c.Connect()
		if err != nil {
			return err
		}
		for events := range qrChan {
			if events.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(events.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("Scan the QR code above to log in", events.Code)
			} else {
				fmt.Println("Login Success", events.Event)
			}
		}
	} else {
		//Already logged in, just connect
		err = c.Connect()
		fmt.Println("Login Success")
		if err != nil {
			return err
		}
	}

	client = c
	//Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	return err
}

func Disconnect() {
	client.Disconnect()
}

// func init() {
// 	ascii := figlet4go.NewAsciiRender()
// 	renderStr, _ := ascii.Render("cece bot")
// 	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
// 	store.DeviceProps.Os = proto.String("cece bot")
// 	fmt.Print(renderStr)
// }
