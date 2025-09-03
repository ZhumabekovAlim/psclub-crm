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

func buildTimeCondition(field string, from, to time.Time, tFrom, tTo string) (string, []interface{}) {
	if tFrom <= tTo {
		cond := fmt.Sprintf("DATE(%s) BETWEEN ? AND ? AND TIME(%s) BETWEEN ? AND ?", field, field)
		return cond, []interface{}{from, to, tFrom, tTo}
	}
	fromNext := from.AddDate(0, 0, 1)
	toNext := to.AddDate(0, 0, 1)
	cond := fmt.Sprintf("((DATE(%s) BETWEEN ? AND ? AND TIME(%s) >= ?) OR (DATE(%s) BETWEEN ? AND ? AND TIME(%s) <= ?))", field, field, field, field)
	args := []interface{}{from, to, tFrom, fromNext, toNext, tTo}
	return cond, args
}

// --- SummaryReport ---
func (r *ReportRepository) SummaryReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) (*models.SummaryReport, error) {
	var result models.SummaryReport
	fmt.Println("SummaryReport called with:", from, to, tFrom, tTo, userID)
	cond, condArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	cond = "b.company_id=? AND b.branch_id=? AND " + cond + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	query := fmt.Sprintf(`
       SELECT
           COALESCE(SUM(b.total_amount), 0) as total,
           COALESCE(SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0) as total_revenue,
           COUNT(DISTINCT client_id) + SUM(CASE WHEN client_id IS NULL THEN 1 ELSE 0 END) as total_clients,
           COALESCE(ROUND(AVG(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100))), 0) as avg_check
       FROM bookings b
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       WHERE %s`, cond)
	args := append([]interface{}{companyID, branchID}, condArgs...)
	if userID > 0 {
		query += " AND b.user_id = ?"
		args = append(args, userID)
	}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&result.Total, &result.TotalRevenue, &result.TotalClients, &result.AvgCheck,
	)
	if err != nil {
		return nil, err
	}

	// Calculate total cost for Bar and Hookah categories
	condCost, costArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condCost = "b.company_id=? AND b.branch_id=? AND " + condCost + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	costQuery := fmt.Sprintf(`
        SELECT COALESCE(SUM(
            bi.price * (1 - bi.discount / 100) * (1 - IFNULL(pt.hold_percent,0)/100) -
            bi.quantity * pi.buy_price
        ),0)
        FROM booking_items bi
        LEFT JOIN bookings b ON bi.booking_id = b.id
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories c ON pi.category_id = c.id
       WHERE %s AND (c.name = 'Бар' OR c.name = 'Кальян')`, condCost)
	costArgs = append([]interface{}{companyID, branchID}, costArgs...)
	if userID > 0 {
		costQuery += " AND b.user_id = ?"
		costArgs = append(costArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, costQuery, costArgs...).Scan(&result.TotalCost)

	// Calculate load percent
	var bookingsCount int
	condCount, countArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condCount = "b.company_id=? AND b.branch_id=? AND " + condCount + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM bookings b WHERE %s", condCount)
	countArgs = append([]interface{}{companyID, branchID}, countArgs...)
	if userID > 0 {
		countQuery += " AND b.user_id = ?"
		countArgs = append(countArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&bookingsCount)

	var tableCount int
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tables WHERE company_id=? AND branch_id=?`, companyID, branchID).Scan(&tableCount)

	var workFrom, workTo string
	_ = r.db.QueryRowContext(ctx, `SELECT work_time_from, work_time_to FROM settings WHERE company_id=? AND branch_id=? LIMIT 1`, companyID, branchID).Scan(&workFrom, &workTo)
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
	condAge, ageArgs := buildTimeCondition("created_at", from, to, tFrom, tTo)
	condAge = "company_id=? AND branch_id=? AND " + condAge + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	ageQuery := fmt.Sprintf(`
                SELECT
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) < 18 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 18 AND 25 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 26 AND 35 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) > 35 THEN 1 ELSE 0 END)
               FROM clients
               WHERE company_id=? AND branch_id=? AND id IN (SELECT DISTINCT client_id FROM bookings WHERE %s`, condAge)
	ageArgs = append([]interface{}{companyID, branchID, companyID, branchID}, ageArgs...)
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

	// Channel statistics
	// Use booking start time for consistent client counting
	condCh, chArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condCh = "b.company_id=? AND b.branch_id=? AND " + condCh + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	chQuery := fmt.Sprintf(`
               SELECT IFNULL(ch.name, ''), COUNT(*)
               FROM clients c
               LEFT JOIN channels ch ON c.channel_id = ch.id
               WHERE c.company_id=? AND c.branch_id=? AND c.id IN (SELECT DISTINCT client_id FROM bookings b WHERE %s`, condCh)
	chArgs = append([]interface{}{companyID, branchID, companyID, branchID}, chArgs...)
	if userID > 0 {
		chQuery += " AND b.user_id = ?"
		chArgs = append(chArgs, userID)
	}
	chQuery += `)`
	chQuery += " GROUP BY IFNULL(ch.name, '')"
	chRows, _ := r.db.QueryContext(ctx, chQuery, chArgs...)
	var chStats []models.ChannelStat
	for chRows.Next() {
		var name sql.NullString
		var count int
		chRows.Scan(&name, &count)
		if !name.Valid {
			name.String = ""
		}
		channelName := name.String
		if channelName == "" {
			channelName = "Гость"
		}
		chStats = append(chStats, models.ChannelStat{Channel: channelName, Clients: count})
	}

	guestCond, guestArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	guestCond = "b.company_id=? AND b.branch_id=? AND " + guestCond + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0 AND b.client_id IS NULL"
	guestQuery := fmt.Sprintf("SELECT COUNT(*) FROM bookings b WHERE %s", guestCond)
	guestArgs = append([]interface{}{companyID, branchID}, guestArgs...)
	if userID > 0 {
		guestQuery += " AND b.user_id = ?"
		guestArgs = append(guestArgs, userID)
	}
	var guestCount int
	_ = r.db.QueryRowContext(ctx, guestQuery, guestArgs...).Scan(&guestCount)
	if guestCount > 0 {
		found := false
		for i := range chStats {
			if chStats[i].Channel == "Гость" {
				chStats[i].Clients += guestCount
				found = true
				break
			}
		}
		if !found {
			chStats = append(chStats, models.ChannelStat{Channel: "Гость", Clients: guestCount})
		}
	}
	result.ChannelStats = chStats

	// Category sales
	condCat, catArgs := buildTimeCondition("bookings.start_time", from, to, tFrom, tTo)
	condCat = "bookings.company_id=? AND bookings.branch_id=? AND " + condCat + " AND bookings.payment_status <> 'UNPAID' AND bookings.payment_type_id <> 0"
	catQuery := fmt.Sprintf(`
               SELECT categories.name, SUM(booking_items.price  * (1 - booking_items.discount / 100) * (1 - IFNULL(pt.hold_percent,0)/100))
               FROM booking_items
               LEFT JOIN bookings ON booking_items.booking_id = bookings.id
               LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
               LEFT JOIN price_items ON booking_items.item_id = price_items.id
               LEFT JOIN categories ON price_items.category_id = categories.id
               WHERE %s`, condCat)
	catArgs = append([]interface{}{companyID, branchID}, catArgs...)
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
	condItem, itemArgs := buildTimeCondition("bookings.start_time", from, to, tFrom, tTo)
	condItem = "bookings.company_id=? AND bookings.branch_id=? AND " + condItem + " AND bookings.payment_status <> 'UNPAID' AND bookings.payment_type_id <> 0"
	itemQuery := fmt.Sprintf(`
    SELECT
        price_items.name,
        SUM(booking_items.quantity),
        SUM(
            CASE
                WHEN categories.name = 'Часы' THEN (booking_items.price * (1 - booking_items.discount / 100)) * (1 - IFNULL(pt.hold_percent,0)/100)
                ELSE booking_items.price  * (1 - booking_items.discount / 100) * (1 - IFNULL(pt.hold_percent,0)/100)
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
                WHEN categories.name = 'Часы' THEN (booking_items.price * (1 - booking_items.discount / 100))*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
                ELSE (booking_items.price * (1 - booking_items.discount / 100))*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
            END
        )
    FROM booking_items
    LEFT JOIN bookings ON booking_items.booking_id = bookings.id
    LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
    LEFT JOIN price_items ON booking_items.item_id = price_items.id
    LEFT JOIN categories ON price_items.category_id = categories.id
    WHERE %s`, condItem)
	itemArgs = append([]interface{}{companyID, branchID}, itemArgs...)
	if userID > 0 {
		itemQuery += " AND bookings.user_id = ?"
		itemArgs = append(itemArgs, userID)
	}

	itemQuery += `
    GROUP BY price_items.name
    ORDER BY SUM(
        CASE
            WHEN categories.name = 'Часы' THEN (booking_items.price * (1 - booking_items.discount / 100))*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
            ELSE (booking_items.price * (1 - booking_items.discount / 100))*(1 - IFNULL(pt.hold_percent,0)/100) - price_items.buy_price
        END
    ) DESC`

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
	condPrev, prevArgs := buildTimeCondition("b.start_time", prevFrom, prevTo, tFrom, tTo)
	condPrev = "b.company_id=? AND b.branch_id=? AND " + condPrev + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	prevQuery := fmt.Sprintf(`
       SELECT COALESCE(SUM(b.total_amount  * (1 - IFNULL(pt.hold_percent,0)/100)),0),
              COUNT(DISTINCT client_id) + SUM(CASE WHEN client_id IS NULL THEN 1 ELSE 0 END),
              COALESCE(AVG(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0)
       FROM bookings b
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       WHERE %s`, condPrev)
	prevArgs = append([]interface{}{companyID, branchID}, prevArgs...)
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

	return &result, nil
}

// --- AdminsReport ---
func (r *ReportRepository) AdminsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) (*models.AdminsReport, error) {

	condAdmin, adminArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condAdmin = "b.company_id=? AND b.branch_id=? AND " + condAdmin + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	query := fmt.Sprintf(`
       SELECT u.id, u.name,
              COUNT(DISTINCT DATE(b.start_time)) AS shifts,
              SUM(
                                CASE
					WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%%кальян%%' THEN bi.quantity
				 	WHEN pi.is_set = 1 THEN bi.quantity * (
                        SELECT COALESCE(SUM(si.quantity),0)
                        FROM set_items si
						JOIN price_items pi2 ON si.item_id = pi2.id
						JOIN categories c2 ON pi2.category_id = c2.id
						WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
				    )
					ELSE 0
				END) AS hookah_qty,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS set_qty,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%%кальян%%' THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS hookah_rev,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS set_rev,
              u.salary_shift, u.salary_hookah, u.hookah_salary_type, u.salary_bar
        FROM users u
       LEFT JOIN bookings b ON b.user_id = u.id
            AND %s
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       LEFT JOIN booking_items bi ON b.id = bi.booking_id
       LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE u.role = 'admin' AND u.company_id=? AND u.branch_id=?`, condAdmin)
	args := append([]interface{}{companyID, branchID}, adminArgs...)
	args = append(args, companyID, branchID)
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
func (r *ReportRepository) SalesReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) (*models.SalesReport, error) {
	condUser, userArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condUser = "b.company_id=? AND b.branch_id=? AND " + condUser + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	userQuery := fmt.Sprintf(`
       SELECT u.id, u.name,
              COUNT(DISTINCT DATE(b.start_time)) AS days,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%%кальян%%' THEN bi.quantity
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS hookahs,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.quantity
                      ELSE 0
                  END) AS sets,
              SUM(
                  CASE
                      WHEN pi.is_set = 0 AND LOWER(categories.name) LIKE '%%кальян%%' THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      WHEN pi.is_set = 1 AND EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
                      ) THEN bi.price * bi.quantity * (1 - IFNULL(pt.hold_percent,0)/100)
                      ELSE 0
                  END) AS hookah_rev,
              SUM(
                  CASE
                      WHEN pi.is_set = 1 AND NOT EXISTS (
                          SELECT 1 FROM set_items si
                          JOIN price_items pi2 ON si.item_id = pi2.id
                          JOIN categories c2 ON pi2.category_id = c2.id
                          WHERE si.price_set_id = pi.id AND LOWER(c2.name) LIKE '%%кальян%%'
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
       WHERE %s AND u.company_id=? AND u.branch_id=?`, condUser)
	userArgs = append([]interface{}{companyID, branchID}, userArgs...)
	userArgs = append(userArgs, companyID, branchID)
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

	condExp, expArgs := buildTimeCondition("e.date", from, to, tFrom, tTo)
	condExp = "e.company_id=? AND e.branch_id=? AND " + condExp
	expQuery := fmt.Sprintf(`
       SELECT IFNULL(ec.name, IFNULL(rc.name, '')) as category, SUM(e.total)
       FROM expenses e
        LEFT JOIN expense_categories ec ON e.category_id = ec.id
        LEFT JOIN repair_categories rc ON e.repair_category_id = rc.id
        WHERE %s
        GROUP BY category`, condExp)
	expArgs = append([]interface{}{companyID, branchID}, expArgs...)
	expRows, err := r.db.QueryContext(ctx, expQuery, expArgs...)
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

	condCat2, catArgs2 := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condCat2 = "b.company_id=? AND b.branch_id=? AND " + condCat2 + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	catQuery2 := fmt.Sprintf(`
        SELECT categories.name, SUM((bi.price * (1 - bi.discount / 100)) * (1 - IFNULL(pt.hold_percent,0)/100))
        FROM booking_items bi
        LEFT JOIN bookings b ON bi.booking_id = b.id
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories ON pi.category_id = categories.id
        WHERE %s`, condCat2)
	catArgs2 = append([]interface{}{companyID, branchID}, catArgs2...)
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

	// Income by payment type
	payCond, payArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	payCond = "b.company_id=? AND b.branch_id=? AND " + payCond + " AND b.payment_status <> 'UNPAID'"
	payQuery := fmt.Sprintf(`
       SELECT IFNULL(pt.name,''), SUM(bp.amount * (1 - IFNULL(pt.hold_percent,0)/100))
       FROM bookings b
       LEFT JOIN booking_payments bp ON b.id = bp.booking_id AND b.company_id = bp.company_id AND b.branch_id = bp.branch_id
       LEFT JOIN payment_types pt ON bp.payment_type_id = pt.id
       WHERE %s AND bp.payment_type_id <> 0`, payCond)
	payArgs = append([]interface{}{companyID, branchID}, payArgs...)
	if userID > 0 {
		payQuery += " AND b.user_id = ?"
		payArgs = append(payArgs, userID)
	}
	payQuery += ` GROUP BY pt.name`
	payRows, err := r.db.QueryContext(ctx, payQuery, payArgs...)
	if err != nil {
		return nil, err
	}
	defer payRows.Close()
	var payIncome []models.CategoryIncome
	for payRows.Next() {
		var inc models.CategoryIncome
		if err := payRows.Scan(&inc.Category, &inc.Total); err != nil {
			return nil, err
		}
		payIncome = append(payIncome, inc)
	}

	const taxPercent = 0
	netProfit := totalInc*(1-taxPercent) - totalExp

	return &models.SalesReport{
		Users:            users,
		Expenses:         expenses,
		IncomeByCategory: incomes,
		IncomeByPayment:  payIncome,
		TotalIncome:      totalInc,
		TotalExpenses:    totalExp,
		NetProfit:        netProfit,
	}, nil
}

// --- AnalyticsReport ---
func (r *ReportRepository) AnalyticsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) (*models.AnalyticsReport, error) {
	// Daily revenue
	condDaily, dailyArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
	condDaily = "b.company_id=? AND b.branch_id=? AND " + condDaily + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
	dailyQuery := fmt.Sprintf(`
       SELECT DATE(b.start_time), SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)) FROM bookings b
       LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
       WHERE %s`, condDaily)
	dailyArgs = append([]interface{}{companyID, branchID}, dailyArgs...)
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
	condHourly, hourlyArgs := buildTimeCondition("start_time", from, to, tFrom, tTo)
	condHourly = "company_id=? AND branch_id=? AND " + condHourly + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	hourlyQuery := fmt.Sprintf(`
       SELECT HOUR(start_time), COUNT(*) FROM bookings
       WHERE %s`, condHourly)
	hourlyArgs = append([]interface{}{companyID, branchID}, hourlyArgs...)
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
	condAnalCat, catArgs := buildTimeCondition("booking_items.created_at", from, to, tFrom, tTo)
	condAnalCat = "bookings.company_id=? AND bookings.branch_id=? AND " + condAnalCat + " AND bookings.payment_status <> 'UNPAID' AND bookings.payment_type_id <> 0"
	catQuery := fmt.Sprintf(`
       SELECT categories.name, SUM(booking_items.quantity), SUM(booking_items.price * (1 - IFNULL(pt.hold_percent,0)/100))
       FROM booking_items
       LEFT JOIN bookings ON booking_items.booking_id = bookings.id
       LEFT JOIN payment_types pt ON bookings.payment_type_id = pt.id
       LEFT JOIN price_items ON booking_items.item_id = price_items.id
       LEFT JOIN categories ON price_items.category_id = categories.id
       WHERE %s`, condAnalCat)
	catArgs = append([]interface{}{companyID, branchID}, catArgs...)
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
func (r *ReportRepository) DiscountsReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) (*models.DiscountsReport, error) {
	var total, count, avg int
	condSum, sumArgs := buildTimeCondition("start_time", from, to, tFrom, tTo)
	condSum = "company_id=? AND branch_id=? AND " + condSum + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	sumQuery := fmt.Sprintf(`
       SELECT COALESCE(SUM(discount),0), COUNT(*), COALESCE(AVG(discount),0)
       FROM bookings
       WHERE discount > 0 AND %s`, condSum)
	sumArgs = append([]interface{}{companyID, branchID}, sumArgs...)
	if userID > 0 {
		sumQuery += " AND user_id = ?"
		sumArgs = append(sumArgs, userID)
	}
	_ = r.db.QueryRowContext(ctx, sumQuery, sumArgs...).Scan(&total, &count, &avg)

	condReason, reasonArgs := buildTimeCondition("start_time", from, to, tFrom, tTo)
	condReason = "company_id=? AND branch_id=? AND " + condReason + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	reasonQuery := fmt.Sprintf(`
       SELECT discount_reason, COUNT(*), SUM(discount), COALESCE(AVG(discount),0)
       FROM bookings
       WHERE discount > 0 AND %s`, condReason)
	reasonArgs = append([]interface{}{companyID, branchID}, reasonArgs...)
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

	condDist, distArgs := buildTimeCondition("start_time", from, to, tFrom, tTo)
	condDist = "company_id=? AND branch_id=? AND " + condDist + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	distQuery := fmt.Sprintf(`
       SELECT discount, COUNT(*) FROM bookings
       WHERE discount > 0 AND %s`, condDist)
	distArgs = append([]interface{}{companyID, branchID}, distArgs...)
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
	condOrders, orderArgs := buildTimeCondition("start_time", from, to, tFrom, tTo)
	condOrders = "company_id=? AND branch_id=? AND " + condOrders + " AND payment_status <> 'UNPAID' AND payment_type_id <> 0"
	orderQuery := fmt.Sprintf(`
       SELECT id, client_id, table_id, user_id, start_time, end_time, note,
               discount, discount_reason, total_amount, bonus_used,
               payment_status, payment_type_id, created_at, updated_at
       FROM bookings
       WHERE discount > 0 AND %s`, condOrders)
	orderArgs = append([]interface{}{companyID, branchID}, orderArgs...)
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
		var clientID, tableID sql.NullInt64
		if err := orderRows.Scan(&b.ID, &clientID, &tableID, &b.UserID, &b.StartTime,
			&b.EndTime, &b.Note, &b.Discount, &b.DiscountReason, &b.TotalAmount,
			&b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		if clientID.Valid {
			b.ClientID = int(clientID.Int64)
		}
		if tableID.Valid {
			b.TableID = int(tableID.Int64)
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

func (r *ReportRepository) TablesReport(ctx context.Context, from, to time.Time, tFrom, tTo string, userID, companyID, branchID int) ([]models.TableReport, error) {
	var workFrom, workTo string
	_ = r.db.QueryRowContext(ctx, `SELECT work_time_from, work_time_to FROM settings WHERE company_id=? AND branch_id=? LIMIT 1`, companyID, branchID).Scan(&workFrom, &workTo)
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
	capacity := hours * float64(days)

	tableRows, err := r.db.QueryContext(ctx, `SELECT id, name FROM tables WHERE company_id=? AND branch_id=?`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer tableRows.Close()

	var result []models.TableReport
	for tableRows.Next() {
		var tID int
		var tName string
		if err := tableRows.Scan(&tID, &tName); err != nil {
			return nil, err
		}

		cond, condArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
		cond = "b.company_id=? AND b.branch_id=? AND b.table_id=? AND " + cond + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
		args := append([]interface{}{companyID, branchID, tID}, condArgs...)
		query := fmt.Sprintf(`SELECT COALESCE(SUM(b.total_amount),0), COALESCE(SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0), COALESCE(COUNT(DISTINCT b.client_id) + SUM(CASE WHEN b.client_id IS NULL THEN 1 ELSE 0 END),0), COALESCE(ROUND(AVG(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100))),0), COUNT(*) FROM bookings b LEFT JOIN payment_types pt ON b.payment_type_id = pt.id WHERE %s`, cond)
		if userID > 0 {
			query += " AND b.user_id = ?"
			args = append(args, userID)
		}
		var total, totalRev, clients, avgCheck float64
		var bookingsCount int
		if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total, &totalRev, &clients, &avgCheck, &bookingsCount); err != nil {
			return nil, err
		}
		loadPercent := 0.0
		if capacity > 0 {
			loadPercent = float64(bookingsCount) * 100 / capacity
		}

		condCost, costArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
		condCost = "b.company_id=? AND b.branch_id=? AND b.table_id=? AND " + condCost + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
		costArgs = append([]interface{}{companyID, branchID, tID}, costArgs...)
		costQuery := fmt.Sprintf(`SELECT COALESCE(SUM(
            bi.price * (1 - bi.discount / 100) * (1 - IFNULL(pt.hold_percent,0)/100) -
            bi.quantity * pi.buy_price
        ),0)
        FROM booking_items bi
        LEFT JOIN bookings b ON bi.booking_id = b.id
        LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
        LEFT JOIN price_items pi ON bi.item_id = pi.id
        LEFT JOIN categories c ON pi.category_id = c.id
        WHERE %s AND (c.name = 'Бар' OR c.name = 'Кальян')`, condCost)
		if userID > 0 {
			costQuery += " AND b.user_id = ?"
			costArgs = append(costArgs, userID)
		}
		var totalCost float64
		_ = r.db.QueryRowContext(ctx, costQuery, costArgs...).Scan(&totalCost)

		condPay, payArgs := buildTimeCondition("b.start_time", from, to, tFrom, tTo)
		condPay = "b.company_id=? AND b.branch_id=? AND b.table_id=? AND " + condPay + " AND b.payment_status <> 'UNPAID' AND b.payment_type_id <> 0"
		payArgs = append([]interface{}{companyID, branchID, tID}, payArgs...)
		payQuery := fmt.Sprintf(`SELECT IFNULL(pt.name,''), COALESCE(SUM(b.total_amount * (1 - IFNULL(pt.hold_percent,0)/100)),0)
            FROM bookings b
            LEFT JOIN payment_types pt ON b.payment_type_id = pt.id
            WHERE %s`, condPay)
		if userID > 0 {
			payQuery += " AND b.user_id = ?"
			payArgs = append(payArgs, userID)
		}
		payQuery += " GROUP BY IFNULL(pt.name,'')"
		payRows, _ := r.db.QueryContext(ctx, payQuery, payArgs...)
		var payStats []models.CategoryIncome
		for payRows.Next() {
			var name sql.NullString
			var amt float64
			payRows.Scan(&name, &amt)
			n := name.String
			if !name.Valid {
				n = ""
			}
			payStats = append(payStats, models.CategoryIncome{Category: n, Total: amt})
		}
		payRows.Close()

		result = append(result, models.TableReport{
			TableID:           tID,
			TableName:         tName,
			Total:             total,
			TotalRevenue:      totalRev,
			AvgCheck:          avgCheck,
			LoadPercent:       loadPercent,
			Visits:            clients,
			PaymentTypeIncome: payStats,
			TotalCost:         totalCost,
		})
	}

	return result, nil
}
