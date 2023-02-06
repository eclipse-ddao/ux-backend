package main

import (
	"log"

	"github.com/eclipse-ddao/eclipse-backend/app"
	"github.com/eclipse-ddao/eclipse-backend/controllers"
	"github.com/eclipse-ddao/eclipse-backend/database"
)

func main() {
	fiberApp := controllers.New()
	db := database.Connect()
	app := app.New(fiberApp, db)

	// Setup the controller to have all app properties
	controller := controllers.Controller{App: app}
	controller.SetupRoutes()
	log.Fatal(fiberApp.Listen(":3000"))

}
