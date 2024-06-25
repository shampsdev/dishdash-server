package domain

type Swipe struct {
	LobbyID string
	CardID  int64
	UserID  string
	Type    SwipeType
}

type SwipeType string

var (
	LIKE    SwipeType = "like"
	DISLIKE SwipeType = "dislike"
)
