package models

import (
	"fmt"

	"github.com/go-playground/validator"
	"golang.org/x/oauth2"
)

type CreateEvent struct {
	Token *oauth2.Token `json:"token" validate:"required"`
	Event CalendarCreate `json:"event" validate:"required"`
}

func (c *CreateEvent) Validate() error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateFormat", validateDateFormat)
	_ = validate.RegisterValidation("timeFormat", validateTimeFormat)

	return validate.Struct(c)
}

type CalendarCreate struct {
	Summary     string     `json:"title"`
	Location    *string     `json:"description"`
	Description *string    `json:"url"`
	Date        CustomDate `json:"date" validate:"required,dateFormat"`
	Time        *CustomTime `json:"time" validate:"omitempty,timeFormat"`
}



func (n *CalendarCreate) Validate() error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateFormat", validateDateFormat)
	_ = validate.RegisterValidation("timeFormat", validateTimeFormat)

	err := validate.Struct(n)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()

	return fmt.Errorf("campo %s es inválido, revisar: (%s)", field, tag)
}

