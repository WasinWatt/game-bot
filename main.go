package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/configor"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/wasinwatt/game-bot/config"
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
	mux.Handle("/bot", gameBot)

	must(err)
	http.ListenAndServe(":3001", mux)
	log.Println("Listening on port: 3001")
}

func (gb *GameBot) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// events, err := gb.bot.ParseRequest(req)
	// must(err)
	// for _, event := range events {
	// 	if event.Type == linebot.EventTypeMessage {
	// 		// Do Something...
	// 	}
	// }
	response, _ := json.Marshal("Connected!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
