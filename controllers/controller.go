package controllers

import (
	"log"

	"github.com/eclipse-ddao/eclipse-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Controller struct {
	*app.App
}

func New() *fiber.App {
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Update this
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	return fiberApp
}

const API_UPDATE_VERSION = 6

func (c *Controller) SetupRoutes() {
	log.Println("Setting up routes")
	c.FiberApp.Get("/health", func(ctx *fiber.Ctx) error {
		ctx.JSON(fiber.Map{
			"success":     true,
			"api_version": API_UPDATE_VERSION,
		})
		return nil
	})

	c.SetupUserRoutes()
	c.SetupDaoRoutes()
	c.SetupFileRoutes()
	c.SetupBigFileRoutes()
	c.SetupStorageProviderRoutes()
}
