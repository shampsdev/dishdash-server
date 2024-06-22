package domain

type Swipe struct {
	LobbyID int64
	CardID  int64
	UserID  string
	Type    SwipeType
}

type SwipeType string

var (
	LIKE    SwipeType = "like"
	DISLIKE SwipeType = "dislike"
)
