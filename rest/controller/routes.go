package controller

import (
	"github.com/ItsClairton/Anny/rest/controller/interactions"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/interactions", interactions.Post)
}
