CREATE TABLE login_history (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT NOW()
);