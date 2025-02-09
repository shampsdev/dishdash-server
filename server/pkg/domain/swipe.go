package domain

type Swipe struct {
	ID      int64     `json:"id"`
	LobbyID string    `json:"lobbyID"`
	CardID  int64     `json:"cardID"`
	UserID  string    `json:"userID"`
	Type    SwipeType `json:"type"`
}

type SwipeType string

var (
	LIKE    SwipeType = "like"
	DISLIKE SwipeType = "dislike"
)
