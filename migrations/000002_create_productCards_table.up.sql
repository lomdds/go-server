CREATE TABLE product_cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    bike_model VARCHAR(255) NOT NULL,
    engine_capacity INTEGER NOT NULL,
    power INTEGER NOT NULL,
    color VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_product_cards_user_id ON product_cards(user_id);