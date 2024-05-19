package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CardType string

const (
	BAR        CardType = "BAR"
	CAFE       CardType = "CAFE"
	RESTAURANT CardType = "RESTAURANT"
)

type CardToCreate struct {
	Title            string   `json:"title" validate:"required,min=1,max=255"`
	ShortDescription string   `json:"short_description" validate:"required,max=255"`
	Description      string   `json:"description" validate:"required"`
	Image            string   `json:"image" validate:"required,url,max=255"`
	Location         string   `json:"location" validate:"required,max=255"`
	Address          string   `json:"address" validate:"required,max=255"`
	Type             CardType `json:"type" validate:"required,cardtype"`
	Price            int      `json:"price" validate:"required,gt=0"`
}

type Card struct {
	ID               int64    `json:"id"`
	Title            string   `json:"title" `
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Image            string   `json:"image"`
	Location         string   `json:"location"`
	Address          string   `json:"address"`
	Type             CardType `json:"type"`
	Price            int      `json:"price"`
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

func validateCardType(fl validator.FieldLevel) bool {
	cardType := CardType(fl.Field().String())
	switch cardType {
	case BAR, CAFE, RESTAURANT:
		return true
	}
	return false
}
