package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	AppConfig *Config
)

// Type dari env
type Config struct {
	AppPort string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string

	JWTSecret string
	JWTRefreshToken string
	JWTExpire string
}

// Function untuk memanggil file ENV  
func LoadEnv(){
 	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found!")
	}

	AppConfig = &Config{
		AppPort: getEnv("PORT", "3000"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "swordfish"),
		DBName: getEnv("DB_NAME", "sistem_project_management"),

		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRES_IN","24h"),
		JWTExpire: getEnv("JWT_EXPIRES_IN","1h"),
	}
}

// Mengambil nilai dari env yang sudah di panggil
func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)

	if exist{
		// Jika ada return value env nya
		return value
	}else{
		// Jika tidak return default
		return fallback
	}
}

// Function untuk connect ke DB
func ConncetDB(){
	 // Ambil konfigurasi database dari struct AppConfig
	cfg := AppConfig

	// Buat string koneksi ke PostgresSQL 
	dsn:= fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName )

	// Koneksi ke DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("Failed to connect to database", err)
	}

	 // Mengambil instance database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance", err)
	}

	// configurasi koneksi ke db
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// Simpan instance koneksi ke variable global DB
	DB = db
}