package routes

import (
	"github.com/gin-gonic/gin"
	"psclub-crm/internal/handlers"
	"psclub-crm/internal/middleware"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	companyHandler *handlers.CompanyHandler,
	clientHandler *handlers.ClientHandler,
	channelHandler *handlers.ChannelHandler,
	userHandler *handlers.UserHandler,
	expCatHandler *handlers.ExpenseCategoryHandler,
	expenseHandler *handlers.ExpenseHandler,
	tableHandler *handlers.TableHandler,
	tableCategoryHandler *handlers.TableCategoryHandler,
	bookingHandler *handlers.BookingHandler,
	categoryHandler *handlers.CategoryHandler,
	subCategoryHandler *handlers.SubcategoryHandler,
	paymentTypeHandler *handlers.PaymentTypeHandler,
	priceListHandler *handlers.PriceItemHandler,
	pricelistHistoryHandler *handlers.PricelistHistoryHandler,
	priceSetHandler *handlers.PriceSetHandler,
	equipmentHandler *handlers.EquipmentHandler,
	repairHandler *handlers.RepairHandler,
	repairCatHandler *handlers.RepairCategoryHandler,
	cashboxHandler *handlers.CashboxHandler,
	settingsHandler *handlers.SettingsHandler,
	reportHandler *handlers.ReportHandler,
	inventoryHandler *handlers.InventoryHandler,
	authSecret string,
) {
	api := r.Group("/api")

	// --- Аутентификация
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	companies := api.Group("/companies")
	{
		companies.POST("", companyHandler.Create)
	}

	api.Use(middleware.Auth(authSecret))

	// --- Клиенты
	clients := api.Group("/clients")
	{
		clients.POST("", clientHandler.CreateClient)
		clients.GET("", clientHandler.GetAllClients)
		clients.GET("/:id", clientHandler.GetClientByID)
		clients.PUT("/:id", clientHandler.UpdateClient)
		clients.DELETE("/:id", clientHandler.DeleteClient)
	}

	// --- Каналы привлечения
	channels := api.Group("/channels")
	{
		channels.POST("", channelHandler.Create)
		channels.GET("", channelHandler.GetAll)
		channels.GET("/:id", channelHandler.GetByID)
		channels.PUT("/:id", channelHandler.Update)
		channels.DELETE("/:id", channelHandler.Delete)
	}

	// --- Сотрудники (Users)
	users := api.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// --- Категории столов
	tableCategories := api.Group("/table-categories")
	{
		tableCategories.POST("", tableCategoryHandler.CreateCategory)
		tableCategories.GET("", tableCategoryHandler.GetAllCategories)
		tableCategories.GET("/:id", tableCategoryHandler.GetCategoryByID)
		tableCategories.PUT("/:id", tableCategoryHandler.UpdateCategory)
		tableCategories.DELETE("/:id", tableCategoryHandler.DeleteCategory)
	}

	// --- Столы
	tables := api.Group("/tables")
	{
		tables.POST("", tableHandler.CreateTable)
		tables.GET("", tableHandler.GetAllTables)
		tables.GET("/:id", tableHandler.GetTableByID)
		tables.PUT("/:id", tableHandler.UpdateTable)
		tables.DELETE("/:id", tableHandler.DeleteTable)
		tables.POST("/change/reorder", tableHandler.ReorderTable)
	}

	// --- Категории товаров/услуг
	categories := api.Group("/categories")
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("", categoryHandler.GetAllCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// --- Подкатегории
	subcategories := api.Group("/subcategories")
	{
		subcategories.POST("", subCategoryHandler.CreateSubcategory)
		subcategories.GET("", subCategoryHandler.GetAllSubcategories)
		subcategories.GET("/:id", subCategoryHandler.GetSubcategoryByID)
		subcategories.GET("/category/:id", subCategoryHandler.GetSubcategoryByCategoryID)
		subcategories.PUT("/:id", subCategoryHandler.UpdateSubcategory)
		subcategories.DELETE("/:id", subCategoryHandler.DeleteSubcategory)
	}

	// --- Типы оплат
	paymentTypes := api.Group("/payment-types")
	{
		paymentTypes.POST("", paymentTypeHandler.CreatePaymentType)
		paymentTypes.GET("", paymentTypeHandler.GetAllPaymentTypes)
		paymentTypes.PUT("/:id", paymentTypeHandler.UpdatePaymentType)
		paymentTypes.DELETE("/:id", paymentTypeHandler.DeletePaymentType)
	}

	// --- Прайс-лист (товары/услуги)
	pricelist := api.Group("/pricelist")
	{
		pricelist.POST("", priceListHandler.CreatePriceItem)
		pricelist.GET("", priceListHandler.GetAllPriceItems)
		pricelist.GET("/category/:id", priceListHandler.GetPriceItemsByCategoryName)
		pricelist.GET("/:id", priceListHandler.GetPriceItemByID)
		pricelist.PUT("/:id", priceListHandler.UpdatePriceItem)
		pricelist.DELETE("/:id", priceListHandler.DeletePriceItem)
		pricelist.POST("/:id/replenish", priceListHandler.Replenish)
	}

	// --- Сеты (наборы товаров)
	sets := api.Group("/sets")
	{
		sets.POST("", priceSetHandler.CreatePriceSet)
		sets.GET("", priceSetHandler.GetAllPriceSets)
		sets.GET("/:id", priceSetHandler.GetPriceSetByID)
		sets.PUT("/:id", priceSetHandler.UpdatePriceSet)
		sets.DELETE("/:id", priceSetHandler.DeletePriceSet)
	}

	// --- Оборудование
	equipment := api.Group("/equipment")
	{
		equipment.POST("", equipmentHandler.Create)
		equipment.GET("", equipmentHandler.GetAll)
		equipment.GET("/:id", equipmentHandler.GetByID)
		equipment.PUT("/:id", equipmentHandler.Update)
		equipment.DELETE("/:id", equipmentHandler.Delete)
		equipment.POST("/inventory", equipmentHandler.PerformInventory)
		equipment.GET("/inventory/history", equipmentHandler.GetHistory)
	}

	// --- История пополнений прайс-листа
	plHistory := api.Group("/pricelist-history")
	{
		plHistory.POST("", pricelistHistoryHandler.Create)
		plHistory.GET("", pricelistHistoryHandler.GetAll)
		plHistory.GET("/item/:id", pricelistHistoryHandler.GetByItem)
		plHistory.GET("/category/:id", pricelistHistoryHandler.GetByCategory)
		plHistory.DELETE("/:id", pricelistHistoryHandler.Delete)

	}

	// --- Бронирования
	bookings := api.Group("/bookings")
	{
		bookings.POST("", bookingHandler.CreateBooking)
		bookings.GET("", bookingHandler.GetAllBookings)
		bookings.GET("/client/:id", bookingHandler.GetBookingsByClientID)
		bookings.GET("/:id", bookingHandler.GetBookingByID)
		bookings.PUT("/:id", bookingHandler.UpdateBooking)
		bookings.DELETE("/:id", bookingHandler.DeleteBooking)
		// Можно добавить эндпоинт для получения позиций бронирования:
		// bookings.GET("/:id/items", bookingHandler.GetBookingItemsByBookingID)
	}

	// --- Категории расходов
	expenseCats := api.Group("/expense-categories")
	{
		expenseCats.POST("", expCatHandler.Create)
		expenseCats.GET("", expCatHandler.GetAll)
		expenseCats.PUT("/:id", expCatHandler.Update)
		expenseCats.DELETE("/:id", expCatHandler.Delete)
	}

	// --- Категории ремонтов
	repairCats := api.Group("/repair-categories")
	{
		repairCats.POST("", repairCatHandler.Create)
		repairCats.GET("", repairCatHandler.GetAll)
		repairCats.PUT("/:id", repairCatHandler.Update)
		repairCats.DELETE("/:id", repairCatHandler.Delete)
	}

	// --- Расходы
	expenses := api.Group("/expenses")
	{
		expenses.POST("", expenseHandler.CreateExpense)
		expenses.GET("", expenseHandler.GetAllExpenses)
		expenses.GET("/:id", expenseHandler.GetExpenseByID)
		expenses.PUT("/:id", expenseHandler.UpdateExpense)
		expenses.DELETE("/:id", expenseHandler.DeleteExpense)
	}

	// --- Инвентаризация
	inventory := api.Group("/inventory")
	{
		inventory.POST("", inventoryHandler.PerformInventory)
		inventory.GET("/history", inventoryHandler.GetHistory)
	}

	// --- Ремонты
	repairs := api.Group("/repairs")
	{
		repairs.POST("", repairHandler.CreateRepair)
		repairs.GET("", repairHandler.GetAllRepairs)
		repairs.GET("/:id", repairHandler.GetRepairByID)
		repairs.PUT("/:id", repairHandler.UpdateRepair)
		repairs.DELETE("/:id", repairHandler.DeleteRepair)
	}

	// --- Касса
	cashbox := api.Group("/cashbox")
	{
		cashbox.GET("", cashboxHandler.GetCashbox)
		cashbox.GET("/day", cashboxHandler.GetDay)
		cashbox.PUT("/:id", cashboxHandler.UpdateCashbox)
		cashbox.POST("/inventory", cashboxHandler.Inventory)
		cashbox.POST("/replenish", cashboxHandler.Replenish)
		cashbox.GET("/history", cashboxHandler.GetHistory)
	}

	// --- Глобальные настройки
	settings := api.Group("/settings")
	{
		settings.POST("", settingsHandler.CreateSettings)
		settings.GET("", settingsHandler.GetSettings)
		settings.GET("/tables-count", settingsHandler.GetTablesCount)
		settings.GET("/notification-time", settingsHandler.GetNotificationTime)
		settings.PUT("/:id", settingsHandler.UpdateSettings)
		settings.DELETE("/:id", settingsHandler.DeleteSettings)
	}

	// --- Отчёты (фильтрация по периодам через query-параметры)
	reports := api.Group("/reports")
	{
		reports.GET("/summary", reportHandler.GetSummaryReport)
		reports.GET("/admins", reportHandler.GetAdminsReport)
		reports.GET("/sales", reportHandler.GetSalesReport)
		reports.GET("/analytics", reportHandler.GetAnalyticsReport)
		reports.GET("/tables", reportHandler.GetTablesReport)
		reports.GET("/discounts", reportHandler.GetDiscountsReport)
	}
}
