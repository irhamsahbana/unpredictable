CREATE TABLE IF NOT EXISTS product_transactions (
    id VARCHAR(255) PRIMARY KEY,
    member_id VARCHAR(255) NOT NULL,
    product_id BIGINT,
    product_grammage_id BIGINT,
    source VARCHAR(50),
    qty INTEGER,
    price_per_unit DECIMAL(19, 4),
    created_at TIMESTAMPTZ,
    is_training_data BOOLEAN DEFAULT FALSE
);