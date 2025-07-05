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
func (r *ReportRepository) SummaryReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID int) (*models.SummaryReport, error) {
	var result models.SummaryReport
	fmt.Println("SummaryReport called with:", from, to, tFrom, tTo, userID)
	query := `
        SELECT
            COALESCE(SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0) as total_revenue,
            COUNT(DISTINCT client_id) as total_clients,
            COALESCE(ROUND(AVG(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100))), 0) as avg_check
        FROM bookings b
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        WHERE DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?`
	args := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		query += " AND b.user_id = ?"
		args = append(args, userID)
	}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&result.TotalRevenue, &result.TotalClients, &result.AvgCheck,
	)
	if err != nil {
		return nil, err
	}

	// Calculate load percent
	var bookingsCount int
	countQuery := `SELECT COUNT(*) FROM bookings b WHERE DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?`
	countArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		countQuery += " AND b.user_id = ?"
		countArgs = append(countArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&bookingsCount)

	var tableCount int
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tables`).Scan(&tableCount)

	var workFrom, workTo string
	_ = r.db.QueryRowContext(ctx, `SELECT work_time_from, work_time_to FROM settings LIMIT 1`).Scan(&workFrom, &workTo)
	wf, _ := time.Parse("15:04:05", workFrom)
	wt, _ := time.Parse("15:04:05", workTo)
	hours := wt.Sub(wf).Hours()
	if hours < 0 {
		hours += 24
	}
	days := int(to.Sub(from).Hours()/24 + 0.5)
	if days <= 0 {
		days = 1
	}
	capacity := float64(tableCount) * hours * float64(days)
	if capacity > 0 {
		result.LoadPercent = float64(bookingsCount) * 100 / capacity
	}

	// Age groups
	var under18, age18to25, age26to35, age36Plus float64
	ageQuery := `
                SELECT
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) < 18 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 18 AND 25 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 26 AND 35 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) > 35 THEN 1 ELSE 0 END)
                FROM clients
                WHERE id IN (SELECT DISTINCT client_id FROM bookings WHERE DATE(created_at) BETWEEN ? AND ? AND TIME(created_at) BETWEEN ? AND ?`
	ageArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		ageQuery += " AND user_id = ?"
		ageArgs = append(ageArgs, userID)
	}
	ageQuery += ")"
	_ = r.db.QueryRowContext(ctx, ageQuery, ageArgs...).Scan(&under18, &age18to25, &age26to35, &age36Plus)
	totalClients := result.TotalClients
	if totalClients > 0 {
		result.AgeUnder18 = float64(under18) * 100 / float64(totalClients)
		result.Age18To25 = float64(age18to25) * 100 / float64(totalClients)
		result.Age26To35 = float64(age26to35) * 100 / float64(totalClients)
		result.Age36Plus = float64(age36Plus) * 100 / float64(totalClients)
	}

	// Category sales
	catQuery := `
                SELECT categories.name, SUM(booking_items.price * (1 - IFNULL(pt.hold_percent,0)/100))
                FROM booking_items
                LEFT JOIN bookings ON booking_items.booking_id = bookings.id
                LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
                LEFT JOIN price_items ON booking_items.item_id = price_items.id
                LEFT JOIN categories ON price_items.category_id = categories.id
                WHERE DATE(booking_items.created_at) BETWEEN ? AND ? AND TIME(booking_items.created_at) BETWEEN ? AND ?`
	catArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		catQuery += " AND bookings.user_id = ?"
		catArgs = append(catArgs, userID)
	}
	catQuery += " GROUP BY categories.name"
	catRows, _ := r.db.QueryContext(ctx, catQuery, catArgs...)
	var catSales []models.CategorySale
	for catRows.Next() {
		var name string
		var revenue int
		catRows.Scan(&name, &revenue)
		catSales = append(catSales, models.CategorySale{Category: name, Revenue: revenue})
	}
	result.CategorySales = catSales

	// Top items by profit
	itemQuery := `
    SELECT 
        price_items.name,
        SUM(booking_items.quantity),
        SUM(
            CASE
                WHEN categories.name = 'Часы' THEN (booking_items.price - booking_items.discount) * (1 - IFNULL(pt.hold_percent,0)/100)
                ELSE (booking_items.price - booking_items.discount) * (1 - IFNULL(pt.hold_percent,0)/100)
            END
        ),
        SUM(
            CASE 
                WHEN categories.name = 'Часы' THEN price_items.buy_price
                ELSE price_items.buy_price
            END
        ),
        SUM(
            CASE
                WHEN categories.name = 'Часы' THEN (booking_items.price - booking_items.discount)*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
                ELSE (booking_items.price - booking_items.discount)*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
            END
        )
    FROM booking_items
    LEFT JOIN bookings ON booking_items.booking_id = bookings.id
    LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
    LEFT JOIN price_items ON booking_items.item_id = price_items.id
    LEFT JOIN categories ON price_items.category_id = categories.id
    WHERE DATE(booking_items.created_at) BETWEEN ? AND ?
      AND TIME(booking_items.created_at) BETWEEN ? AND ?`

	itemArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		itemQuery += " AND bookings.user_id = ?"
		itemArgs = append(itemArgs, userID)
	}

	itemQuery += `
    GROUP BY price_items.name
    ORDER BY SUM(
        CASE
            WHEN categories.name = 'Часы' THEN (booking_items.price - booking_items.discount)*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
            ELSE (booking_items.price - booking_items.discount)*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
        END
    ) DESC
    LIMIT 5`

	itemRows, err := r.db.QueryContext(ctx, itemQuery, itemArgs...)
	if err != nil {

	}
	defer itemRows.Close()

	var topItems []models.ProfitItem
	for itemRows.Next() {
		var it models.ProfitItem
		if err := itemRows.Scan(&it.Name, &it.Quantity, &it.Revenue, &it.Expense, &it.Profit); err != nil {

		}
		topItems = append(topItems, it)
	}
	result.TopItems = topItems

	var prevRevenue, prevClients, prevAvgCheck float64
	prevFrom := from.Add(-(to.Sub(from)))
	prevTo := from
	prevQuery := `
        SELECT COALESCE(SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0), COUNT(DISTINCT client_id), COALESCE(AVG(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0)
        FROM bookings b
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        WHERE DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?`
	prevArgs := []interface{}{prevFrom, prevTo, tFrom, tTo}
	if userID > 0 {
		prevQuery += " AND b.user_id = ?"
		prevArgs = append(prevArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, prevQuery, prevArgs...).Scan(&prevRevenue, &prevClients, &prevAvgCheck)
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
func (r *ReportRepository) AdminsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID int) (*models.AdminsReport, error) {

	query := `
       SELECT u.id, u.name,
              COUNT(DISTINCT DATE(b.start_time)) AS shifts,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%кальян%' THEN bi.quantity
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%кальян%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS hookah_qty,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%кальян%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS set_qty,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%кальян%' THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%кальян%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS hookah_rev,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%кальян%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS set_rev,
              u.salary_shift, u.salary_hookah, u.hookah_salary_type, u.salary_bar
        FROM users u
       LEFT JOIN bookings b ON b.user_id = u.id
            AND DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       LEFT JOIN booking_items bi ON b.id = bi.booking_id
       LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE u.role = 'admin'`
	args := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		query += " AND u.id = ?"
		args = append(args, userID)
	}
	query += `
       GROUP BY u.id, u.name, u.salary_shift, u.salary_hookah, u.hookah_salary_type, u.salary_bar
       ORDER BY u.name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report models.AdminsReport
	for rows.Next() {
		var (
			id                    int
			name                  string
			shifts, hookahs, sets int
			hookahRev, setRev     float64
			shiftSalary           int
			hookahValue           float64
			hookahType            string
			setPercent            float64
		)
		if err := rows.Scan(&id, &name, &shifts, &hookahs, &sets, &hookahRev, &setRev, &shiftSalary, &hookahValue, &hookahType, &setPercent); err != nil {
			return nil, err
		}

		shiftTotal := shifts * shiftSalary
		var hookahTotal int
		if hookahType == "percent" {
			hookahTotal = int(hookahRev * hookahValue / 100)
		} else {
			hookahTotal = int(float64(hookahs) * hookahValue)
		}
		setTotal := int(setRev * setPercent / 100)
		salary := shiftTotal + hookahTotal + setTotal

		var hookahDetail string
		if hookahType == "percent" {
			hookahDetail = fmt.Sprintf("%d₸ × %.0f%% = %d₸", int(hookahRev), hookahValue, hookahTotal)
		} else {
			hookahDetail = fmt.Sprintf("%d × %.0f₸ = %d₸", hookahs, hookahValue, hookahTotal)
		}
		detail := fmt.Sprintf("Смены: %d × %d = %d₸, кальяны: %s, сеты: %d₸ × %.0f%% = %d₸",
			shifts, shiftSalary, shiftTotal,
			hookahDetail,
			int(setRev), setPercent, setTotal)

		report.Admins = append(report.Admins, models.AdminReportRow{
			Name:         name,
			Shifts:       shifts,
			HookahsSold:  hookahs,
			SetsSold:     sets,
			Salary:       salary,
			SalaryDetail: detail,
		})
	}

	return &report, nil

}

// --- SalesReport ---
func (r *ReportRepository) SalesReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID int) (*models.SalesReport, error) {
	userQuery := `
       SELECT u.id, u.name,
              COUNT(DISTINCT DATE(b.start_time)) AS days,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%' THEN bi.quantity
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS hookahs,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS sets,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%' THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS hookah_rev,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%\u043a\u0430\u043b\u044c\u044f\u043d%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS set_rev,
              u.salary_shift, u.salary_hookah, u.hookah_salary_type, u.salary_bar
       FROM bookings b
       JOIN users u ON b.user_id = u.id
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       LEFT JOIN booking_items bi ON b.id = bi.booking_id
       LEFT JOIN price_items pi ON bi.item_id = pi.id
       LEFT JOIN categories ON pi.category_id = categories.id
       WHERE DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?`
	userArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		userQuery += " AND b.user_id = ?"
		userArgs = append(userArgs, userID)
	}
	userQuery += `
       GROUP BY u.id, u.name, u.salary_shift, u.salary_hookah, u.hookah_salary_type, u.salary_bar`
	rows, err := r.db.QueryContext(ctx, userQuery, userArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.UserSales
	for rows.Next() {
		var (
			id          int
			name        string
			days        int
			hookahs     int
			sets        int
			hookahRev   float64
			setRev      float64
			shiftSalary int
			hookahValue float64
			hookahType  string
			setPercent  float64
		)
		if err := rows.Scan(&id, &name, &days, &hookahs, &sets, &hookahRev, &setRev, &shiftSalary, &hookahValue, &hookahType, &setPercent); err != nil {
			return nil, err
		}
		shiftTotal := days * shiftSalary
		var hookahTotal int
		if hookahType == "percent" {
			hookahTotal = int(hookahRev * hookahValue / 100)
		} else {
			hookahTotal = int(float64(hookahs) * hookahValue)
		}
		setTotal := int(setRev * setPercent / 100)
		users = append(users, models.UserSales{
			Name:        name,
			DaysWorked:  days,
			HookahsSold: hookahs,
			SetsSold:    sets,
			Salary:      float64(shiftTotal + hookahTotal + setTotal),
		})
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

	catQuery2 := `
        SELECT categories.name, SUM((bi.price-bi.discount) * (1 - IFNULL(pt.hold_percent,0)/100))
        FROM booking_items bi
        LEFT JOIN bookings b ON bi.booking_id = b.id
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE DATE(bi.created_at) BETWEEN ? AND ? AND TIME(bi.created_at) BETWEEN ? AND ?`
	catArgs2 := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		catQuery2 += " AND b.user_id = ?"
		catArgs2 = append(catArgs2, userID)
	}
	catQuery2 += `
        GROUP BY categories.name`
	catRows, err := r.db.QueryContext(ctx, catQuery2, catArgs2...)
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
func (r *ReportRepository) AnalyticsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID int) (*models.AnalyticsReport, error) {
	// Daily revenue
	dailyQuery := `
        SELECT DATE(b.start_time), SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)) FROM bookings b
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        WHERE DATE(b.start_time) BETWEEN ? AND ? AND TIME(b.start_time) BETWEEN ? AND ?`
	dailyArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		dailyQuery += " AND b.user_id = ?"
		dailyArgs = append(dailyArgs, userID)
	}
	dailyQuery += `
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)`
	dailyRows, _ := r.db.QueryContext(ctx, dailyQuery, dailyArgs...)
	var daily []models.DataPoint
	for dailyRows.Next() {
		var day string
		var sum int
		dailyRows.Scan(&day, &sum)
		daily = append(daily, models.DataPoint{Label: day, Value: sum})
	}

	// Hourly load
	hourlyQuery := `
        SELECT HOUR(start_time), COUNT(*) FROM bookings
        WHERE DATE(start_time) BETWEEN ? AND ? AND TIME(start_time) BETWEEN ? AND ?`
	hourlyArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		hourlyQuery += " AND user_id = ?"
		hourlyArgs = append(hourlyArgs, userID)
	}
	hourlyQuery += `
        GROUP BY HOUR(start_time)
        ORDER BY HOUR(start_time)`
	hourlyRows, _ := r.db.QueryContext(ctx, hourlyQuery, hourlyArgs...)
	var hourly []models.DataPoint
	for hourlyRows.Next() {
		var hour int
		var count int
		hourlyRows.Scan(&hour, &count)
		hourly = append(hourly, models.DataPoint{Label: fmt.Sprintf("%02d:00", hour), Value: count})
	}

	// Category stats
	catQuery := `
        SELECT categories.name, SUM(booking_items.quantity), SUM(booking_items.price * (1 - IFNULL(pt.hold_percent,0)/100))
        FROM booking_items
        LEFT JOIN bookings ON booking_items.booking_id = bookings.id
        LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
        LEFT JOIN price_items ON booking_items.item_id = price_items.id
        LEFT JOIN categories ON price_items.category_id = categories.id
        WHERE DATE(booking_items.created_at) BETWEEN ? AND ? AND TIME(booking_items.created_at) BETWEEN ? AND ?`
	catArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		catQuery += " AND bookings.user_id = ?"
		catArgs = append(catArgs, userID)
	}
	catQuery += `
        GROUP BY categories.name`
	catRows, _ := r.db.QueryContext(ctx, catQuery, catArgs...)
	var cats []models.CategoryStat
	for catRows.Next() {
		var name string
		var qty, revenue float64
		catRows.Scan(&name, &qty, &revenue)
		avgCheck := 0.0
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
func (r *ReportRepository) DiscountsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID int) (*models.DiscountsReport, error) {
	var total, count, avg int
	sumQuery := `
        SELECT COALESCE(SUM(discount),0), COUNT(*), COALESCE(AVG(discount),0)
        FROM bookings
        WHERE discount > 0 AND DATE(start_time) BETWEEN ? AND ? AND TIME(start_time) BETWEEN ? AND ?`
	sumArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		sumQuery += " AND user_id = ?"
		sumArgs = append(sumArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, sumQuery, sumArgs...).Scan(&total, &count, &avg)

	reasonQuery := `
        SELECT discount_reason, COUNT(*), SUM(discount), COALESCE(AVG(discount),0)
        FROM bookings
        WHERE discount > 0 AND DATE(start_time) BETWEEN ? AND ? AND TIME(start_time) BETWEEN ? AND ?`
	reasonArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		reasonQuery += " AND user_id = ?"
		reasonArgs = append(reasonArgs, userID)
	}
	reasonQuery += `
        GROUP BY discount_reason ORDER BY SUM(discount) DESC LIMIT 5`
	rows, _ := r.db.QueryContext(ctx, reasonQuery, reasonArgs...)
	var reasons []models.ReasonRow
	for rows.Next() {
		var reason sql.NullString
		var cnt, sum, av float64
		rows.Scan(&reason, &cnt, &sum, &av)
		reasons = append(reasons, models.ReasonRow{
			Reason: reason.String, Count: int(cnt), Sum: sum, Avg: av,
		})
	}

	distQuery := `
        SELECT discount, COUNT(*) FROM bookings
        WHERE discount > 0 AND DATE(start_time) BETWEEN ? AND ? AND TIME(start_time) BETWEEN ? AND ?`
	distArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		distQuery += " AND user_id = ?"
		distArgs = append(distArgs, userID)
	}
	distQuery += `
        GROUP BY discount
        ORDER BY discount DESC`
	distRows, _ := r.db.QueryContext(ctx, distQuery, distArgs...)
	var dist []models.DataPoint
	for distRows.Next() {
		var sum, cnt int
		distRows.Scan(&sum, &cnt)
		dist = append(dist, models.DataPoint{
			Label: fmt.Sprintf("%d₸", sum), Value: cnt,
		})
	}

	// Retrieve all orders with a discount within the period
	orderQuery := `
        SELECT id, client_id, table_id, user_id, start_time, end_time, note,
               discount, discount_reason, total_amount, bonus_used,
               payment_status, payment_type_id, created_at, updated_at
        FROM bookings
        WHERE discount > 0 AND DATE(start_time) BETWEEN ? AND ? AND TIME(start_time) BETWEEN ? AND ?`
	orderArgs := []interface{}{from, to, tFrom, tTo}
	if userID > 0 {
		orderQuery += " AND user_id = ?"
		orderArgs = append(orderArgs, userID)
	}
	orderQuery += `
        ORDER BY created_at DESC`
	orderRows, err := r.db.QueryContext(ctx, orderQuery, orderArgs...)
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
