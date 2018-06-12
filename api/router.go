package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/WasinWatt/game-bot/service"
	"github.com/WasinWatt/game-bot/user"
	"github.com/line/line-bot-sdk-go/linebot"
	mgo "gopkg.in/mgo.v2"
)

// Handler is a api handler
type Handler struct {
	Client     *linebot.Client
	Session    *mgo.Session
	controller *service.Controller
}

// NewHandler creates new hanlder
func NewHandler(lineClient *linebot.Client, session *mgo.Session, controller *service.Controller) *Handler {
	return &Handler{
		Client:     lineClient,
		Session:    session,
		controller: controller,
	}
}

// MakeHandler make default handler
func (h *Handler) MakeHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", checkRequest())
	mux.Handle("/line", h.manageLineRequest())
	return mux
}

func (h *Handler) manageLineRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		events, err := h.Client.ParseRequest(req)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				userID := event.Source.UserID
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					err := h.handleTextMessage(message, userID)
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

func (h *Handler) handleTextMessage(message *linebot.TextMessage, userID string) error {
	var words []string
	words = strings.Split(message.Text, " ")
	command := strings.ToLower(words[0])
	if command == "create" {
		if len(words) < 3 {
			replyDefaultMessage(h.Client, userID)
			return nil
		}
		err := h.controller.CreateRoom(h.Session, words[1], userID)
		if err != nil {
			return err
		}

		u := &user.User{
			ID:     userID,
			RoomID: words[1],
			Name:   words[2],
		}
		err = h.controller.Join(h.Session, u)
		if err != nil {
			return err
		}

		reply := "Room: " + words[1] + " creation successful!"

		replyMessage(h.Client, userID, reply)
		return nil
	}

	if command == "join" {
		if len(words) < 3 {
			replyDefaultMessage(h.Client, userID)
		} else {
			u := &user.User{
				ID:     userID,
				RoomID: words[1],
				Name:   words[2],
			}
			err := h.controller.Join(h.Session, u)
			if err != nil {
				replyInternalErrorMessage(h.Client, userID)
				return err
			}

			reply := "Join room: " + words[1] + " successful!"
			replyMessage(h.Client, userID, reply)
		}
		return nil
	}

	if command == "leave" || command == "quit" {
		isOwner, players, err := h.controller.Leave(h.Session, userID)
		if err == service.ErrNotFound {
			replyMessage(h.Client, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
			return err
		}

		reply := "มึงออกจากห้องแล้วจ้าาาา"

		if isOwner {
			if err != nil {
				replyInternalErrorMessage(h.Client, userID)
				return err
			}

			for i := range players {
				go func(id string) {
					replyMessage(h.Client, id, reply)
				}(players[i].ID)
			}
		} else {
			replyMessage(h.Client, userID, reply)
		}
		return nil

	}
	if command == "list" {
		x, err := h.controller.GetUser(h.Session, userID)
		if err == service.ErrNotFound {
			replyMessage(h.Client, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
			return err
		}

		players, err := h.controller.GetAllUsersByRoomID(h.Session, x.RoomID)
		if err == service.ErrNotFound {
			replyMessage(h.Client, userID, "Nobody is in the room.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
			return err
		}

		reply := ""
		for i := range players {
			reply = reply + players[i].Name + "\n"
			log.Println(players[i].Name)
		}
		replyMessage(h.Client, userID, reply)
		return nil
	}

	if command == "start" || command == "begin" {
		players, err := h.controller.GetAllUsersByRoomID(h.Session, userID)

		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
			return err
		}

		if len(players) < 5 {
			for i := range players {
				replyMessage(h.Client, players[i].ID, "Need at least 5 players to begin the game")
			}
			return nil
		}

		v, err := h.controller.GetVocab(h.Session)

		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
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

			err := h.controller.AddUserRole(h.Session, players[i].ID, shuffledList[i])
			if err != nil {
				return err
			}

			replyMessage(h.Client, players[i].ID, userWord)
		}
		return nil
	}

	if command == "add" || command == "vocab" {
		if len(words) < 3 {
			replyDefaultMessage(h.Client, userID)
			return nil
		}

		err := h.controller.AddVocab(h.Session, words[1], words[2])
		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
		}

		reply := "Add vocab successful!"
		replyMessage(h.Client, userID, reply)
		return nil

	}

	if command == "mockstart" || command == "mockbegin" {
		players, err := h.controller.GetAllUsersByRoomID(h.Session, userID)

		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
			return err
		}

		v, err := h.controller.GetVocab(h.Session)

		if err != nil {
			replyInternalErrorMessage(h.Client, userID)
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
			err := h.controller.AddUserRole(h.Session, players[i].ID, shuffledList[i])
			if err != nil {
				return err
			}

			replyMessage(h.Client, players[i].ID, userWord)
		}
		return nil
	}

	replyDefaultMessage(h.Client, userID)
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
	replySticker(client, userID, "2", "520")
	replyMessage(client, userID, message)
}

func replySticker(client *linebot.Client, userID string, packageID string, stickerID string) {
	client.PushMessage(userID, linebot.NewStickerMessage(packageID, stickerID)).Do()
}

func replyMessage(client *linebot.Client, userID string, message string) {
	client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
}
