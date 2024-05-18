package domain

import "dishdash.ru/internal/dto"

type Card struct {
	ID               int
	Title            string
	ShortDescription string
	Description      string
	Image            string
	Location         string
	Address          string
	Type             dto.CardType
	Price            int
}

func CardToDto(c Card) dto.Card {
	return dto.Card{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Location:         c.Location,
		Address:          c.Address,
		Type:             c.Type,
		Price:            c.Price,
	}
}

func CardFromDto(c dto.CardToCreate) *Card {
	return &Card{
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Location:         c.Location,
		Address:          c.Address,
		Type:             c.Type,
		Price:            c.Price,
	}
}
