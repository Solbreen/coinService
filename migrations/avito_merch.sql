CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    coins INT DEFAULT 1000
);

CREATE TABLE inventory (
    user_id INT REFERENCES users(id),
    item_name TEXT NOT NULL,
    quantity INT DEFAULT 0,
    PRIMARY KEY (user_id, item_name)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT REFERENCES users(id),
    to_user_id INT REFERENCES users(id),
    amount INT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);