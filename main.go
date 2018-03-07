package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/WasinWatt/game-bot/config"
	"github.com/jinzhu/configor"
	"github.com/line/line-bot-sdk-go/linebot"
)

// GameBot is the app's bot
type GameBot struct {
	bot *linebot.Client
}

func main() {
	config := &config.Config{}
	configor.Load(config, "config.yml")

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccToken)

	gameBot := &GameBot{bot}
	must(err)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal("Status OK")
		w.WriteHeader(200)
		w.Write(response)
	})
	mux.Handle("/bot", gameBot)

	must(err)

	// addr := 3000
	http.ListenAndServe(":3000", mux)
	log.Println("Listening on port: 3000")
}

func (gb *GameBot) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	events, err := gb.bot.ParseRequest(req)
	must(err)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			userID := event.Source.UserID
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = gb.bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage(message.ID+":"+message.Text+" OK!"+" By: "+userID)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	// response, _ := json.Marshal("Connected!")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
	// w.Write(response)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
