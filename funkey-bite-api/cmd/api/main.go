package main

import (
	"log"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

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

	if err := utils.ConfigureJWTSecretFromEnv(); err != nil {
		log.Fatal(err)
	}

	// Initialize database
	db := database.InitializeDatabase()
	defer database.CloseDatabase(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	cateringRepo := repository.NewCateringRepository(db)
	adminRepo := repository.NewAdminRepository(db)
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
	inventoryService := services.NewInventoryService(inventoryRepo, menuRepo)
	orderService := services.NewOrderService(orderRepo, menuRepo, inventoryService, notificationService)
	userService := services.NewUserService(*userRepo, *orderRepo)
	menuService := services.NewMenuService(*menuRepo)
	cateringService := services.NewCateringService(*cateringRepo, notificationService)
	adminService := services.NewAdminService(adminRepo, *orderRepo, *userRepo, *cateringRepo, *menuRepo)
	settingsService := services.NewSettingsService(*settingsRepo)
	promotionService := services.NewPromotionService(promotionRepo)

	// Initialize handlers
	authHandler := v1.NewAuthHandler(authService, userService)
	orderHandler := v1.NewOrderHandler(orderService, authService, settingsService, promotionService)
	menuHandler := v1.NewMenuHandler(menuService)
	cateringHandler := v1.NewCateringHandler(cateringService)
	adminHandler := v1.NewAdminHandler(adminService)
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
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/auth/check", authHandler.CheckUser)

		public.GET("/settings", settingsHandler.GetPublicSettings)
		public.GET("/settings/hours", settingsHandler.GetOpeningHours)

		public.POST("/catering/requests", cateringHandler.CreateRequest)
		public.GET("/promotions/validate", promotionHandler.ValidatePromotion)
		public.GET("/promotions/active", promotionHandler.GetActivePromotions)
		public.POST("/admin/auth/login", adminHandler.AdminLogin)

	}

	menuRoutes := public.Group("/menu")
	{
		menuRoutes.GET("/", menuHandler.GetMenu)
		menuRoutes.GET("/search", menuHandler.SearchMenu)
		menuRoutes.GET("/featured", menuHandler.GetFeaturedItems)
		menuRoutes.GET("/tags", menuHandler.GetMenuByTags)
		menuRoutes.GET("/:id", menuHandler.GetMenuItem)
		menuRoutes.GET("/category", menuHandler.GetMenuByCategory)
	}
	orderGroup := public.Group("/orders")
	orderGroup.Use(middleware.OptionalAuthMiddleware())
	{
		orderGroup.POST("/", orderHandler.CreateOrder)
		orderGroup.GET("/track/:orderNumber", middleware.TrackingRateLimitMiddleware(), orderHandler.TrackOrder)
		orderGroup.PATCH("/:id/cancel", orderHandler.CancelOrder)
	}

	userRoutes := public.Group("/auth")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/profile", authHandler.GetProfile)
		userRoutes.PUT("/profile", authHandler.UpdateProfile)
		userRoutes.PATCH("/password", authHandler.ChangePassword)
		userRoutes.GET("/orders", orderHandler.GetUserOrders)
		userRoutes.GET("/orders/:id", orderHandler.GetUserOrder)
		userRoutes.GET("/catering-requests", cateringHandler.GetUserRequests)
	}

	admin := public.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware(adminRepo))
	{
		admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
		admin.GET("/dashboard/stats/today", adminHandler.GetTodayStats)

		admin.GET("/reports/sales", adminHandler.GetSalesReport)

		admin.GET("/orders", adminHandler.GetAllOrders)
		admin.PATCH("/orders/:id/status", adminHandler.UpdateOrderStatus)

		admin.GET("/users", adminHandler.GetAllUsers)
		admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)

		admin.POST("/menu/items", adminHandler.CreateMenuItem)
		admin.PUT("/menu/items/:id", adminHandler.UpdateMenuItem)
		admin.DELETE("/menu/items/:id", adminHandler.DeleteMenuItem)

		admin.GET("/catering/requests", cateringHandler.GetAllRequests)
		admin.PATCH("/catering/requests/:id/status", cateringHandler.UpdateRequestStatus)

		admin.GET("/settings", settingsHandler.GetSettings)
		admin.PUT("/settings", settingsHandler.UpdateSettings)

		admin.POST("/promotions", promotionHandler.CreatePromotion)
		admin.GET("/promotions", promotionHandler.GetPromotions)
		admin.GET("/promotions/:id", promotionHandler.GetPromotion)
		admin.PUT("/promotions/:id", promotionHandler.UpdatePromotion)
		admin.DELETE("/promotions/:id", promotionHandler.DeletePromotion)

		admin.POST("/auth/logout", adminHandler.AdminLogout)
		admin.PATCH("/auth/password", adminHandler.UpdateAdminPassword)

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
