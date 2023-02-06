package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	DB       *gorm.DB
	FiberApp *fiber.App
}

func New(f *fiber.App, d *gorm.DB) *App {

	defaultApp := &App{
		FiberApp: f,
		DB:       d,
	}
	return defaultApp
}
