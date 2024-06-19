package domain

type Swipe struct {
	ID      int64
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
