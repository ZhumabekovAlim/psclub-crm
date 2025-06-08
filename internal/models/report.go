package models

type SummaryReport struct {
	TotalRevenue   int     `json:"total_revenue"`
	TotalClients   int     `json:"total_clients"`
	AvgCheck       int     `json:"avg_check"`
	LoadPercent    int     `json:"load_percent"`
	RevenueChange  float64 `json:"revenue_change_percent"`
	ClientsChange  float64 `json:"clients_change_percent"`
	AvgCheckChange float64 `json:"avg_check_change_percent"`
	LoadChange     float64 `json:"load_change_percent"`
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
type SalesReport struct {
	TopSales []SaleItem `json:"top_sales"`
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
}
type ReasonRow struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
	Sum    int    `json:"sum"`
	Avg    int    `json:"avg"`
}
