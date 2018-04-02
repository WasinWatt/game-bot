package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/WasinWatt/game-bot/api"
	"github.com/WasinWatt/game-bot/config"
	"github.com/jinzhu/configor"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	config := &config.Config{}
	configor.Load(config, "config.yml")

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccToken)
	must(err)

	gameBot := &api.GameBot{
		Client: bot,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal("Status OK")
		w.WriteHeader(200)
		w.Write(response)
	})

	mux.Handle("/api/", http.StripPrefix("/api", api.MakeAPIHandler(gameBot)))

	must(err)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "3000"
	}
	http.ListenAndServe(":"+addr, mux)
	log.Println("Listening on port: " + addr)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
