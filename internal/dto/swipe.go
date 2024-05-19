package dto

type SwipeType string

const (
	LIKE    SwipeType = "LIKE"
	DISLIKE SwipeType = "DISLIKE"
)

type Swipe struct {
	ID        int64     `json:"ID"`
	LobbyID   int64     `json:"lobbyID"`
	CardID    int64     `json:"cardID"`
	UserID    string    `json:"userID"`
	SwipeType SwipeType `json:"swipeType"`
}
