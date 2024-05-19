package domain

import "dishdash.ru/internal/dto"

type Swipe struct {
	ID        int64
	LobbyID   int64
	CardID    int64
	UserID    string
	SwipeType dto.SwipeType
}

func (s *Swipe) ToDto() dto.Swipe {
	return dto.Swipe{
		LobbyID:   s.LobbyID,
		CardID:    s.CardID,
		UserID:    s.UserID,
		SwipeType: s.SwipeType,
	}
}

func (s *Swipe) ParseDto(sDto dto.Swipe) {
	s.ID = sDto.ID
	s.LobbyID = sDto.LobbyID
	s.UserID = sDto.UserID
	s.CardID = sDto.CardID
	s.SwipeType = sDto.SwipeType
}
