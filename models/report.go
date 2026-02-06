package models

type BestSellerItem struct {
	ProductName  string `json:"product_name"`
	QuantitySold int    `json:"quantity_sold"`
}

type DailyReport struct {
	TotalRevenue int            `json:"total_revenue"`
	TotalSales   int            `json:"total_sales"`
	BestSellers  BestSellerItem `json:"best_sellers"`
}

type ProductStat struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

type TrendStat struct {
	Date         string `json:"date"`
	TotalRevenue int    `json:"daily_revenue"`
}

type ReportSummary struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalTransactions int     `json:"total_transactions"`
	AOV               float64 `json:"avg_order_value"`
	TotalItemsSold    int     `json:"total_items_sold"`
}

type FullReport struct {
	Summary           ReportSummary `json:"summary"`
	Top5Products      []ProductStat `json:"top_5_products"`
	DailyRevenueTrend []TrendStat   `json:"daily_revenue_trend"`
}
