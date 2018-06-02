package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/WasinWatt/game-bot/room"
	"github.com/WasinWatt/game-bot/user"
	"github.com/WasinWatt/game-bot/vocab"
	"github.com/line/line-bot-sdk-go/linebot"
	mgo "gopkg.in/mgo.v2"
)

// Handler is a api handler
type Handler struct {
	Client  *linebot.Client
	Session *mgo.Session
}

// NewHandler creates new hanlder
func NewHandler(lineClient *linebot.Client, session *mgo.Session) *Handler {
	return &Handler{
		Client:  lineClient,
		Session: session,
	}
}

// MakeAPIHandler make default handler
func (h *Handler) MakeHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", checkRequest())
	mux.Handle("/line", manageLineRequest(h.Client, h.Session))
	return mux
}

func manageLineRequest(client *linebot.Client, session *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		events, err := client.ParseRequest(req)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				userID := event.Source.UserID
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					err := handleTextMessage(client, session, message, userID)
					if err != nil {
						w.Header().Set("Content-type", "application/json; charset=utf-8")
						w.WriteHeader(http.StatusInternalServerError)
					}
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

func handleTextMessage(client *linebot.Client, session *mgo.Session, message *linebot.TextMessage, userID string) error {
	var words []string
	words = strings.Split(message.Text, " ")
	command := strings.ToLower(words[0])
	if command == "create" {
		if len(words) < 3 {
			replyDefaultMessage(client, userID)
			return nil
		}
		err := room.Create(session, words[1], userID)
		if err != nil {
			return err
		}

		u := &user.User{
			ID:     userID,
			RoomID: words[1],
			Name:   words[2],
		}
		err = user.JoinRoom(session, u)
		if err != nil {
			return err
		}

		reply := "Room: " + words[1] + " creation successful!"
		replyMessage(client, userID, reply)

	} else if command == "join" {
		if len(words) < 3 {
			replyDefaultMessage(client, userID)
		} else {
			u := &user.User{
				ID:     userID,
				RoomID: words[1],
				Name:   words[2],
			}
			err := user.JoinRoom(session, u)
			if err != nil {
				replyInternalErrorMessage(client, userID)
				return err
			}

			reply := "Join room: " + words[1] + " successful!"
			replyMessage(client, userID, reply)
		}

	} else if command == "leave" || command == "quit" {
		isOwner, players, err := user.Leave(session, userID)
		if err == user.ErrNotFound {
			replyMessage(client, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		reply := "มึงออกจากห้องแล้วจ้าาาา"

		if isOwner {
			if err != nil {
				replyInternalErrorMessage(client, userID)
				return err
			}

			for i := range players {
				go func(id string) {
					replyMessage(client, id, reply)
				}(players[i].ID)
			}
		} else {
			replyMessage(client, userID, reply)
		}

	} else if command == "list" {
		x, err := user.Get(session, userID)
		if err == user.ErrNotFound {
			replyMessage(client, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		players, err := user.GetAllByRoomID(session, x.RoomID)
		if err == user.ErrNotFound {
			replyMessage(client, userID, "Nobody is in the room.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		reply := ""
		for i := range players {
			reply = reply + players[i].Name + "\n"
			log.Println(players[i].Name)
		}
		replyMessage(client, userID, reply)

	} else if command == "start" || command == "begin" {
		players, err := user.GetAllByRoomID(session, userID)

		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		if len(players) < 5 {
			for i := range players {
				go func(id string) {
					replyMessage(client, id, "Need at least 5 players to begin the game")
				}(players[i].ID)
			}
			return nil
		}

		v, err := vocab.Get(session)

		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		var normalWord string
		var undercoverWord string
		if rand.Intn(2) == 0 {
			normalWord = v.First
			undercoverWord = v.Second
		} else {
			normalWord = v.Second
			undercoverWord = v.First
		}

		roleNumList := make([]int, 10)
		roleNumList[0] = 0
		roleNumList[0] = 1
		roleNumList[0] = 2
		roleNumList[0] = 2
		roleNumList[0] = 2

		if len(players) == 6 {
			roleNumList = append(roleNumList, 1)
		}

		if len(players) == 7 {
			roleNumList = append(roleNumList, 0)
		}

		shuffledList := make([]int, len(players))
		perm := rand.Perm(len(players))
		for i, v := range perm {
			shuffledList[v] = roleNumList[i]
		}

		for i := range shuffledList {
			var userWord string
			switch shuffledList[i] {
			case 0:
				userWord = ""
			case 1:
				userWord = undercoverWord
			case 2:
				userWord = normalWord
			}

			err := user.AddRole(session, players[i].ID, shuffledList[i])
			if err != nil {
				return err
			}

			replyMessage(client, players[i].ID, userWord)
		}

	} else if command == "add" || command == "vocab" {
		if len(words) < 3 {
			replyDefaultMessage(client, userID)
			return nil
		}

		err := vocab.Add(session, words[1], words[2])
		if err != nil {
			replyInternalErrorMessage(client, userID)
		}

		reply := "Add vocab successful!"
		replyMessage(client, userID, reply)

	} else if command == "mockstart" || command == "mockbegin" {
		players, err := user.GetAllByRoomID(session, userID)

		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		v, err := vocab.Get(session)

		if err != nil {
			replyInternalErrorMessage(client, userID)
			return err
		}

		var normalWord string
		var undercoverWord string
		if rand.Intn(2) == 0 {
			normalWord = v.First
			undercoverWord = v.Second
		} else {
			normalWord = v.Second
			undercoverWord = v.First
		}

		roleNumList := make([]int, 10)
		roleNumList[0] = 1
		roleNumList[0] = 2

		shuffledList := make([]int, len(players))
		perm := rand.Perm(len(players))
		for i, v := range perm {
			shuffledList[v] = roleNumList[i]
		}

		for i := range shuffledList {
			var userWord string
			switch shuffledList[i] {
			case 0:
				userWord = ""
			case 1:
				userWord = undercoverWord
			case 2:
				userWord = normalWord
			}
			err := user.AddRole(session, players[i].ID, shuffledList[i])
			if err != nil {
				return err
			}

			replyMessage(client, players[i].ID, userWord)
		}

	} else {
		replyDefaultMessage(client, userID)

	}

	return nil
}

func replyInternalErrorMessage(client *linebot.Client, userID string) {
	message := `ระบบขัดข้อง กรุณาลองใหม่`
	replyMessage(client, userID, message)
}

func replyDefaultMessage(client *linebot.Client, userID string) {
	message := `ทำตามคำสั่งด้านล่างเท่านั้นนะจ๊ะ
- create {เลขห้อง} : สร้างห้องเพื่อเล่นเกม
- join {เลขห้อง} {ชื่อที่ใช้เล่นเกม} : เข้าห้องเพื่อรอเล่นเกม
- leave : ออกจากห้องเกมปัจจุบัน
- help : แสดงข้อความคำสั่งทั้งหมด`

	replyMessage(client, userID, message)
}

func replyMessage(client *linebot.Client, userID string, message string) {
	client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
}
