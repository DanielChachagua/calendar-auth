package controllers

import (
	"calendar_auth/models"
	"calendar_auth/services"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func GetUrl(c *fiber.Ctx) error {
	redirectUrl := c.Query("redirect_url")
	if redirectUrl == "" {
		return c.Status(400).JSON(fiber.Map{
			"details": "redirect_url es requerida",
		})
		
	}

	url, err := services.GetCalendarUrl(redirectUrl)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"url": url,
	})
}

func GetToken(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{
			"details": "code es requerido",
		})
	}

	redirectUrl := c.Query("redirect_url")
	if redirectUrl == "" {
		return c.Status(400).JSON(fiber.Map{
			"details": "redirect_url es requerida",
		})
		
	}

	token, err := services.GetCalendarToken(code, redirectUrl)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token": token,
	})
}

func GetEvents(c *fiber.Ctx) error {
	var token oauth2.Token
	if err := c.BodyParser(&token); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	events, err := services.GetCalendarEvents(&token)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"items": &events,
	})
}

func CreateEvent(c *fiber.Ctx) error {
	var event models.CreateEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	if err := event.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err := services.CreateCalendarEvent(event.Token, event.Event)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"detail": "Evento creado con exito",
	})
}

func UpdateEvent(c *fiber.Ctx) error {
	eventId := c.Params("event_id")
	if eventId == "" {
		return c.Status(400).JSON(fiber.Map{
			"details": "event_id es requerido",
		})
	}

	var event models.UpdateEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	if err := event.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err := services.UpdateCalendarEvent(event.Token, event.Event, eventId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"detail": "Evento actualizado con exito",
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	eventId := c.Params("event_id")
	if eventId == "" {
		return c.Status(400).JSON(fiber.Map{
			"details": "event_id es requerido",
		})
	}

	var token oauth2.Token
	if err := c.BodyParser(&token); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err := services.DeleteCalendarEvent(&token, eventId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Evento eliminado",
	})
}