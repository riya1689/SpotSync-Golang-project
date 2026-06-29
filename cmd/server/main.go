package main 

import (
	"log" 
	"os"  

	"spotSync-golang-Project/config"    
	"spotSync-golang-Project/internal/handler"   
	"spotSync-golang-Project/internal/middleware"
	"spotSync-golang-Project/internal/models"     
	"spotSync-golang-Project/internal/repository" 
	"spotSync-golang-Project/internal/service"   

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"           
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() { 
	if err := godotenv.Load(); err != nil {
		log.Println("No environment variables found") 
	}

	// Connect to Database
	config.ConnectDB()

	// Migrate Database 
	models.MigrateModels(config.DB) 

	// Initialize Validator
	validate := validator.New()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB) 
	zoneRepo := repository.NewZoneRepository(config.DB) 
	resRepo := repository.NewReservationRepository(config.DB) 

	// Initialize Services
	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	resService := service.NewReservationService(resRepo)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authService, validate)
	zoneHandler := handler.NewZoneHandler(zoneService, validate)
	resHandler := handler.NewReservationHandler(resService, validate)

	// Setup Echo
	e := echo.New()

	// Global Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// API Group
	api := e.Group("/api/v1")

	// Public Routes for Auth
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Public Routes for Zones
	api.GET("/zones", zoneHandler.GetAllZones)
	api.GET("/zones/:id", zoneHandler.GetZoneByID)

	// Protected Routes
	protected := api.Group("") 
	protected.Use(middleware.JWTAuth()) 

	// Admin Routes for Zones 
	adminZones := protected.Group("/zones")
	adminZones.Use(middleware.AdminOnly())
	adminZones.POST("", zoneHandler.CreateZone)

	// Protected Routes for Reservations
	reservations := protected.Group("/reservations")
	reservations.POST("", resHandler.ReserveSpot)
	reservations.GET("/my-reservations", resHandler.GetMyReservations)
	reservations.DELETE("/:id", resHandler.CancelReservation)

	// Admin Routes for Reservations
	reservations.GET("", resHandler.GetAllReservations, middleware.AdminOnly())

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
