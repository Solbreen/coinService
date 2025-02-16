package database

import (
	"coinService/internal/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func Connect() (*DB, error) {
	connStr := "user=user dbname=avito_merch password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &DB{conn: db}, nil
}

func (db *DB) Authenticate(username, password string) (*models.User, error) {
	var user models.User
	err := db.conn.QueryRow("SELECT id, username FROM users WHERE username = $1 AND password_hash = $2", username, password).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to authenticate user: %v", err)
	}
	return &user, nil
}

func (db *DB) GetUserInfo(userID int) (*models.InfoResponse, error) {
	var info models.InfoResponse

	err := db.conn.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&info.Coins)
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query("SELECT item_name, quantity FROM inventory WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		info.Inventory = append(info.Inventory, item)
	}

	rows, err = db.conn.Query(`
        SELECT u1.username, u2.username, t.amount
        FROM transactions t
        JOIN users u1 ON t.from_user_id = u1.id
        JOIN users u2 ON t.to_user_id = u2.id
        WHERE t.from_user_id = $1 OR t.to_user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.FromUser, &transaction.ToUser, &transaction.Amount); err != nil {
			return nil, err
		}
		if transaction.FromUser == "user1" {
			info.CoinHistory.Sent = append(info.CoinHistory.Sent, transaction)
		} else {
			info.CoinHistory.Received = append(info.CoinHistory.Received, transaction)
		}
	}

	return &info, nil
}

func (db *DB) SendCoin(fromUserID int, toUser string, amount int) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	var toUserID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = $1", toUser).Scan(&toUserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("user not found")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, fromUserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct coins")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, toUserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add coins")
	}

	_, err = tx.Exec("INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", fromUserID, toUserID, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record transaction")
	}

	return tx.Commit()
}

func (db *DB) BuyItem(userID int, item string) error {
	var price int
	switch item {
	case "t-shirt":
		price = 80
	case "cup":
		price = 20
	case "book":
		price = 50
	case "pen":
		price = 10
	case "powerbank":
		price = 200
	case "hoody":
		price = 300
	case "umbrella":
		price = 200
	case "socks":
		price = 10
	case "wallet":
		price = 50
	case "pink-hoody":
		price = 500
	default:
		return fmt.Errorf("unknown item: %s", item)
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("not enough coins")
	}

	_, err = tx.Exec("INSERT INTO inventory (user_id, item_name, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_name) DO UPDATE SET quantity = inventory.quantity + 1", userID, item)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update inventory")
	}

	return tx.Commit()
}
