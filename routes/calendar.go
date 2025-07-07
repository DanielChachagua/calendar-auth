package routes

import (
	"calendar_auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func CalendarRoutes(app *fiber.App) {
	calendar := app.Group("/calendar")

	calendar.Get("/get_url", controllers.GetUrl)
	
	calendar.Post("/get_token", controllers.GetToken)

	calendar.Post("/get_events", controllers.GetEvents)

	calendar.Post("/create", controllers.CreateEvent)

	calendar.Post("/update", controllers.UpdateEvent)

	calendar.Post("/delete/", controllers.DeleteEvent)
}