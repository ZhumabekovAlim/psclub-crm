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

	// Прайс-лист
	priceRepo := repositories.NewPriceItemRepository(db)
	historyRepo := repositories.NewPriceItemHistoryRepository(db)
	priceService := services.NewPriceItemService(priceRepo, historyRepo)
	priceHandler := handlers.NewPriceItemHandler(priceService)

	// Сеты товаров
	priceSetRepo := repositories.NewPriceSetRepository(db)
	priceSetService := services.NewPriceSetService(priceSetRepo)
	priceSetHandler := handlers.NewPriceSetHandler(priceSetService)

	// Бронирования
	bookingItemRepo := repositories.NewBookingItemRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	bookingService := services.NewBookingService(
		bookingRepo,
		bookingItemRepo,
		clientRepo,
		settingsRepo,
	)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Расходы
	expenseRepo := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	// Ремонты
	repairRepo := repositories.NewRepairRepository(db)
	repairService := services.NewRepairService(repairRepo)
	repairHandler := handlers.NewRepairHandler(repairService)

	// Касса
	cashboxRepo := repositories.NewCashboxRepository(db)
	cashboxService := services.NewCashboxService(cashboxRepo)
	cashboxHandler := handlers.NewCashboxHandlerCashboxHandler(cashboxService)

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
		userHandler,
		expenseHandler,
		tableHandler,
		tableCategoryHandler,
		bookingHandler,
		categoryHandler,
		subCategoryHandler,
		priceHandler,
		priceSetHandler,
		repairHandler,
		cashboxHandler,
		settingsHandler,
		reportHandler,
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
