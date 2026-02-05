package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
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
	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) VALUES ($1, $2, $3, $4, $5) RETURNING id", transactionID, details[i].ProductID, details[i].Price, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
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
