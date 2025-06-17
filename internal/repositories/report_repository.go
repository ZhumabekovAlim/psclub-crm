package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"psclub-crm/internal/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Здесь будет только последняя актуальная версия SalesReport
func (r *ReportRepository) SalesReport(ctx context.Context, from, to time.Time) (*models.SalesReport, error) {
	query := `
        SELECT u.name,
               COUNT(DISTINCT DATE(b.start_time)) AS days,
               SUM(CASE WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%кальян%' THEN bi.quantity ELSE 0 END) AS hookahs,
               SUM(CASE WHEN pi.is_set = 1 THEN bi.quantity ELSE 0 END) AS sets,
               ROUND(u.salary_shift * COUNT(DISTINCT DATE(b.start_time)) +
                     u.salary_hookah * SUM(CASE WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%кальян%' THEN bi.quantity ELSE 0 END) +
                     u.salary_bar * SUM(CASE WHEN pi.is_set = 1 THEN bi.quantity ELSE 0 END)) AS salary
        FROM bookings b
        JOIN users u ON b.user_id = u.id
        LEFT JOIN booking_items bi ON b.id = bi.booking_id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE b.created_at BETWEEN ? AND ?
        GROUP BY u.id, u.name`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserSales
	for rows.Next() {
		var row models.UserSales
		if err := rows.Scan(&row.Name, &row.DaysWorked, &row.HookahsSold, &row.SetsSold, &row.Salary); err != nil {
			return nil, err
		}
		users = append(users, row)
	}
// --- SalesReport ---
func (r *ReportRepository) SalesReport(ctx context.Context, from, to time.Time) (*models.SalesReport, error) {
	userQuery := `
        SELECT u.name,
               COUNT(DISTINCT DATE(b.start_time)) AS days,
               SUM(CASE WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%' THEN bi.quantity ELSE 0 END) AS hookahs,
               SUM(CASE WHEN pi.is_set = 1 THEN bi.quantity ELSE 0 END) AS sets,
               ROUND(u.salary_shift * COUNT(DISTINCT DATE(b.start_time)) +
                     u.salary_hookah * SUM(CASE WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%' THEN bi.quantity ELSE 0 END) +
                     u.salary_bar * SUM(CASE WHEN pi.is_set = 1 THEN bi.quantity ELSE 0 END)) AS salary
        FROM bookings b
        JOIN users u ON b.user_id = u.id
        LEFT JOIN booking_items bi ON b.id = bi.booking_id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE b.created_at BETWEEN ? AND ?
        GROUP BY u.id, u.name`
	rows, err := r.db.QueryContext(ctx, userQuery, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.UserSales
	for rows.Next() {
		var row models.UserSales
		if err := rows.Scan(&row.Name, &row.DaysWorked, &row.HookahsSold, &row.SetsSold, &row.Salary); err != nil {
			return nil, err
		}
		users = append(users, row)
	}

	expRows, err := r.db.QueryContext(ctx, `
        SELECT title, SUM(total) FROM expenses
        WHERE date BETWEEN ? AND ?
        GROUP BY title`, from, to)
	if err != nil {
		return nil, err
	}
	defer expRows.Close()
	var expenses []models.ExpenseTotal
	var totalExp float64
	for expRows.Next() {
		var e models.ExpenseTotal
		if err := expRows.Scan(&e.Title, &e.Total); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
		totalExp += e.Total
	}

	catRows, err := r.db.QueryContext(ctx, `
        SELECT categories.name, SUM(bi.price * bi.quantity)
        FROM booking_items bi
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE bi.created_at BETWEEN ? AND ?
        GROUP BY categories.name`, from, to)
	if err != nil {
		return nil, err
	}
	defer catRows.Close()
	var incomes []models.CategoryIncome
	var totalInc float64
	for catRows.Next() {
		var inc models.CategoryIncome
		if err := catRows.Scan(&inc.Category, &inc.Total); err != nil {
			return nil, err
		}
		incomes = append(incomes, inc)
		totalInc += inc.Total
	}

	const taxPercent = 0.10
	netProfit := totalInc*(1-taxPercent) - totalExp

	return &models.SalesReport{
		Users:            users,
		Expenses:         expenses,
		IncomeByCategory: incomes,
		TotalIncome:      totalInc,
		TotalExpenses:    totalExp,
		NetProfit:        netProfit,
	}, nil
}

// --- AnalyticsReport ---
func (r *ReportRepository) AnalyticsReport(ctx context.Context, from, to time.Time) (*models.AnalyticsReport, error) {
	// Daily revenue
	dailyRows, _ := r.db.QueryContext(ctx, `
        SELECT DATE(created_at), SUM(total_amount) FROM bookings
        WHERE created_at BETWEEN ? AND ?
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)
    `, from, to)
	var daily []models.DataPoint
	for dailyRows.Next() {
		var day string
		var sum int
		dailyRows.Scan(&day, &sum)
		daily = append(daily, models.DataPoint{Label: day, Value: sum})


	}

	catRows, err := r.db.QueryContext(ctx, `
        SELECT categories.name, SUM(bi.price * bi.quantity)
        FROM booking_items bi
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE bi.created_at BETWEEN ? AND ?
        GROUP BY categories.name`, from, to)
	if err != nil {
		return nil, err
	}
	defer catRows.Close()

	var incomes []models.CategoryIncome
	var totalInc float64
	for catRows.Next() {
		var inc models.CategoryIncome
		if err := catRows.Scan(&inc.Category, &inc.Total); err != nil {
			return nil, err
		}
		incomes = append(incomes, inc)
		totalInc += inc.Total
	}

	const taxPercent = 0.10
	netProfit := totalInc*(1-taxPercent) - totalExp


	// Retrieve all orders with a discount within the period
	orderRows, err := r.db.QueryContext(ctx, `
        SELECT id, client_id, table_id, user_id, start_time, end_time, note,
               discount, discount_reason, total_amount, bonus_used,
               payment_status, payment_type_id, created_at, updated_at
        FROM bookings
        WHERE discount > 0 AND created_at BETWEEN ? AND ?
        ORDER BY created_at DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer orderRows.Close()
	var orders []models.Booking
	for orderRows.Next() {
		var b models.Booking
		if err := orderRows.Scan(&b.ID, &b.ClientID, &b.TableID, &b.UserID, &b.StartTime,
			&b.EndTime, &b.Note, &b.Discount, &b.DiscountReason, &b.TotalAmount,
			&b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, b)
	}

	return &models.DiscountsReport{
		TotalDiscount:     total,
		DiscountCount:     count,
		AvgDiscount:       avg,
		TopReasons:        reasons,
		DistributionBySum: dist,
		Orders:            orders,

	}, nil
}
