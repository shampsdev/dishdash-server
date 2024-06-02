package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CardToCreate struct {
	Title            string `json:"title" validate:"required,min=1,max=255"`
	ShortDescription string `json:"shortDescription" validate:"required,max=255"`
	Description      string `json:"description" validate:"required"`
	Image            string `json:"image" validate:"required,url,max=255"`
	Location         string `json:"location" validate:"required,max=255"`
	Address          string `json:"address" validate:"required,max=255"`
	Price            int    `json:"price" validate:"required,gt=0"`
}

type Card struct {
	ID               int64  `json:"id"`
	Title            string `json:"title" `
	ShortDescription string `json:"shortDescription"`
	Description      string `json:"description"`
	Image            string `json:"image"`
	Location         string `json:"location"`
	Address          string `json:"address"`
	Price            int    `json:"price"`
	Tags             []Tag  `json:"tags"`
}

func ValidateCardToCreate(card CardToCreate) validator.ValidationErrors {
	err := validate.Struct(card)
	if err != nil {
		var errs validator.ValidationErrors
		_ = errors.As(err, &errs)
		return errs
	}
	return nil
}
