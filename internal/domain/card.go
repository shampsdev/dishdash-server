package domain

import (
	"dishdash.ru/internal/dto"
	geo "github.com/kellydunn/golang-geo"
)

// Card represents a card with various attributes.
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

func (c *Card) ToDto() dto.Card {
	cardDto := dto.Card{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Address:          c.Address,
		Type:             c.Type,
		Price:            c.Price,
	}

	cardDto.Location = Point2String(c.Location)
	return cardDto
}

func (c *Card) ParseDto(cardDto dto.Card) error {
	c.ID = cardDto.ID
	c.Title = cardDto.Title
	c.ShortDescription = cardDto.ShortDescription
	c.Description = cardDto.Description
	c.Image = cardDto.Image
	c.Address = cardDto.Address
	c.Type = cardDto.Type
	c.Price = cardDto.Price

	return ParsePoint(cardDto.Location, c.Location)
}

func (c *Card) ParseDtoToCreate(cardDto dto.CardToCreate) error {
	c.Title = cardDto.Title
	c.ShortDescription = cardDto.ShortDescription
	c.Description = cardDto.Description
	c.Image = cardDto.Image
	c.Address = cardDto.Address
	c.Type = cardDto.Type
	c.Price = cardDto.Price

	return ParsePoint(cardDto.Location, c.Location)
}
