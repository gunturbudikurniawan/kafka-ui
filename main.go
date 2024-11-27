package main

import (
	"fmt"
	"go-clean-architecture/controller"
	"go-clean-architecture/kafka"
	"go-clean-architecture/user"
	"log"
	"os"

	_ "go-clean-architecture/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// @title Go Clean Architecture API
// @version 1.0
// @description API dokumentasi untuk aplikasi Go Clean Architecture
// @host localhost:8000
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimezone,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connection to Database is good")

	// Kafka Configuration
	brokerAddresses := []string{"localhost:9092"}
	topic := "RegisterUsers"

	// Kafka Producer
	producer := kafka.NewKafkaProducer(brokerAddresses, topic)
	defer producer.Close()

	// Kafka Service
	kafkaService := user.NewKafkaService(producer)

	// Initialize User Service and Controller
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, kafkaService)
	userController := controller.NewUserController(userService)

	// Setup Routes
	router := gin.Default()
	// Route Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("api/v1")
	{
		v1.POST("/user", userController.RegisterUserInput)
		v1.POST("/users", userController.RegisterUsersInput)

	}

	err = router.Run(":8000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
