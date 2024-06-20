package card

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/pkg/filter"
)

type tagOutput struct {
	Id   int64  `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
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

func tagToOutput(t *domain.Tag) tagOutput {
	return tagOutput{
		Id:   t.ID,
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
		Tags:             filter.Map(c.Tags, tagToOutput),
	}
}
