package domain

type LobbySettings struct {
	LobbyID     string
	PriceMin    int
	PriceMax    int
	MaxDistance float64
	Tags        []*Tag
}
