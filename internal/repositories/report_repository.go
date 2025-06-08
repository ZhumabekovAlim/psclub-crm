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

// --- SummaryReport ---
func (r *ReportRepository) SummaryReport(ctx context.Context, from, to time.Time) (*models.SummaryReport, error) {
	var result models.SummaryReport

	err := r.db.QueryRowContext(ctx, `
        SELECT 
            COALESCE(SUM(total_amount),0) as total_revenue,
            COUNT(DISTINCT client_id) as total_clients,
            COALESCE(AVG(total_amount),0) as avg_check,
            68 as load_percent
        FROM bookings
        WHERE created_at BETWEEN ? AND ?
    `, from, to).Scan(
		&result.TotalRevenue, &result.TotalClients, &result.AvgCheck, &result.LoadPercent,
	)
	if err != nil {
		return nil, err
	}

	var prevRevenue, prevClients, prevAvgCheck int
	prevFrom := from.Add(-(to.Sub(from)))
	prevTo := from
	_ = r.db.QueryRowContext(ctx, `
        SELECT COALESCE(SUM(total_amount),0), COUNT(DISTINCT client_id), COALESCE(AVG(total_amount),0)
        FROM bookings WHERE created_at BETWEEN ? AND ?
    `, prevFrom, prevTo).Scan(&prevRevenue, &prevClients, &prevAvgCheck)
	if prevRevenue > 0 {
		result.RevenueChange = float64(result.TotalRevenue-prevRevenue) * 100.0 / float64(prevRevenue)
	}
	if prevClients > 0 {
		result.ClientsChange = float64(result.TotalClients-prevClients) * 100.0 / float64(prevClients)
	}
	if prevAvgCheck > 0 {
		result.AvgCheckChange = float64(result.AvgCheck-prevAvgCheck) * 100.0 / float64(prevAvgCheck)
	}
	result.LoadChange = 3 // Пример
	return &result, nil
}

// --- AdminsReport ---
func (r *ReportRepository) AdminsReport(ctx context.Context, from, to time.Time) (*models.AdminsReport, error) {
	// Пример: возвращаем двух админов, статично
	report := &models.AdminsReport{
		Admins: []models.AdminReportRow{
			{Name: "Иван Петров", Shifts: 15, HookahsSold: 25, SetsSold: 12, Salary: 32500, SalaryDetail: "15 × 1000 = 15000₸, кальяны: 2500₸, сеты: 1800₸"},
			{Name: "Мария Сидорова", Shifts: 12, HookahsSold: 18, SetsSold: 8, Salary: 26400, SalaryDetail: "12 × 1000 = 12000₸, кальяны: 1800₸, сеты: 1200₸"},
		},
	}
	return report, nil
}

// --- SalesReport ---
func (r *ReportRepository) SalesReport(ctx context.Context, from, to time.Time) (*models.SalesReport, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT name, SUM(quantity) as qty, SUM(quantity * price) as revenue
        FROM booking_items
        LEFT JOIN price_items ON booking_items.item_id = price_items.id
        WHERE booking_items.created_at BETWEEN ? AND ?
        GROUP BY name ORDER BY revenue DESC LIMIT 10
    `, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SaleItem
	for rows.Next() {
		var s models.SaleItem
		rows.Scan(&s.Name, &s.Quantity, &s.Revenue)
		if s.Quantity > 0 {
			s.AvgCheck = s.Revenue / s.Quantity
		}
		items = append(items, s)
	}
	return &models.SalesReport{TopSales: items}, nil
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

	// Hourly load
	hourlyRows, _ := r.db.QueryContext(ctx, `
        SELECT HOUR(start_time), COUNT(*) FROM bookings
        WHERE created_at BETWEEN ? AND ?
        GROUP BY HOUR(start_time)
        ORDER BY HOUR(start_time)
    `, from, to)
	var hourly []models.DataPoint
	for hourlyRows.Next() {
		var hour int
		var count int
		hourlyRows.Scan(&hour, &count)
		hourly = append(hourly, models.DataPoint{Label: fmt.Sprintf("%02d:00", hour), Value: count})
	}

	// Category stats
	catRows, _ := r.db.QueryContext(ctx, `
        SELECT categories.name, SUM(booking_items.quantity), SUM(booking_items.price * booking_items.quantity)
        FROM booking_items
        LEFT JOIN price_items ON booking_items.item_id = price_items.id
        LEFT JOIN categories ON price_items.category_id = categories.id
        WHERE booking_items.created_at BETWEEN ? AND ?
        GROUP BY categories.name
    `, from, to)
	var cats []models.CategoryStat
	for catRows.Next() {
		var name string
		var qty, revenue int
		catRows.Scan(&name, &qty, &revenue)
		avgCheck := 0
		if qty > 0 {
			avgCheck = revenue / qty
		}
		cats = append(cats, models.CategoryStat{
			Category: name, Quantity: qty, Revenue: revenue, AvgCheck: avgCheck,
		})
	}
	return &models.AnalyticsReport{
		DailyRevenue:  daily,
		HourlyLoad:    hourly,
		CategoryStats: cats,
	}, nil
}

// --- DiscountsReport ---
func (r *ReportRepository) DiscountsReport(ctx context.Context, from, to time.Time) (*models.DiscountsReport, error) {
	var total, count, avg int
	_ = r.db.QueryRowContext(ctx, `
        SELECT COALESCE(SUM(discount),0), COUNT(*), COALESCE(AVG(discount),0)
        FROM bookings
        WHERE discount > 0 AND created_at BETWEEN ? AND ?
    `, from, to).Scan(&total, &count, &avg)

	rows, _ := r.db.QueryContext(ctx, `
        SELECT discount_reason, COUNT(*), SUM(discount), COALESCE(AVG(discount),0)
        FROM bookings
        WHERE discount > 0 AND created_at BETWEEN ? AND ?
        GROUP BY discount_reason ORDER BY SUM(discount) DESC LIMIT 5
    `, from, to)
	var reasons []models.ReasonRow
	for rows.Next() {
		var reason sql.NullString
		var cnt, sum, av int
		rows.Scan(&reason, &cnt, &sum, &av)
		reasons = append(reasons, models.ReasonRow{
			Reason: reason.String, Count: cnt, Sum: sum, Avg: av,
		})
	}

	distRows, _ := r.db.QueryContext(ctx, `
        SELECT discount, COUNT(*) FROM bookings
        WHERE discount > 0 AND created_at BETWEEN ? AND ?
        GROUP BY discount
        ORDER BY discount DESC
    `, from, to)
	var dist []models.DataPoint
	for distRows.Next() {
		var sum, cnt int
		distRows.Scan(&sum, &cnt)
		dist = append(dist, models.DataPoint{
			Label: fmt.Sprintf("%d₸", sum), Value: cnt,
		})
	}

	return &models.DiscountsReport{
		TotalDiscount:     total,
		DiscountCount:     count,
		AvgDiscount:       avg,
		TopReasons:        reasons,
		DistributionBySum: dist,
	}, nil
}
