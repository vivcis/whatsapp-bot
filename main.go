package main

import (
	_ "github.com/mattn/go-sqlite3"
	"vivcis/github.com/whatsapp-bot/whatsapp"
)

func main() {
	err := whatsapp.Connect()
	if err != nil {
		panic(err)
	}
	defer whatsapp.Disconnect()
}
