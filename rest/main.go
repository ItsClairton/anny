package rest

import (
	"os"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/rest/controller"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/gofiber/fiber/v2"
)

var Module = &core.Module{StartFunc: func() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			if err, ok := e.(*fiber.Error); ok {
				return c.Status(err.Code).JSON(&fiber.Map{"data": nil, "error": err.Message})
			}

			logger.ErrorF("%s %s: %+v", c.Method(), c.Path(), e)
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": nil, "error": e.Error()})
		},
	})

	controller.New(app)

	if err := app.Listen(utils.Fmt(":%s", os.Getenv("API_PORT"))); err != nil {
		logger.FatalF("Não foi possível iniciar a API na porta %s: %v", os.Getenv("API_PORT"), err)
	}
}}
