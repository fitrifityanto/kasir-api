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
