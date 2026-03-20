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
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	// List Controller
	listRepo := repositories.NewListRepository()
	listPositionRepo := repositories.NewListPositionsRepository()
	listService := services.NewListService(listRepo, boardRepo, listPositionRepo)
	listController := controllers.NewListController(listService)

	// Card Controller
	cardRepo := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepo, listRepo, userRepo)
	cardController := controllers.NewCardController(cardService)

	routes.Setup(app, userController, boardController, listController, cardController)

	log.Println("Server is running on port : ", port)
	log.Fatal(app.Listen(":" + port))
}
