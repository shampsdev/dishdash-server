package domain

type LobbySettings struct {
	ID          int64
	LobbyID     string
	PriceMin    int
	PriceMax    int
	MaxDistance float64
	Tags        []*Tag
}
