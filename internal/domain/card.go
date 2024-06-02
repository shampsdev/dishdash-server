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
	Price            int
	Tags             []*Tag
}

func (c *Card) ToDto() dto.Card {
	tagsDto := make([]dto.Tag, len(c.Tags))
	for i, tag := range c.Tags {
		tagsDto[i] = tag.ToDto()
	}

	cardDto := dto.Card{
		ID:               c.ID,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		Image:            c.Image,
		Address:          c.Address,
		Price:            c.Price,
		Tags:             tagsDto,
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
	c.Price = cardDto.Price
	c.Location = new(geo.Point)

	tags := make([]*Tag, len(cardDto.Tags))
	for i, tagDto := range cardDto.Tags {
		tag := new(Tag)
		tag.ParseDto(tagDto)
		tags[i] = tag
	}
	c.Tags = tags

	return ParsePoint(cardDto.Location, c.Location)
}

func (c *Card) ParseDtoToCreate(cardDto dto.CardToCreate) error {
	c.Title = cardDto.Title
	c.ShortDescription = cardDto.ShortDescription
	c.Description = cardDto.Description
	c.Image = cardDto.Image
	c.Address = cardDto.Address
	c.Price = cardDto.Price
	c.Location = new(geo.Point)

	return ParsePoint(cardDto.Location, c.Location)
}
