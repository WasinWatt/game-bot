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
					err := handleTextMessage(gb, message, userID)
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

func handleTextMessage(gb *GameBot, message *linebot.TextMessage, userID string) error {
	var words []string
	words = strings.Split(message.Text, " ")
	command := strings.ToLower(words[0])
	if command == "create" {
		if len(words) < 3 {
			replyDefaultMessage(gb, userID)
			return nil
		}
		err := room.Create(gb.Session, words[1], userID)
		if err != nil {
			return err
		}

		u := &user.User{
			ID:     userID,
			RoomID: words[1],
			Name:   words[2],
		}
		err = user.JoinRoom(gb.Session, u)
		if err != nil {
			return err
		}

		reply := "Room: " + words[1] + " creation successful!"
		replyMessage(gb, userID, reply)

	} else if command == "join" {
		if len(words) < 3 {
			replyDefaultMessage(gb, userID)
		} else {
			u := &user.User{
				ID:     userID,
				RoomID: words[1],
				Name:   words[2],
			}
			err := user.JoinRoom(gb.Session, u)
			if err != nil {
				replyInternalErrorMessage(gb, userID)
				return err
			}

			reply := "Join room: " + words[1] + " successful!"
			replyMessage(gb, userID, reply)
		}

	} else if command == "leave" || command == "quit" {
		isOwner, players, err := user.Leave(gb.Session, userID)
		if err == user.ErrNotFound {
			replyMessage(gb, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(gb, userID)
			return err
		}

		reply := "มึงออกจากห้องแล้วจ้าาาา"

		if isOwner {
			if err != nil {
				replyInternalErrorMessage(gb, userID)
				return err
			}

			for i := range players {
				go func(id string) {
					replyMessage(gb, id, reply)
				}(players[i].ID)
			}
		} else {
			replyMessage(gb, userID, reply)
		}

	} else if command == "list" {
		x, err := user.Get(gb.Session, userID)
		if err == user.ErrNotFound {
			replyMessage(gb, userID, "You are not in any room. Please join first.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(gb, userID)
			return err
		}

		players, err := user.GetAllByRoomID(gb.Session, x.RoomID)
		if err == user.ErrNotFound {
			replyMessage(gb, userID, "Nobody is in the room.")
			return nil
		}
		if err != nil {
			replyInternalErrorMessage(gb, userID)
			return err
		}

		reply := ""
		for i := range players {
			reply = reply + players[i].Name + "\n"
			log.Println(players[i].Name)
		}
		replyMessage(gb, userID, reply)

		// } else if command == "start" || command == "begin" {
		// 	players, err := user.GetAllByRoomID(gb.Session, userID)

		// 	if err != nil {
		// 		replyInternalErrorMessage(gb, userID)
		// 		return err
		// 	}

		// 	if len(players) < 5 {
		// 		for i := range players {
		// 			go func(id string) {
		// 				replyMessage(gb, id, "Need at least 5 players to begin the game")
		// 			}(players[i].ID)
		// 		}
		// 		return nil
		// 	}

		// 	v, err := vocab.Get(gb.Session)

		// 	if err != nil {
		// 		replyInternalErrorMessage(gb, userID)
		// 		return err
		// 	}

		// 	var normalWord string
		// 	var undercoverWord string
		// 	if rand.Intn(2) == 0 {
		// 		normalWord = v.First
		// 		undercoverWord = v.Second
		// 	} else {
		// 		normalWord = v.Second
		// 		undercoverWord = v.First
		// 	}

		// 	roleNumList := make([]int, 10)
		// 	roleNumList[0] = 0
		// 	roleNumList[0] = 1
		// 	roleNumList[0] = 2
		// 	roleNumList[0] = 2
		// 	roleNumList[0] = 2

		// 	if len(players) == 6 {
		// 		roleNumList = append(roleNumList, 1)
		// 	}

		// 	if len(players) == 7 {
		// 		roleNumList = append(roleNumList, 0)
		// 	}

		// 	shuffledList := make([]int, len(players))
		// 	perm := rand.Perm(len(players))
		// 	for i, v := range perm {
		// 		shuffledList[v] = roleNumList[i]
		// 	}

		// 	for i := range shuffledList {
		// 		var userWord string
		// 		switch shuffledList[i] {
		// 		case 0:
		// 			userWord = ""
		// 		case 1:
		// 			userWord = undercoverWord
		// 		case 2:
		// 			userWord = normalWord
		// 		}
		// 		go func(i int, userWord string) {
		// 			err := user.AddRole(gb.Session, players[i].ID, shuffledList[i], userWord)
		// 			replyMessage(gb, players[i].ID, userWord)
		// 		}(i, userWord)
		// 	}

	} else {
		replyDefaultMessage(gb, userID)

	}

	return nil
}

func replyInternalErrorMessage(gb *GameBot, userID string) {
	message := `ระบบขัดข้อง กรุณาลองใหม่`
	replyMessage(gb, userID, message)
}

func replyDefaultMessage(gb *GameBot, userID string) {
	message := `ทำตามคำสั่งด้านล่างเท่านั้นนะจ๊ะ
- create {เลขห้อง} : สร้างห้องเพื่อเล่นเกม
- join {เลขห้อง} {ชื่อที่ใช้เล่นเกม} : เข้าห้องเพื่อรอเล่นเกม
- leave : ออกจากห้องเกมปัจจุบัน
- help : แสดงข้อความคำสั่งทั้งหมด`

	replyMessage(gb, userID, message)
}

func replyMessage(gb *GameBot, userID string, message string) {
	gb.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
}
