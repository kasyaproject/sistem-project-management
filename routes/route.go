package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kasyaproject/sistem-project-management/controllers"
)

func Setup(app *fiber.App, uc *controllers.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Routes
	app.Post("/v1/auth/register", uc.Register)
}
