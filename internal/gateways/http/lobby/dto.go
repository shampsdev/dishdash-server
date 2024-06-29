package lobby

import (
	"time"

	"dishdash.ru/pkg/filter"

	"dishdash.ru/internal/domain"
)

type tagOutput struct {
	ID   int64  `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
}

type lobbySettingsOutput struct {
	ID          int64       `json:"id"`
	PriceMin    int         `json:"priceMin"`
	PriceMax    int         `json:"priceMax"`
	MaxDistance float64     `json:"maxDistance"`
	Tags        []tagOutput `json:"tags"`
}

type cardOutput struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	ShortDescription string            `json:"shortDescription"`
	Description      string            `json:"description"`
	Image            string            `json:"image"`
	Location         domain.Coordinate `json:"location"`
	Address          string            `json:"address"`
	PriceMin         int               `json:"priceMin"`
	PriceMax         int               `json:"priceMax"`
	Tags             []tagOutput       `json:"tags"`
}

type matchOutput struct {
	ID      int64  `json:"id"`
	CardID  int64  `json:"cardID"`
	LobbyID string `json:"lobbyID"`
}

type swipeOutput struct {
	UserID  string `json:"userID"`
	CardID  int64  `json:"cardID"`
	LobbyID string `json:"lobbyID"`
	Type    string `json:"type"`
}

type finalVoteOutput struct {
	ID      int64  `json:"id"`
	CardID  int64  `json:"cardID"`
	UserID  string `json:"userID"`
	LobbyID string `json:"lobbyID"`
}

type lobbyOutput struct {
	ID            string               `json:"id"`
	CreatedAt     time.Time            `json:"createdAt"`
	Location      domain.Coordinate    `json:"location"`
	LobbySettings *lobbySettingsOutput `json:"lobbySettings"`
	Cards         []cardOutput         `json:"cards"`
	Matches       []*matchOutput       `json:"matches"`
	FinalVotes    []*finalVoteOutput   `json:"finalVotes"`
	Swipes        []*swipeOutput       `json:"swipes"`
}

type nearestLobbyOutput struct {
	Dist  float64     `json:"distance"`
	Lobby lobbyOutput `json:"lobby"`
}

type findLobbyInput struct {
	Dist     float64           `json:"dist"`
	Location domain.Coordinate `json:"location"`
}

func lobbySettingsToOutput(settings *domain.LobbySettings) *lobbySettingsOutput {
	if settings == nil {
		return nil
	}
	return &lobbySettingsOutput{
		PriceMin:    settings.PriceMin,
		PriceMax:    settings.PriceMax,
		MaxDistance: settings.MaxDistance,
		Tags:        filter.Map(settings.Tags, tagToOutput),
	}
}

func tagToOutput(t domain.Tag) tagOutput {
	return tagOutput{
		ID:   t.ID,
		Icon: t.Icon,
		Name: t.Name,
	}
}

func cardToOutput(c *domain.Card) cardOutput {
	return cardOutput{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Location:         c.Location,
		Address:          c.Address,
		PriceMin:         c.PriceMin,
		PriceMax:         c.PriceMax,
		Tags: filter.Map(c.Tags, func(t *domain.Tag) tagOutput {
			return tagToOutput(*t)
		}),
	}
}

func matchToOutput(match *domain.Match) *matchOutput {
	if match == nil {
		return nil
	}
	return &matchOutput{
		ID:      match.ID,
		CardID:  match.CardID,
		LobbyID: match.LobbyID,
	}
}

func swipeToOutput(swipe *domain.Swipe) *swipeOutput {
	if swipe == nil {
		return nil
	}
	return &swipeOutput{
		UserID:  swipe.UserID,
		CardID:  swipe.CardID,
		LobbyID: swipe.LobbyID,
		Type:    string(swipe.Type),
	}
}

func finalVoteToOutput(finalVote *domain.FinalVote) *finalVoteOutput {
	if finalVote == nil {
		return nil
	}
	return &finalVoteOutput{
		ID:      finalVote.ID,
		CardID:  finalVote.CardID,
		UserID:  finalVote.UserID,
		LobbyID: finalVote.LobbyID,
	}
}

func lobbyToOutput(lobby *domain.Lobby) lobbyOutput {
	return lobbyOutput{
		ID:            lobby.ID,
		CreatedAt:     lobby.CreatedAt,
		Location:      lobby.Location,
		LobbySettings: lobbySettingsToOutput(lobby.LobbySettings),
		Cards:         filter.Map(lobby.Cards, cardToOutput),
		Matches:       filter.Map(lobby.Matches, matchToOutput),
		FinalVotes:    filter.Map(lobby.FinalVotes, finalVoteToOutput),
		Swipes:        filter.Map(lobby.Swipes, swipeToOutput),
	}
}
