package user

// User is the user struct
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	RoomID string `json:"roomId"`
}
