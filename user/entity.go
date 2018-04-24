package user

// User is the user struct
type User struct {
	ID     string `bson:"id"`
	Name   string `bson:"name"`
	RoomID string `bson:"roomId"`
	Role   int    `bson:"role"`
}
