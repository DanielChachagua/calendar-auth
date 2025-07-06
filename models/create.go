package models

import (
	"github.com/go-playground/validator"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type CreateEvent struct {
	Token *oauth2.Token `json:"token" validate:"required"`
	Event *calendar.Event `json:"event" validate:"required"`
}

func (c *CreateEvent) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}