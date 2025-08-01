package app

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"psclub-crm/internal/config"
	"psclub-crm/internal/handlers"
	"psclub-crm/internal/middleware"
	"psclub-crm/internal/repositories"
	"psclub-crm/internal/routes"
	"psclub-crm/internal/services"
	"time"
)

func Run() {
	cfg := config.LoadConfig()
	dsn := cfg.Database.DSN
	port := cfg.Server.Port

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer db.Close()
	// ========== Инициализация зависимостей ==========

	// Клиенты
	clientRepo := repositories.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	// Каналы привлечения
	channelRepo := repositories.NewChannelRepository(db)
	channelService := services.NewChannelService(channelRepo)
	channelHandler := handlers.NewChannelHandler(channelService)

	// Сотрудники (Users)
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	authService := services.NewAuthService(
		userRepo,
		tokenRepo,
		cfg.Auth.AccessSecret,
		cfg.Auth.RefreshSecret,
		time.Duration(cfg.Auth.AccessTTL)*time.Second,
		time.Duration(cfg.Auth.RefreshTTL)*time.Second,
	)
	authHandler := handlers.NewAuthHandler(authService)

	// Категории столов
	tableCategoryRepo := repositories.NewTableCategoryRepository(db)
	tableCategoryService := services.NewTableCategoryService(tableCategoryRepo)
	tableCategoryHandler := handlers.NewTableCategoryHandler(tableCategoryService)

	// Столы
	tableRepo := repositories.NewTableRepository(db)
	tableService := services.NewTableService(tableRepo)
	tableHandler := handlers.NewTableHandler(tableService)

	// Категории товаров/услуг
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Подкатегории
	subCategoryRepo := repositories.NewSubcategoryRepository(db)
	subCategoryService := services.NewSubcategoryService(subCategoryRepo)
	subCategoryHandler := handlers.NewSubcategoryHandler(subCategoryService)

	// Payment types
	paymentTypeRepo := repositories.NewPaymentTypeRepository(db)
	paymentTypeService := services.NewPaymentTypeService(paymentTypeRepo)
	paymentTypeHandler := handlers.NewPaymentTypeHandler(paymentTypeService)

	// Прайс-лист
	priceRepo := repositories.NewPriceItemRepository(db)
	historyRepo := repositories.NewPriceItemHistoryRepository(db)
	plHistoryRepo := repositories.NewPricelistHistoryRepository(db)
	priceService := services.NewPriceItemService(priceRepo, historyRepo, plHistoryRepo)

	// Сеты товаров
	priceSetRepo := repositories.NewPriceSetRepository(db)
	priceSetService := services.NewPriceSetService(priceSetRepo, priceRepo, categoryRepo)
	priceSetHandler := handlers.NewPriceSetHandler(priceSetService)

	// Бронирования (инициализируем позже, после кассы)
	bookingItemRepo := repositories.NewBookingItemRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	bookingPaymentRepo := repositories.NewBookingPaymentRepository(db)

	// Категории расходов и сами расходы
	expCatRepo := repositories.NewExpenseCategoryRepository(db)
	expCatService := services.NewExpenseCategoryService(expCatRepo)
	expCatHandler := handlers.NewExpenseCategoryHandler(expCatService)

	// Категории ремонтов
	repCatRepo := repositories.NewRepairCategoryRepository(db)
	repCatService := services.NewRepairCategoryService(repCatRepo)
	repCatHandler := handlers.NewRepairCategoryHandler(repCatService)

	expenseRepo := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	// Инвентаризация
	invHistRepo := repositories.NewInventoryHistoryRepository(db)
	inventoryService := services.NewInventoryService(priceRepo, invHistRepo, expenseService, expCatService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Оборудование
	equipmentRepo := repositories.NewEquipmentRepository(db)
	equipmentInvHistRepo := repositories.NewEquipmentInventoryHistoryRepository(db)
	equipmentService := services.NewEquipmentService(equipmentRepo)
	equipmentInvService := services.NewEquipmentInventoryService(equipmentRepo, equipmentInvHistRepo, expenseService, expCatService)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService, equipmentInvService)

	// Прайс-лист handlers depend on expense and category services
	priceHandler := handlers.NewPriceItemHandler(priceService, expenseService, expCatService, categoryService)
	plHistoryHandler := handlers.NewPricelistHistoryHandler(priceService, expenseService)

	// Ремонты
	repairRepo := repositories.NewRepairRepository(db)
	repairService := services.NewRepairService(repairRepo)
	repairHandler := handlers.NewRepairHandler(repairService, expenseService, expCatService)

	// Касса
	cashboxRepo := repositories.NewCashboxRepository(db)
	cashboxHistRepo := repositories.NewCashboxHistoryRepository(db)
	cashboxService := services.NewCashboxService(cashboxRepo, cashboxHistRepo, expenseService, expCatService, settingsRepo)
	cashboxHandler := handlers.NewCashboxHandlerCashboxHandler(cashboxService)

	bookingService := services.NewBookingService(
		bookingRepo,
		bookingItemRepo,
		clientRepo,
		settingsRepo,
		priceRepo,
		priceSetRepo,
		categoryRepo,
		bookingPaymentRepo,
		paymentTypeRepo,
		cashboxService,
	)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Настройки
	settingsRepo = repositories.NewSettingsRepository(db)
	settingsService := services.NewSettingsService(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsService)

	// Отчеты
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// ========== Роутер и middlewares ==========
	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	routes.SetupRoutes(
		router,
		authHandler,
		clientHandler,
		channelHandler,
		userHandler,
		expCatHandler,
		expenseHandler,
		tableHandler,
		tableCategoryHandler,
		bookingHandler,
		categoryHandler,
		subCategoryHandler,
		paymentTypeHandler,
		priceHandler,
		plHistoryHandler,
		priceSetHandler,
		equipmentHandler,
		repairHandler,
		repCatHandler,
		cashboxHandler,
		settingsHandler,
		reportHandler,
		inventoryHandler,
	)

	listenAddr := fmt.Sprintf(":%d", port)
	log.Printf("Server started on %s", listenAddr)
	if err := router.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
