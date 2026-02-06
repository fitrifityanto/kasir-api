package repositories

import (
	"database/sql"
	"encoding/json"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {

	query := `
WITH stats AS (
    SELECT 
        COALESCE(SUM(total_amount), 0) AS total_revenue,
        COUNT(id) AS total_sales
    FROM transactions
    WHERE created_at >= (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::date
    AND created_at < ((CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta') + INTERVAL '1 day')::date
),
top_product AS (
    SELECT 
        product_name,
        SUM(quantity) AS quantity_sold
    FROM transaction_details td
    JOIN transactions t ON td.transaction_id = t.id
    WHERE t.created_at >= (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::date
    AND t.created_at < ((CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta') + INTERVAL '1 day')::date
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
	var r models.DailyReport
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

func (repo *ReportRepository) GetReport(startDate, endDate string) (*models.FullReport, error) {
	now := time.Now()
	if startDate == "" {
		startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	query := `
WITH date_range AS (
    SELECT 
        ($1::date AT TIME ZONE 'Asia/Jakarta') AS start_ts,
        ($2::date + INTERVAL '1 day') AT TIME ZONE 'Asia/Jakarta' AS end_ts
),
filtered_transactions AS (
    SELECT 
        t.id, 
        t.total_amount, 
        (t.created_at AT TIME ZONE 'Asia/Jakarta')::date AS day,
        t.created_at
    FROM transactions t, date_range dr
    WHERE t.created_at >= dr.start_ts AND t.created_at < dr.end_ts
),
summary AS (
    SELECT 
        COALESCE(SUM(total_amount), 0) AS total_revenue,
        COUNT(id) AS total_transactions,
        COALESCE(SUM(total_amount) / NULLIF(COUNT(id), 0), 0) AS avg_order_value
    FROM filtered_transactions
),
top_products AS (
    SELECT jsonb_agg(jsonb_build_object('product_name', product_name, 'quantity', total_qty)) as list
    FROM (
        SELECT td.product_name, SUM(td.quantity) as total_qty
        FROM transaction_details td
        JOIN filtered_transactions ft ON td.transaction_id = ft.id
        GROUP BY 1 ORDER BY 2 DESC LIMIT 5
    ) AS sub
),
daily_trend AS (
    SELECT jsonb_agg(jsonb_build_object('date', day, 'daily_revenue', daily_revenue)) as trend
    FROM (
        SELECT day, SUM(total_amount) AS daily_revenue
        FROM filtered_transactions
        GROUP BY 1 ORDER BY 1 ASC
    ) AS sub
),
total_qty AS (
    SELECT COALESCE(SUM(td.quantity), 0) as total_items_sold
    FROM transaction_details td
    JOIN filtered_transactions ft ON td.transaction_id = ft.id
)
SELECT 
    s.total_revenue, 
    s.total_transactions, 
    s.avg_order_value, 
    tq.total_items_sold,
    COALESCE(tp.list, '[]'::jsonb) as top_5_products,
    COALESCE(dt.trend, '[]'::jsonb) as daily_revenue_trend
FROM summary s, top_products tp, daily_trend dt, total_qty tq;
	`
	var r models.FullReport
	var topProductsJSON, dailyTrendJSON []byte
	err := repo.db.QueryRow(query, startDate, endDate).Scan(
		&r.Summary.TotalRevenue,
		&r.Summary.TotalTransactions,
		&r.Summary.AOV,
		&r.Summary.TotalItemsSold,
		&topProductsJSON,
		&dailyTrendJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(topProductsJSON, &r.Top5Products); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dailyTrendJSON, &r.DailyRevenueTrend); err != nil {
		return nil, err
	}
	return &r, nil
}
