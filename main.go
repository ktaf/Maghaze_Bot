// Package main provides
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

// SalesKeyboard the main keyboard with buttons for listing status of things
var SalesKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("routers"),
		tgbotapi.NewKeyboardButton("Switches"),
		tgbotapi.NewKeyboardButton("Wires"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
	),
)

// url2File url, fname and returns int64
func url2File(url, fname string) int64 {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	outfile, _ := os.Create(fname)
	defer outfile.Close()

	n, _ := io.Copy(outfile, resp.Body)
	return n
}

// my check function
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// this is the main
func main() {

	// in1 := []byte("hello go\n")

	// err := ioutil.WriteFile("notes.txt", in1, 0644)
	// check(err)

	tgbot := os.Getenv("TGBOT")

	newbot, _ := tgbotapi.NewBotAPI(tgbot)

	// newbot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := newbot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = SalesKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		default:
			msg.Text = "Sorry unknown commmand :("
		}

		// handling Documents
		switch update.Message.Document {
		case nil:
			continue
		default:
			msgtime := update.Message.Time()
			url, _ := newbot.GetFileDirectURL(update.Message.Document.FileID)
			msg.Text = "got a doc at : " + msgtime.Format("Mon Jan 2 15:04:05 MST 2006") + "\n" + url
			url2File(url, update.Message.Document.FileName)
		}
		// telegram slash commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome to our shop!"
				msg.ReplyMarkup = SalesKeyboard
			case "login":
				password := update.Message.CommandArguments()
				fmt.Printf("argument = %+v\n", password)
				msg.Text = "your password is :" + password
			}
		}
		// msg.Text = "grrr XD"
		newbot.Send(msg)
	}

}
