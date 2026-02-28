package seed

import (
	"log"

	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/utils"
)

func SeedAdmin() {
	// Hashing password sebelum di  kirim ke DB
	password, _ := utils.HashPassword("admin123")

	// 
	admin := models.User{
		Name: "Admin Aplikasi",
		Email: "admin@example.com",
		Password: password,
		Role: "admin",
	}

	// Insert data ke DB dan kondisi jika email sudah ada
	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil{
		log.Println("Failed to sedd admin",err)
	} else {
		log.Println("Admin user succesfully seeded")
	}
}