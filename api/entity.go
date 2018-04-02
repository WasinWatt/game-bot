package api

import (
	"github.com/line/line-bot-sdk-go/linebot"
	mgo "gopkg.in/mgo.v2"
)

// GameBot is the app's api
type GameBot struct {
	Client  *linebot.Client
	Session *mgo.Session
}
