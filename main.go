package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/WasinWatt/game-bot/room"
	"github.com/WasinWatt/game-bot/vocab"

	"github.com/WasinWatt/game-bot/api"
	"github.com/WasinWatt/game-bot/config"
	"github.com/WasinWatt/game-bot/user"
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
	userRepo := user.NewRepository()
	roomRepo := room.NewRepository()
	vocabRepo := vocab.NewRepository()

	// Service initialize
	userServ := user.NewService(userRepo, roomRepo)
	roomServ := room.NewService(roomRepo)
	vocabServ := vocab.NewService(vocabRepo)

	apiHandler := api.NewHandler(bot, session, userServ, roomServ, vocabServ)

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

func ensureIndex(session *mgo.Session) {
	c := session.DB("undercover").C("users")
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	must(err)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
