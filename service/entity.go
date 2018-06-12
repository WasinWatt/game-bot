package service

import (
	"errors"

	"github.com/WasinWatt/game-bot/mongo"
)

// Errors
var (
	ErrNotFound = errors.New("user: not found")
)

// New creates new service controller
func New(repo *mongo.Repository) *Controller {
	return &Controller{repo}
}

// Controller controls business flows
type Controller struct {
	repo *mongo.Repository
}
