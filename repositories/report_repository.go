package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	var r models.DailyReport
	query := `
WITH stats AS (
    SELECT 
        COALESCE(SUM(total_amount), 0) AS total_revenue,
        COUNT(id) AS total_sales
    FROM transactions
    WHERE created_at::date = CURRENT_DATE
),
top_product AS (
    SELECT 
        product_name,
        SUM(quantity) AS quantity_sold
    FROM transaction_details td
    JOIN transactions t ON td.transaction_id = t.id
    WHERE t.created_at::date = CURRENT_DATE
    GROUP BY product_name
    ORDER BY quantity_sold DESC, SUM(subtotal) DESC
    LIMIT 1
)
SELECT 
    s.total_revenue, 
    s.total_sales,
    COALESCE(tp.product_name, 'No Sales') AS product_name,
    COALESCE(tp.quantity_sold, 0) AS quantity_sold
FROM stats s
LEFT JOIN top_product tp ON true;
	`

	err := repo.db.QueryRow(query).Scan(
		&r.TotalRevenue,
		&r.TotalSales,
		&r.BestSellers.ProductName,
		&r.BestSellers.QuantitySold,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}
