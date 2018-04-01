package api

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

// MakeAPIHandler make default handler
func MakeAPIHandler(gb *GameBot) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", manageLineRequest(gb))
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

func handleTextMessage(gb *GameBot, event linebot.Event, message *linebot.TextMessage, userID string) {
	_, err := gb.Client.ReplyMessage(event.ReplyToken,
		linebot.NewTextMessage(message.ID+":"+message.Text+" OK!"+" By: "+userID)).Do()
	if err != nil {
		log.Fatal(err)
	}
}
