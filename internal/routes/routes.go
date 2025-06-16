package routes

import (
	"github.com/gin-gonic/gin"
	"psclub-crm/internal/handlers"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	clientHandler *handlers.ClientHandler,
	userHandler *handlers.UserHandler,
	expenseHandler *handlers.ExpenseHandler,
	tableHandler *handlers.TableHandler,
	tableCategoryHandler *handlers.TableCategoryHandler,
	bookingHandler *handlers.BookingHandler,
	categoryHandler *handlers.CategoryHandler,
	subCategoryHandler *handlers.SubcategoryHandler,
	priceListHandler *handlers.PriceItemHandler,
	repairHandler *handlers.RepairHandler,
	cashboxHandler *handlers.CashboxHandler,
	settingsHandler *handlers.SettingsHandler,
	reportHandler *handlers.ReportHandler,
) {
	api := r.Group("/api")

	// --- Аутентификация
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// --- Клиенты
	clients := api.Group("/clients")
	{
		clients.POST("", clientHandler.CreateClient)
		clients.GET("", clientHandler.GetAllClients)
		clients.GET("/:id", clientHandler.GetClientByID)
		clients.PUT("/:id", clientHandler.UpdateClient)
		clients.DELETE("/:id", clientHandler.DeleteClient)
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
		subcategories.PUT("/:id", subCategoryHandler.UpdateSubcategory)
		subcategories.DELETE("/:id", subCategoryHandler.DeleteSubcategory)
	}

	// --- Прайс-лист (товары/услуги)
	pricelist := api.Group("/pricelist")
	{
		pricelist.POST("", priceListHandler.CreatePriceItem)
		pricelist.GET("", priceListHandler.GetAllPriceItems)
		pricelist.GET("/:id", priceListHandler.GetPriceItemByID)
		pricelist.PUT("/:id", priceListHandler.UpdatePriceItem)
		pricelist.DELETE("/:id", priceListHandler.DeletePriceItem)
	}

	//// --- История закупа (прихода на склад)
	//priceItemHistory := api.Group("/price-item-history")
	//{
	//	priceItemHistory.POST("", priceItemHistoryHandler.CreateHistory)
	//	priceItemHistory.GET("", priceItemHistoryHandler.GetAllHistories)
	//	priceItemHistory.GET("/:id", priceItemHistoryHandler.GetHistoryByID)
	//	priceItemHistory.DELETE("/:id", priceItemHistoryHandler.DeleteHistory)
	//}

	// --- Бронирования
	bookings := api.Group("/bookings")
	{
		bookings.POST("", bookingHandler.CreateBooking)
		bookings.GET("", bookingHandler.GetAllBookings)
		bookings.GET("/:id", bookingHandler.GetBookingByID)
		bookings.PUT("/:id", bookingHandler.UpdateBooking)
		bookings.DELETE("/:id", bookingHandler.DeleteBooking)
		// Можно добавить эндпоинт для получения позиций бронирования:
		// bookings.GET("/:id/items", bookingHandler.GetBookingItemsByBookingID)
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
		cashbox.PUT("/:id", cashboxHandler.UpdateCashbox)
	}

	// --- Глобальные настройки
	settings := api.Group("/settings")
	{
		settings.GET("", settingsHandler.GetSettings)
		settings.PUT("/:id", settingsHandler.UpdateSettings)
	}

	// --- Отчёты (фильтрация по периодам через query-параметры)
	reports := api.Group("/reports")
	{
		reports.GET("/summary", reportHandler.GetSummaryReport)
		reports.GET("/admins", reportHandler.GetAdminsReport)
		reports.GET("/sales", reportHandler.GetSalesReport)
		reports.GET("/analytics", reportHandler.GetAnalyticsReport)
		reports.GET("/discounts", reportHandler.GetDiscountsReport)
	}
}
