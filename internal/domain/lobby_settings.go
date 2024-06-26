package domain

type LobbySettings struct {
	PriceMin    int
	PriceMax    int
	MaxDistance float64
	Tags        []Tag
}
