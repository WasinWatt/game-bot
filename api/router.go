package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

// MakeAPIHandler make default handler
func MakeAPIHandler(gb *GameBot) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", manageLineRequest(gb))
	mux.HandleFunc("/health", checkRequest())
	return mux
}

func manageLineRequest(gb *GameBot) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		events, err := gb.Client.ParseRequest(req)
		if err != nil {
			log.Fatal(err)
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				userID := event.Source.UserID
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					handleTextMessage(gb, event, message, userID)
				}
			}
		}
	}
}

func checkRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal("/api is ok")
		w.WriteHeader(200)
		w.Write(response)
	}
}

func handleTextMessage(gb *GameBot, event linebot.Event, message *linebot.TextMessage, userID string) {
	_, err := gb.Client.ReplyMessage(event.ReplyToken,
		linebot.NewTextMessage(message.ID+":"+message.Text+" OK!"+" By: "+userID)).Do()
	if err != nil {
		log.Fatal(err)
	}
}
