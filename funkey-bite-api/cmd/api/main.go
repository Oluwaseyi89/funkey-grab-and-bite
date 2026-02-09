package main

import (
	// "database/sql"
	// "fmt"
	"log"
	"time"

	// "os"

	"funkey-grab-and-bite/funkey-bite-api/internal/database" // Add this import

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	// "github.com/ulule/limiter"
	// "github.com/ulule/limiter/drivers/store/memory"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers/middleware"
	v1 "funkey-grab-and-bite/funkey-bite-api/internal/handlers/v1"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
	settingsRepo := repository.NewSettingsRepository(db)
	promotionRepo := repository.NewPromotionRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	emailService := utils.NewEmailService()
	smsService := utils.NewSMSService()

	notificationService := services.NewNotificationService(
		emailService,
		smsService,
		*orderRepo,
		*userRepo,
		*cateringRepo,
		*notificationRepo,
	)

	// Initialize services
	authService := services.NewAuthService(*userRepo)
	orderService := services.NewOrderService(orderRepo, *menuRepo, notificationService)
	userService := services.NewUserService(*userRepo, *orderRepo)
	menuService := services.NewMenuService(*menuRepo)                                                    // Add this
	cateringService := services.NewCateringService(*cateringRepo, notificationService)                   // Add this
	adminService := services.NewAdminService(adminRepo, *orderRepo, *userRepo, *cateringRepo, *menuRepo) // Add this
	settingsService := services.NewSettingsService(*settingsRepo)
	promotionService := services.NewPromotionService(*promotionRepo)
	inventoryService := services.NewInventoryService(*inventoryRepo, *menuRepo)

	// Initialize handlers
	authHandler := v1.NewAuthHandler(authService, userService)
	orderHandler := v1.NewOrderHandler(orderService, authService)
	menuHandler := v1.NewMenuHandler(menuService)
	cateringHandler := v1.NewCateringHandler(cateringService) // Add this
	adminHandler := v1.NewAdminHandler(adminService)          // Add this
	settingsHandler := v1.NewSettingsHandler(settingsService)
	promotionHandler := v1.NewPromotionHandler(promotionService)
	inventoryHandler := v1.NewInventoryHandler(inventoryService)

	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Hour,
		Limit:  100,
	}
	limiterInstance := limiter.New(store, rate)

	// Setup router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Public routes
	public := r.Group("/api/v1")
	public.Use(ginLimiter.NewMiddleware(limiterInstance))

	public.GET("/order/track/:phone/:orderNumber", middleware.TrackingRateLimitMiddleware(), orderHandler.TrackOrderPublic)

	{
		// public.GET("/menu", menuHandler.GetMenu)
		// public.GET("/menu/:id", menuHandler.GetMenuItem)
		// public.GET("/categories", menuHandler.GetCategories)
		// public.GET("/menu/category", menuHandler.GetMenuByCategory)
		// public.GET("/menu/search", menuHandler.SearchMenu)
		// public.GET("/menu/featured", menuHandler.GetFeaturedItems)

		// Auth routes
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/auth/check", authHandler.CheckUser)

		public.GET("/settings", settingsHandler.GetPublicSettings)
		public.GET("/settings/hours", settingsHandler.GetOpeningHours)

		// Order with optional auth middleware
		// orderGroup := public.Group("/orders")
		// orderGroup.Use(middleware.OptionalAuthMiddleware())
		// {
		// 	orderGroup.POST("/", orderHandler.CreateOrder)
		// }

		// Catering (no auth required)
		public.POST("/catering/requests", cateringHandler.CreateRequest)
		public.GET("/promotions/validate", promotionHandler.ValidatePromotion)
		public.GET("/promotions/active", promotionHandler.GetActivePromotions)
		public.POST("/admin/auth/login", adminHandler.AdminLogin)

	}

	// Add new menu routes
	menuRoutes := public.Group("/menu")
	{
		menuRoutes.GET("/", menuHandler.GetMenu)
		menuRoutes.GET("/search", menuHandler.SearchMenu)
		menuRoutes.GET("/featured", menuHandler.GetFeaturedItems)
		menuRoutes.GET("/tags", menuHandler.GetMenuByTags)
		menuRoutes.GET("/:id", menuHandler.GetMenuItem)
		menuRoutes.GET("/category", menuHandler.GetMenuByCategory)
	}
	// To this (already exists):
	orderGroup := public.Group("/orders")
	orderGroup.Use(middleware.OptionalAuthMiddleware())
	{
		orderGroup.POST("/", orderHandler.CreateOrder)
		// Add guest order tracking
		orderGroup.GET("/track/:orderNumber", middleware.TrackingRateLimitMiddleware(), orderHandler.TrackOrder)
		orderGroup.PATCH("/:id/cancel", orderHandler.CancelOrder) // Add this line

		// orderGroup.GET("/track/:phone/:orderNumber", orderHandler.TrackOrderPublic)
	}

	// Add user profile routes
	userRoutes := public.Group("/auth")
	userRoutes.Use(middleware.AuthMiddleware()) // Protect these routes
	{
		userRoutes.GET("/profile", authHandler.GetProfile)
		userRoutes.PUT("/profile", authHandler.UpdateProfile)
		userRoutes.PATCH("/password", authHandler.ChangePassword)
		userRoutes.GET("/orders", orderHandler.GetUserOrders)
		userRoutes.GET("/orders/:id", orderHandler.GetUserOrder)
		userRoutes.GET("/catering-requests", cateringHandler.GetUserRequests)
	}

	// Protected routes (for admin dashboard)
	// admin := public.Group("/admin")
	// admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	admin := public.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// Admin endpoints for managing orders, menu, etc.
		// Dashboard
		admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
		admin.GET("/dashboard/stats/today", adminHandler.GetTodayStats)

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

		admin.GET("/settings", settingsHandler.GetSettings)
		admin.PUT("/settings", settingsHandler.UpdateSettings)

		admin.POST("/promotions", promotionHandler.CreatePromotion)
		admin.GET("/promotions", promotionHandler.GetPromotions)
		admin.GET("/promotions/:id", promotionHandler.GetPromotion)
		admin.PUT("/promotions/:id", promotionHandler.UpdatePromotion)
		admin.DELETE("/promotions/:id", promotionHandler.DeletePromotion)

		// Admin Auth (protected)
		admin.POST("/auth/logout", adminHandler.AdminLogout)
		admin.PATCH("/auth/password", adminHandler.UpdateAdminPassword)

		// Admin User Management
		admin.GET("/users/admins", adminHandler.GetAdminUsers)
		admin.POST("/users/admins", adminHandler.CreateAdminUser)
		admin.PUT("/users/admins/:id", adminHandler.UpdateAdminUser)
		admin.DELETE("/users/admins/:id", adminHandler.DeleteAdminUser)

	}

	inventoryRoutes := admin.Group("/inventory")
	{
		inventoryRoutes.GET("/", inventoryHandler.GetInventory)
		inventoryRoutes.GET("/dashboard", inventoryHandler.GetDashboard)
		inventoryRoutes.GET("/low-stock", inventoryHandler.GetLowStock)
		inventoryRoutes.GET("/alerts", inventoryHandler.GetAlerts)
		inventoryRoutes.GET("/check", inventoryHandler.CheckAvailability)
		inventoryRoutes.GET("/menu-item/:menuItemId", inventoryHandler.GetInventoryByMenuItem)
		inventoryRoutes.GET("/:id", inventoryHandler.GetInventoryItem)

		inventoryRoutes.POST("/", inventoryHandler.CreateInventoryItem)
		inventoryRoutes.POST("/restock", inventoryHandler.RestockItem)

		inventoryRoutes.PATCH("/stock", inventoryHandler.UpdateStock)
		inventoryRoutes.PATCH("/alerts/:id/resolve", inventoryHandler.ResolveAlert)
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
