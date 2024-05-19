package domain

import (
	"errors"

	"dishdash.ru/internal/dto"

	geo "github.com/kellydunn/golang-geo"
)

type Card struct {
	ID               int64
	Title            string
	ShortDescription string
	Description      string
	Image            string
	Location         *geo.Point
	Address          string
	Type             dto.CardType
	Price            int
}

func CardToDto(c Card) dto.Card {
	// TODO move to Card.ToDto
	cardDto := dto.Card{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Address:          c.Address,
		Location:         "",
		Type:             c.Type,
		Price:            c.Price,
	}

	cardDto.Location = Point2String(c.Location)
	return cardDto
}

func CardFromDtoToCreate(c dto.CardToCreate) (*Card, error) {
	card := &Card{
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Location:         &geo.Point{},
		Address:          c.Address,
		Type:             c.Type,
		Price:            c.Price,
	}

	var err error
	card.Location, err = ParsePoint(c.Location)
	if err != nil {
		err = errors.New("can't parse location")
	}
	return card, err
}

func CardFromDto(c dto.Card) (*Card, error) {
	card := &Card{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Location:         &geo.Point{},
		Address:          c.Address,
		Type:             c.Type,
		Price:            c.Price,
	}

	var err error
	card.Location, err = ParsePoint(c.Location)
	if err != nil {
		err = errors.New("can't parse location")
	}
	return card, err
}
