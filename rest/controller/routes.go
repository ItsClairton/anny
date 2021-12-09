package controller

import (
	"github.com/ItsClairton/Anny/rest/controller/player"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/player/:id", player.Index)
}
