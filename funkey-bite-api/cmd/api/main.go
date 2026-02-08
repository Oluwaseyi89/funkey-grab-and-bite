package main

import (
	// "database/sql"
	// "fmt"
	"log"
	// "os"

	"funkey-grab-and-bite/funkey-bite-api/internal/database" // Add this import

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers/middleware"
	v1 "funkey-grab-and-bite/funkey-bite-api/internal/handlers/v1"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

func main() {

	err := godotenv.Load() // loads .env into environment
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db := database.InitializeDatabase()
	defer database.CloseDatabase(db) // Add this to ensure proper cleanup

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	cateringRepo := repository.NewCateringRepository(db) // Add this
	adminRepo := repository.NewAdminRepository(db)       // Add this

	// Initialize services
	authService := services.NewAuthService(*userRepo)
	orderService := services.NewOrderService(*orderRepo, *menuRepo)
	userService := services.NewUserService(*userRepo, *orderRepo)
	menuService := services.NewMenuService(*menuRepo)                                                    // Add this
	cateringService := services.NewCateringService(*cateringRepo)                                        // Add this
	adminService := services.NewAdminService(adminRepo, *orderRepo, *userRepo, *cateringRepo, *menuRepo) // Add this

	// Initialize handlers
	authHandler := v1.NewAuthHandler(authService, userService)
	orderHandler := v1.NewOrderHandler(orderService, authService)
	menuHandler := v1.NewMenuHandler(menuService)
	cateringHandler := v1.NewCateringHandler(cateringService) // Add this
	adminHandler := v1.NewAdminHandler(adminService)          // Add this

	// Setup router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/menu", menuHandler.GetMenu)
		public.GET("/menu/:id", menuHandler.GetMenuItem)
		public.GET("/categories", menuHandler.GetCategories)
		public.GET("/menu/category", menuHandler.GetMenuByCategory)
		public.GET("/menu/search", menuHandler.SearchMenu)
		public.GET("/menu/featured", menuHandler.GetFeaturedItems)

		// Auth routes
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/auth/check", authHandler.CheckUser)

		// Order with optional auth middleware
		orderGroup := public.Group("/orders")
		orderGroup.Use(middleware.OptionalAuthMiddleware())
		{
			orderGroup.POST("/", orderHandler.CreateOrder)
		}

		// Catering (no auth required)
		public.POST("/catering/requests", cateringHandler.CreateRequest)
	}

	// Protected routes (for admin dashboard)
	admin := public.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Admin endpoints for managing orders, menu, etc.
		// Dashboard
		admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)

		// Reports
		admin.GET("/reports/sales", adminHandler.GetSalesReport)

		// Orders
		admin.GET("/orders", adminHandler.GetAllOrders)
		admin.PATCH("/orders/:id/status", adminHandler.UpdateOrderStatus)

		// Users
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)

		// Menu Management
		admin.POST("/menu/items", adminHandler.CreateMenuItem)
		admin.PUT("/menu/items/:id", adminHandler.UpdateMenuItem)
		admin.DELETE("/menu/items/:id", adminHandler.DeleteMenuItem)

		// Catering Management
		admin.GET("/catering/requests", cateringHandler.GetAllRequests)
		admin.PATCH("/catering/requests/:id/status", cateringHandler.UpdateRequestStatus)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Here's the logical order to resolve these issues:

// Implement initializeDatabase() function - Database connection setup

// Create middleware package implementations - CORSMiddleware, LoggerMiddleware, AuthMiddleware, AdminMiddleware

// Implement v1.NewMenuHandler and menu endpoints - Menu API handlers

// Implement v1.NewCateringHandler with authentication - Catering requests with user auth

// Create admin routes implementation - Admin-only endpoints for data management

// Fix repository pointer issues - Remove * from *userRepo, *orderRepo, *menuRepo

// Add user profile routes - /auth/profile, /auth/orders, etc.
