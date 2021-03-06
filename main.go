package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/WasinWatt/game-bot/mongo"
	"github.com/WasinWatt/game-bot/service"

	"github.com/WasinWatt/game-bot/api"
	"github.com/WasinWatt/game-bot/config"
	"github.com/jinzhu/configor"
	"github.com/line/line-bot-sdk-go/linebot"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	config := &config.Config{}
	configor.Load(config, "config.yml")

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccToken)
	must(err)

	dialInfo, err := mgo.ParseURL(config.MongoURI)
	must(err)

	session, err := mgo.DialWithInfo(dialInfo)
	must(err)
	defer session.Close()

	log.Println("Connected to DB")

	// Repo initialize
	repo := mongo.New()

	// Service controller initialize
	controller := service.New(repo)

	apiHandler := api.NewHandler(bot, session, controller)

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

	mux.Handle("/api/", http.StripPrefix("/api", apiHandler.MakeHandler()))

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
