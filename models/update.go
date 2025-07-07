package models

import (
	"fmt"

	"github.com/go-playground/validator"
	"golang.org/x/oauth2"
)

type UpdateEvent struct {
	Token *oauth2.Token  `json:"token" validate:"required"`
	Event CalendarUpdate `json:"event" validate:"required"`
}

func (c *UpdateEvent) Validate() error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateFormat", validateDateFormat)
	_ = validate.RegisterValidation("timeFormat", validateTimeFormat)

	return validate.Struct(c)
}

type CalendarUpdate struct {
	ID          string      `json:"id"`
	Summary     string      `json:"summary"`
	Location    *string     `json:"location"`
	Description *string     `json:"description"`
	Date        CustomDate  `json:"date" validate:"required,dateFormat"`
	Time        *CustomTime `json:"time" validate:"omitempty,timeFormat"`
}

func (n *CalendarUpdate) Validate() error {
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
