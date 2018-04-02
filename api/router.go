package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/WasinWatt/game-bot/room"
	"github.com/WasinWatt/game-bot/user"
	"github.com/line/line-bot-sdk-go/linebot"
)

// MakeAPIHandler make default handler
func MakeAPIHandler(gb *GameBot) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", checkRequest())
	mux.Handle("/line", manageLineRequest(gb))
	return mux
}

func manageLineRequest(gb *GameBot) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		events, err := gb.Client.ParseRequest(req)
		if err != nil {
			log.Println(err)
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
	var words []string
	words = strings.Split(message.Text, " ")
	switch words[0] {
	case "create":
		if len(words) < 2 {
			replyDefaultMessage(gb, event)
			return
		}
		err := room.Create(gb.Session, userID, words[1])
		if err != nil {
			log.Fatal(err)
		}

		reply := "Room: " + words[1] + " creation successful!"
		replyMessage(gb, event, reply)

	case "join":
		if len(words) < 3 {
			replyDefaultMessage(gb, event)
		} else {
			u := &user.User{
				ID:     userID,
				RoomID: words[1],
				Name:   words[2],
			}
			err := user.JoinRoom(gb.Session, u)
			if err != nil {
				log.Fatal(err)
			}

			reply := "Join room: " + words[1] + " successful!"
			replyMessage(gb, event, reply)
		}
	default:
		replyDefaultMessage(gb, event)
	}

}

func replyDefaultMessage(gb *GameBot, event linebot.Event) {
	message := `ทำตามคำสั่งด้านล่างเท่านั้นนะจ๊ะ
- create {เลขห้อง} : สร้างห้องเพื่อเล่นเกม
- join {เลขห้อง} {ชื่อที่ใช้เล่นเกม} : เข้าห้องเพื่อรอเล่นเกม
- leave : ออกจากห้องเกมปัจจุบัน
- help : แสดงข้อความคำสั่งทั้งหมด`

	replyMessage(gb, event, message)
}

func replyMessage(gb *GameBot, event linebot.Event, message string) {
	_, err := gb.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.Fatal(err)
	}
}
