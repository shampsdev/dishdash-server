package domain

type LobbySettings struct {
	ID          int
	PriceMin    int
	PriceMax    int
	MaxDistance float64
	Tags        []Tag
}
