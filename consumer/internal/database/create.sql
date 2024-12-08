-- Таблица для категорий
CREATE TABLE CategoryItem (
    id INT PRIMARY KEY,
    category VARCHAR(255) NOT NULL
);

-- Таблица для товаров
CREATE TABLE Item (
    id INT PRIMARY KEY,
    track_number VARCHAR(255),
    category_id INT,
    price INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    size VARCHAR(50),
    total_price INT NOT NULL,
    brand VARCHAR(255),
    status INT,
    FOREIGN KEY (category_id) REFERENCES CategoryItem(id)
);

-- Таблица для курьеров
CREATE TABLE DeliveryMan (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(255)
);

-- Таблица для платежей
CREATE TABLE Payment (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(100),
    amount INT NOT NULL,
    payment_dt INT,
    bank VARCHAR(100),
    delivery_cost INT,
    custom_fee INT
);

-- Таблица для заказов
CREATE TABLE Order (
    id SERIAL PRIMARY KEY,
    delivery_id INT,
    payment_id INT,
    locale VARCHAR(10),
    delivery_service VARCHAR(100),
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (delivery_id) REFERENCES DeliveryMan(id),
    FOREIGN KEY (payment_id) REFERENCES Payment(id)
);

-- Таблица для статусов заказов
CREATE TABLE OrderStatus (
    id SERIAL PRIMARY KEY,
    order_id INT,
    status VARCHAR(50),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Order(id)
);