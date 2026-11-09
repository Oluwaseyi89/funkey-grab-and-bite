package app

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/database"
	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers/middleware"
	v1 "funkey-grab-and-bite/funkey-bite-api/internal/handlers/v1"
	"funkey-grab-and-bite/funkey-bite-api/internal/jobs"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
	"funkey-grab-and-bite/funkey-bite-api/internal/storage"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// SetupEngine handles the complete dependency injection and routing configurations.
// It returns the configured gin instance and a cleanup function to close the database safely.
func SetupEngine() (*gin.Engine, func()) {
	err := godotenv.Load() // loads .env into environment
	if err != nil {
		log.Println("No .env file found")
	}

	if err := utils.ConfigureJWTSecretFromEnv(); err != nil {
		log.Fatal(err)
	}

	// Initialize database
	db := database.InitializeDatabase()
	maintenanceCtx, stopMaintenance := context.WithCancel(context.Background())
	var worker *jobs.Worker
	cleanup := func() {
		stopMaintenance()
		if worker != nil {
			worker.Wait()
		}
		database.CloseDatabase(db)
	}

	// Initialize static file storage (nil when STORAGE_PROVIDER is unset)
	imageStorage, err := storage.NewFromEnv(context.Background())
	if err != nil {
		log.Fatalf("Invalid storage configuration: %v", err)
	}
	if imageStorage != nil {
		models.SetImageURLResolver(func(value string) string {
			return storage.ResolveURL(imageStorage, value)
		})
		log.Printf("✅ Image storage configured (provider=%s)", os.Getenv("STORAGE_PROVIDER"))
	} else {
		log.Println("Image storage not configured; menu image uploads are disabled")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	cateringRepo := repository.NewCateringRepository(db)
	cateringPackageRepo := repository.NewCateringPackageRepository(db)
	cateringContractRepo := repository.NewCateringContractRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)
	promotionRepo := repository.NewPromotionRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	idempotencyRepo := repository.NewIdempotencyRepository(db)
	rateLimitRepo := repository.NewRateLimitRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	jobRepo := repository.NewJobRepository(db)
	worker = jobs.NewWorker(jobRepo)
	checkoutRepo := repository.NewCheckoutRepository(db, orderRepo, userRepo, inventoryRepo, promotionRepo, auditRepo, idempotencyRepo, jobRepo)

	emailService := utils.NewEmailService()
	smsService := utils.NewSMSService()
	paymentGateway := utils.NewPaymentGateway()

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
	orderService := services.NewOrderService(orderRepo, menuRepo, inventoryService, notificationService, checkoutRepo)
	userService := services.NewUserService(*userRepo, *orderRepo)
	menuService := services.NewMenuService(*menuRepo)
	cateringService := services.NewCateringService(cateringRepo, cateringPackageRepo, menuRepo, notificationService)
	adminService := services.NewAdminService(adminRepo, *orderRepo, *userRepo, *cateringRepo, *menuRepo, imageStorage)
	settingsService := services.NewSettingsService(*settingsRepo)
	promotionService := services.NewPromotionService(promotionRepo)
	paymentService := services.NewPaymentService(paymentGateway, orderRepo, checkoutRepo)
	invoiceService := services.NewInvoiceService(settingsService, emailService)
	teamService := services.NewTeamService(teamRepo, imageStorage)
	cateringContractService := services.NewCateringContractService(cateringRepo, cateringContractRepo, settingsService, emailService, paymentGateway)

	// Initialize handlers
	accountService := services.NewAccountService(userRepo, repository.NewAccountCodeRepository(db), jobEmailQueue{jobs: jobRepo, worker: worker}, auditRepo)
	authHandler := v1.NewAuthHandler(authService, userService, accountService)
	orderHandler := v1.NewOrderHandler(orderService, authService, settingsService, promotionService, paymentService, invoiceService)
	paymentHandler := v1.NewPaymentHandler(paymentService)
	adminPaymentsHandler := v1.NewAdminPaymentsHandler(repository.NewAdminPaymentRepository(db), checkoutRepo, paymentService, jobRepo, auditRepo)
	menuHandler := v1.NewMenuHandler(menuService)
	cateringHandler := v1.NewCateringHandler(cateringService)
	cateringContractHandler := v1.NewCateringContractHandler(cateringContractService, cateringHandler)
	paymentService.SetCateringHandler(cateringContractHandler.HandlePaystackCharge)
	teamHandler := v1.NewTeamHandler(teamService)
	adminHandler := v1.NewAdminHandler(adminService)
	uploadHandler := v1.NewUploadHandler(imageStorage)
	settingsHandler := v1.NewSettingsHandler(settingsService)
	promotionHandler := v1.NewPromotionHandler(promotionService)
	inventoryHandler := v1.NewInventoryHandler(inventoryService)
	healthHandler := v1.NewHealthHandler()
	realtimeHandler := v1.NewRealtimeHandler(adminRepo)

	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Hour,
		Limit:  100,
	}
	limiterInstance := limiter.New(store, rate)

	// Setup router
	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())

	trustedProxies := []string{"127.0.0.1", "::1"}
	if configuredTrustedProxies := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES")); configuredTrustedProxies != "" {
		trustedProxies = []string{}
		for _, proxy := range strings.Split(configuredTrustedProxies, ",") {
			value := strings.TrimSpace(proxy)
			if value != "" {
				trustedProxies = append(trustedProxies, value)
			}
		}
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	// Behind a platform proxy (e.g. Render, which sits behind Cloudflare), the
	// connection comes from the proxy, so per-IP limits would treat every
	// customer as one. CLIENT_IP_HEADER names the header the platform sets to
	// the real client IP (e.g. True-Client-IP); it must be one clients cannot set.
	if header := strings.TrimSpace(os.Getenv("CLIENT_IP_HEADER")); header != "" {
		r.TrustedPlatform = header
	}

	startMaintenance(maintenanceCtx, idempotencyRepo, rateLimitRepo, jobRepo)
	registerJobs(worker, orderService, notificationService, invoiceService, paymentService, emailService)
	cateringContractService.SetEmailQueue(jobEmailQueue{jobs: jobRepo, worker: worker})
	worker.Start(maintenanceCtx)

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ResponseEnvelopeMiddleware())

	// Stable health endpoint used by frontend fallback logic
	r.GET("/health", healthHandler.GetHealthLegacy)

	// Public routes
	public := r.Group("/api/v1")
	public.Use(ginLimiter.NewMiddleware(limiterInstance))
	public.GET("/health", healthHandler.GetHealth)

	// Per phone number or signed-in customer, on top of the per-IP limits; each
	// group has its own allowance.
	limitAuthByPhone := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "auth")
	limitOrdersByPhone := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "orders")
	limitLookupsByPhone := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "lookup")
	limitCateringByPhone := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "catering")
	limitBookingByPhone := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "booking")
	limitAccountByUser := middleware.IdentityRateLimitMiddleware(rateLimitRepo, "account")

	public.GET("/order/track/:phone/:orderNumber", middleware.TrackingRateLimitMiddleware(), limitLookupsByPhone, orderHandler.TrackOrderPublic)

	{
		public.POST("/auth/register", limitAuthByPhone, authHandler.Register)
		public.POST("/auth/register/confirm", middleware.TrackingRateLimitMiddleware(), limitAuthByPhone, authHandler.ConfirmRegistration)
		public.POST("/auth/password/forgot", middleware.TrackingRateLimitMiddleware(), limitAuthByPhone, authHandler.ForgotPassword)
		public.POST("/auth/password/reset", middleware.TrackingRateLimitMiddleware(), limitAuthByPhone, authHandler.ResetPassword)
		public.POST("/auth/login", limitAuthByPhone, authHandler.Login)
		public.GET("/auth/check", authHandler.CheckUser)

		public.GET("/settings", settingsHandler.GetPublicSettings)
		public.GET("/settings/hours", settingsHandler.GetOpeningHours)

		public.GET("/catering/packages", cateringHandler.ListPackages)
		public.GET("/team", teamHandler.ListPublic)
		public.POST("/catering/requests", limitCateringByPhone, cateringHandler.CreateRequest)
		public.POST("/catering/requests/lookup", middleware.TrackingRateLimitMiddleware(), limitBookingByPhone, cateringHandler.LookupRequest)
		public.POST("/catering/booking", middleware.TrackingRateLimitMiddleware(), limitBookingByPhone, cateringContractHandler.ViewBooking)
		public.POST("/catering/booking/accept", middleware.TrackingRateLimitMiddleware(), limitBookingByPhone, cateringContractHandler.AcceptQuote)
		public.POST("/catering/booking/pay", middleware.TrackingRateLimitMiddleware(), limitBookingByPhone, cateringContractHandler.PayOnline)
		public.POST("/catering/booking/pay/verify", middleware.TrackingRateLimitMiddleware(), limitBookingByPhone, cateringContractHandler.VerifyPayment)
		public.GET("/promotions/validate", promotionHandler.ValidatePromotion)
		public.GET("/promotions/active", promotionHandler.GetActivePromotions)
		public.POST("/admin/auth/login", adminHandler.AdminLogin)
		public.GET("/admin/realtime/ws", realtimeHandler.ConnectAdmin)

		// Paystack calls this directly; authenticity comes from the
		// x-paystack-signature check, not from a bearer token.
		public.POST("/payments/webhook", paymentHandler.Webhook)
	}

	menuRoutes := public.Group("/menu")
	{
		menuRoutes.GET("", menuHandler.GetMenu)
		menuRoutes.GET("/categories", menuHandler.GetCategories)
		menuRoutes.GET("/search", menuHandler.SearchMenu)
		menuRoutes.GET("/featured", menuHandler.GetFeaturedItems)
		menuRoutes.GET("/tags", menuHandler.GetMenuByTags)
		menuRoutes.GET("/:id", menuHandler.GetMenuItem)
		menuRoutes.GET("/category", menuHandler.GetMenuByCategory)
	}

	orderGroup := public.Group("/orders")
	orderGroup.Use(middleware.OptionalAuthMiddleware())
	{
		// Idempotency runs first so a retried checkout replays its result
		// without counting against the rate limits again.
		orderGroup.POST("",
			middleware.Idempotency(idempotencyRepo, middleware.IdempotencyOptions{Replay: orderHandler.ReplayCreateOrder}),
			middleware.OrderRateLimitMiddleware(rateLimitRepo), limitOrdersByPhone, orderHandler.CreateOrder)
		orderGroup.GET("/track/:orderNumber", middleware.TrackingRateLimitMiddleware(), limitLookupsByPhone, orderHandler.TrackOrder)
		orderGroup.POST("/lookup", middleware.TrackingRateLimitMiddleware(), limitLookupsByPhone, orderHandler.LookupOrder)
		orderGroup.POST("/payment-details", middleware.TrackingRateLimitMiddleware(), limitLookupsByPhone, orderHandler.RenewPayment)
		orderGroup.POST("/payment-status", middleware.TrackingRateLimitMiddleware(), limitLookupsByPhone, orderHandler.RefreshPaymentStatus)
		orderGroup.PATCH("/:id/cancel", limitOrdersByPhone, orderHandler.CancelOrder)
	}

	userRoutes := public.Group("/auth")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/profile", authHandler.GetProfile)
		userRoutes.PUT("/profile", limitAccountByUser, authHandler.UpdateProfile)
		userRoutes.PATCH("/password", limitAccountByUser, authHandler.ChangePassword)
		userRoutes.GET("/orders", orderHandler.GetUserOrders)
		userRoutes.GET("/orders/:id", orderHandler.GetUserOrder)
		userRoutes.POST("/orders/:id/payment-status", middleware.TrackingRateLimitMiddleware(), limitAccountByUser, orderHandler.RefreshMyPaymentStatus)
		userRoutes.GET("/catering-requests", cateringHandler.GetUserRequests)
	}

	admin := public.Group("/admin")
	admin.Use(
		middleware.AdminAuthMiddleware(adminRepo),
		// Writes sent with an Idempotency-Key aren't repeated on retry.
		middleware.Idempotency(idempotencyRepo, middleware.IdempotencyOptions{}),
		// Every successful change is recorded. Payment decisions write a
		// fuller entry in their own transaction.
		middleware.AuditWrites(auditRepo,
			"POST /api/v1/admin/orders/:id/payment-review",
			"POST /api/v1/admin/orders/:id/payment-check",
		),
	)
	{
		admin.GET("/payments", adminPaymentsHandler.ListPayments)
		admin.GET("/payments/unmatched", adminPaymentsHandler.ListUnmatchedPayments)
		admin.POST("/orders/:id/payment-review", adminPaymentsHandler.ResolvePaymentReview)
		admin.POST("/orders/:id/payment-check", adminPaymentsHandler.RecheckPayment)
		admin.GET("/jobs", adminPaymentsHandler.ListJobs)
		admin.POST("/jobs/:id/retry", adminPaymentsHandler.RetryJob)
		admin.GET("/audit", adminPaymentsHandler.ListAudit)

		admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
		admin.GET("/dashboard/stats/today", adminHandler.GetTodayStats)
		admin.GET("/reports/sales", adminHandler.GetSalesReport)
		admin.GET("/orders", adminHandler.GetAllOrders)
		admin.GET("/orders/:id", orderHandler.GetOrder)
		admin.PATCH("/orders/:id/status", adminHandler.UpdateOrderStatus)
		admin.PATCH("/orders/:id/payment", orderHandler.MarkPaymentCollected)
		admin.GET("/orders/:id/invoice", orderHandler.GetOrderInvoice)
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)
		admin.GET("/menu/items", adminHandler.GetMenuItems)
		admin.GET("/menu/items/:id", adminHandler.GetMenuItem)
		admin.POST("/menu/items", adminHandler.CreateMenuItem)
		admin.PUT("/menu/items/:id", adminHandler.UpdateMenuItem)
		admin.DELETE("/menu/items/:id", adminHandler.DeleteMenuItem)
		admin.POST("/uploads/menu-image", uploadHandler.UploadMenuImage)
		admin.POST("/uploads/team-photo", uploadHandler.UploadTeamPhoto)
		admin.GET("/team", teamHandler.ListAll)
		admin.POST("/team", teamHandler.Create)
		admin.PUT("/team/:id", teamHandler.Update)
		admin.DELETE("/team/:id", teamHandler.Delete)
		admin.POST("/menu/categories", adminHandler.CreateMenuCategory)
		admin.PUT("/menu/categories/:id", adminHandler.UpdateMenuCategory)
		admin.GET("/catering/requests", cateringHandler.GetAllRequests)
		admin.GET("/catering/requests/:id", cateringHandler.GetRequest)
		admin.PATCH("/catering/requests/:id", cateringHandler.UpdateRequest)
		admin.PATCH("/catering/requests/:id/status", cateringContractHandler.ChangeStatus)
		admin.GET("/catering/requests/:id/contract", cateringContractHandler.GetContract)
		admin.POST("/catering/requests/:id/quotes", cateringContractHandler.SendQuote)
		admin.POST("/catering/requests/:id/payments", cateringContractHandler.RecordPayment)
		admin.GET("/catering/packages", cateringHandler.ListAllPackages)
		admin.GET("/catering/packages/:id", cateringHandler.GetPackage)
		admin.POST("/catering/packages", cateringHandler.CreatePackage)
		admin.PUT("/catering/packages/:id", cateringHandler.UpdatePackage)
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
		inventoryRoutes.GET("", inventoryHandler.GetInventory)
		inventoryRoutes.GET("/dashboard", inventoryHandler.GetDashboard)
		inventoryRoutes.GET("/low-stock", inventoryHandler.GetLowStock)
		inventoryRoutes.GET("/alerts", inventoryHandler.GetAlerts)
		inventoryRoutes.GET("/check", inventoryHandler.CheckAvailability)
		inventoryRoutes.GET("/menu-item/:menuItemId", inventoryHandler.GetInventoryByMenuItem)
		inventoryRoutes.GET("/:id", inventoryHandler.GetInventoryItem)
		inventoryRoutes.POST("", inventoryHandler.CreateInventoryItem)
		inventoryRoutes.POST("/restock", inventoryHandler.RestockItem)
		inventoryRoutes.PATCH("/stock", inventoryHandler.UpdateStock)
		inventoryRoutes.PATCH("/alerts/:id/resolve", inventoryHandler.ResolveAlert)
	}

	return r, cleanup
}
