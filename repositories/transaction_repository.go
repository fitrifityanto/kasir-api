package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisialisasi subtotal --> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetails, 0)
	// loop setiap item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product untuk dapetin price
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// cek stock
		if stock < item.Quantity {
			return nil, fmt.Errorf("product %s only has %d stock", productName, stock)
		}
		// hitung current total = quantity * price
		// ditambahin kedalalm subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangin stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// itemnya dimasukin ke transactionDetails
		details = append(details, models.TransactionDetails{
			ProductID:   productID,
			ProductName: productName,
			Price:       price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	// insert transactionDetails
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO transaction_details (transaction_id, product_id, product_name, price, quantity, subtotal) VALUES ")

	values := []any{}

	for i, detail := range details {
		n := i * 6
		fmt.Fprintf(&queryBuilder, "($%d, $%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5, n+6)

		if i < len(details)-1 {
			queryBuilder.WriteString(",")
		}

		details[i].TransactionID = transactionID
		values = append(values, transactionID, detail.ProductID, detail.ProductName, detail.Price, detail.Quantity, detail.Subtotal)

	}

	queryBuilder.WriteString(" RETURNING id")
	query := queryBuilder.String()
	rows, err := tx.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
		i++
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		CreatedAt:   createdAt,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}
