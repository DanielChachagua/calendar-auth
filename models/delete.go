package models

import (
	"github.com/go-playground/validator"
	"golang.org/x/oauth2"
)

type DeleteEvent struct {
	Token *oauth2.Token `json:"token" validate:"required"`
	EventIds []string `json:"event_ids" validate:"required,min=1"`
}

func (c *DeleteEvent) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}