CREATE DATABASE subscriptions;
USE subscriptions;

CREATE TABLE users (
    id VARCHAR(21),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    dob DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT uc_email UNIQUE (email)
);

CREATE INDEX idx_user_email ON users(email);

CREATE TABLE products (
    id VARCHAR(21),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    duration_days integer NOT NULL,
    price DECIMAL(6,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT uc_products_name_duration_days UNIQUE (name, duration_days)
);

CREATE TABLE vouchers (
    code VARCHAR(21),
    type ENUM('FIXED_AMOUNT', 'PERCENTAGE') NOT NULL,
    value DECIMAL(6,2) NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    max_uses integer,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (code),
    CONSTRAINT chk_vouchers_value CHECK (value>0.00 AND ((type = 'PERCENTAGE' AND value<=100.00) OR (type = 'FIXED_AMOUNT'))),
    CONSTRAINT chk_vouchers_valid_validity_range CHECK (valid_until>valid_from)
);

CREATE TABLE products_vouchers (
    product_id VARCHAR(21),
    voucher_code VARCHAR(21),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, voucher_code),
    CONSTRAINT fk_products_vouchers_product_id FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_products_vouchers_voucher_code FOREIGN KEY (voucher_code) REFERENCES vouchers(code)
);

CREATE TABLE subscription_plans (
    id VARCHAR(21),
    user_id VARCHAR(21) NOT NULL,
    product_id VARCHAR(21) NOT NULL,
    status ENUM('ACTIVE', 'EXPIRED', 'PAUSED', 'CANCELLED') NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    net_price DECIMAL(6,2) NOT NULL,
    gross_price DECIMAL(6,2) NOT NULL,
    tax DECIMAL(6,2) NOT NULL,
    discount DECIMAL(6,2) DEFAULT 0.0,
    voucher_code VARCHAR(21),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_subscription_plans_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_subscription_plans_product_id FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_subscription_plans_voucher_code FOREIGN KEY (voucher_code) REFERENCES vouchers(code),
    CONSTRAINT chk_subscription_plans_duration_range CHECK (end_date>start_date)
);

CREATE TABLE subscription_plan_pauses (
    id VARCHAR(21),
    subscription_plan_id VARCHAR(21) NOT NULL,
    pause_date TIMESTAMP NOT NULL,
    end_date_at_pause TIMESTAMP NOT NULL,
    resumed TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT chk_subscription_plan_pauses_time_ranges CHECK (end_date_at_pause>pause_date)
);

CREATE INDEX idx_subscription_plan_pauses_subscription_plan_id ON subscription_plan_pauses(subscription_plan_id);