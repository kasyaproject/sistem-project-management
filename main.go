package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/controllers"
	"github.com/kasyaproject/sistem-project-management/database/seed"
	"github.com/kasyaproject/sistem-project-management/repositories"
	"github.com/kasyaproject/sistem-project-management/routes"
	"github.com/kasyaproject/sistem-project-management/services"
)

func main() {
	// Load env & Connect DB
	config.LoadEnv()
	config.ConncetDB()
	port := config.AppConfig.AppPort

	// Jalankan seeder
	seed.SeedAdmin()

	app := fiber.New()

	// User Controller
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Board Controller
	boardRepo := repositories.NewBoardRepository()
	boardService := services.NewBoardService(boardRepo, userRepo)
	boardController := controllers.NewBoardController(boardService)

	routes.Setup(app, userController, boardController)

	log.Println("Server is running on port : ", port)
	log.Fatal(app.Listen(":" + port))
}
