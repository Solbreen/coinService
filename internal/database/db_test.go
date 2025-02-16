package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)", "testuser", "testpass", 1000)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO inventory (user_id, item_name, quantity) VALUES ($1, $2, $3)", 1, "t-shirt", 2)
	if err != nil {
		t.Fatalf("Failed to add item to inventory: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", 1, 2, 100)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	info, err := db.GetUserInfo(1)
	if err != nil {
		t.Fatalf("Failed to get user info: %v", err)
	}

	assert.Equal(t, 1000, info.Coins, "Expected 1000 coins")
	assert.Equal(t, 1, len(info.Inventory), "Expected 1 item in inventory")
	assert.Equal(t, "t-shirt", info.Inventory[0].Type, "Expected t-shirt in inventory")
	assert.Equal(t, 2, info.Inventory[0].Quantity, "Expected 2 t-shirts in inventory")
	assert.Equal(t, 1, len(info.CoinHistory.Sent), "Expected 1 sent transaction")
	assert.Equal(t, 0, len(info.CoinHistory.Received), "Expected 0 received transactions")

	_, err = db.conn.Exec("DELETE FROM users WHERE username = 'testuser'")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

func TestSendCoin(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)", "user1", "pass1", 1000)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	_, err = db.conn.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)", "user2", "pass2", 1000)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = db.SendCoin(1, "user2", 100)
	if err != nil {
		t.Fatalf("Failed to send coins: %v", err)
	}

	var user1Coins, user2Coins int
	err = db.conn.QueryRow("SELECT coins FROM users WHERE id = 1").Scan(&user1Coins)
	if err != nil {
		t.Fatalf("Failed to get user1 coins: %v", err)
	}
	err = db.conn.QueryRow("SELECT coins FROM users WHERE id = 2").Scan(&user2Coins)
	if err != nil {
		t.Fatalf("Failed to get user2 coins: %v", err)
	}

	assert.Equal(t, 900, user1Coins, "Expected user1 to have 900 coins")
	assert.Equal(t, 1100, user2Coins, "Expected user2 to have 1100 coins")

	var amount int
	err = db.conn.QueryRow("SELECT amount FROM transactions WHERE from_user_id = 1 AND to_user_id = 2").Scan(&amount)
	if err != nil {
		t.Fatalf("Failed to get transaction: %v", err)
	}
	assert.Equal(t, 100, amount, "Expected transaction amount to be 100")

	_, err = db.conn.Exec("DELETE FROM users WHERE username IN ('user1', 'user2')")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

func TestBuyItem(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)", "testuser", "testpass", 1000)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = db.BuyItem(1, "t-shirt")
	if err != nil {
		t.Fatalf("Failed to buy item: %v", err)
	}

	var coins int
	err = db.conn.QueryRow("SELECT coins FROM users WHERE id = 1").Scan(&coins)
	if err != nil {
		t.Fatalf("Failed to get user coins: %v", err)
	}
	assert.Equal(t, 920, coins, "Expected user to have 920 coins after buying t-shirt")

	var quantity int
	err = db.conn.QueryRow("SELECT quantity FROM inventory WHERE user_id = 1 AND item_name = 't-shirt'").Scan(&quantity)
	if err != nil {
		t.Fatalf("Failed to get inventory: %v", err)
	}
	assert.Equal(t, 1, quantity, "Expected 1 t-shirt in inventory")

	_, err = db.conn.Exec("DELETE FROM users WHERE username = 'testuser'")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

func TestAuthenticate(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.conn.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)", "testuser", "testpass", 1000)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user, err := db.Authenticate("testuser", "testpass")
	if err != nil {
		t.Fatalf("Failed to authenticate user: %v", err)
	}

	assert.Equal(t, 1, user.ID, "Expected user ID to be 1")
	assert.Equal(t, "testuser", user.Username, "Expected username to be 'testuser'")

	_, err = db.Authenticate("testuser", "wrongpass")
	assert.NotNil(t, err, "Expected error for wrong password")

	_, err = db.conn.Exec("DELETE FROM users WHERE username = 'testuser'")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}
