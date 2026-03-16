package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/controllers"
	"github.com/kasyaproject/sistem-project-management/database/seed"
	"github.com/kasyaproject/sistem-project-management/repositories"
	"github.com/kasyaproject/sistem-project-management/routes"
	"github.com/kasyaproject/sistem-project-management/services"
)

func main() {
	//
	config.LoadEnv()
	config.ConncetDB()
	port := config.AppConfig.AppPort

	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)
	fmt.Println("JWT Expire:", duration)
	fmt.Println(time.Now())
	// Jalankan seeder
	seed.SeedAdmin()

	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, userController)

	log.Println("Server is running on port : ", port)
	log.Fatal(app.Listen(":" + port))
}
