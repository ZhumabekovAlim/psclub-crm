package models

type SummaryReport struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalClients   int            `json:"total_clients"`
	AvgCheck       int            `json:"avg_check"`
	LoadPercent    int            `json:"load_percent"`
	AgeUnder18     float64        `json:"age_under_18_percent"`
	Age18To25      float64        `json:"age_18_25_percent"`
	Age26To35      float64        `json:"age_26_35_percent"`
	Age36Plus      float64        `json:"age_36_plus_percent"`
	CategorySales  []CategorySale `json:"category_sales"`
	TopItems       []ProfitItem   `json:"top_items"`
	RevenueChange  float64        `json:"revenue_change_percent"`
	ClientsChange  float64        `json:"clients_change_percent"`
	AvgCheckChange float64        `json:"avg_check_change_percent"`
	LoadChange     float64        `json:"load_change_percent"`
}

type AdminsReport struct {
	Admins []AdminReportRow `json:"admins"`
}

type AdminReportRow struct {
	Name         string `json:"name"`
	Shifts       int    `json:"shifts"`
	HookahsSold  int    `json:"hookahs_sold"`
	SetsSold     int    `json:"sets_sold"`
	Salary       int    `json:"salary"`
	SalaryDetail string `json:"salary_detail"`
}

type SaleItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Revenue  int    `json:"revenue"`
	AvgCheck int    `json:"avg_check"`
	Category string `json:"category,omitempty"`
}

type UserSales struct {
	Name        string `json:"name"`
	DaysWorked  int    `json:"days_worked"`
	HookahsSold int    `json:"hookahs_sold"`
	SetsSold    int    `json:"sets_sold"`
	Salary      int    `json:"salary"`
}

type ExpenseTotal struct {
	Title string  `json:"title"`
	Total float64 `json:"total"`
}

type CategoryIncome struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type SalesReport struct {
	Users            []UserSales      `json:"users"`
	Expenses         []ExpenseTotal   `json:"expenses,omitempty"`
	IncomeByCategory []CategoryIncome `json:"income_by_category,omitempty"`
	TotalIncome      float64          `json:"total_income,omitempty"`
	TotalExpenses    float64          `json:"total_expenses,omitempty"`
	NetProfit        float64          `json:"net_profit,omitempty"`
}

type AnalyticsReport struct {
	DailyRevenue  []DataPoint    `json:"daily_revenue"`
	HourlyLoad    []DataPoint    `json:"hourly_load"`
	CategoryStats []CategoryStat `json:"category_stats"`
}

type DataPoint struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	Revenue  int    `json:"revenue"`
	AvgCheck int    `json:"avg_check"`
}

type DiscountsReport struct {
	TotalDiscount     int         `json:"total_discount"`
	DiscountCount     int         `json:"discount_count"`
	AvgDiscount       int         `json:"avg_discount"`
	TopReasons        []ReasonRow `json:"top_reasons"`
	DistributionBySum []DataPoint `json:"distribution_by_sum"`
	Orders            []Booking   `json:"orders"`
}

type ReasonRow struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
	Sum    int    `json:"sum"`
	Avg    int    `json:"avg"`
}

type CategorySale struct {
	Category string `json:"category"`
	Revenue  int    `json:"revenue"`
}

type ProfitItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
	Expense  float64 `json:"expense"`
	Profit   float64 `json:"profit"`
}
