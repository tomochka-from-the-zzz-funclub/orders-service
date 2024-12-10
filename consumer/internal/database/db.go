package database

import (
	"consumer/internal/config"
	myErrors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type Postgres struct {
	Connection *sql.DB
}

func NewPostgres(cfg config.Config) *Postgres {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		myLog.Log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil //, myErrors.ErrCreatePostgresConnection
	}
	time.Sleep(time.Minute)
	err = db.Ping()
	if err != nil {
		myLog.Log.Fatalf("Failed to ping PostgreSQL: %v", err)
		return nil //, myErrors.ErrPing
	} else {
		myLog.Log.Debugf("ping success")
	}
	time.Sleep(time.Minute)
	query :=
		` 
	CREATE TABLE IF NOT EXISTS  category_item (
    id SERIAL PRIMARY KEY,
    category VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY ,
    track_number VARCHAR(255),
    category_id INT,
    price INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    size VARCHAR(50),
    total_price INT NOT NULL,
    brand VARCHAR(255),
    status INT
);

CREATE TABLE IF NOT EXISTS delivery_man (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS  payment (
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

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    payment_id INT,
	items_id INT[],
    locale VARCHAR(10),
    delivery_service VARCHAR(100),
    date_created VARCHAR(40)
);

CREATE TABLE IF NOT EXISTS order_status (
    order_id INT,
    status VARCHAR(50),
    updated_at VARCHAR(40)
);

CREATE TABLE IF NOT EXISTS delivery (
    order_id INT,
    delivery_id VARCHAR(50),
    updated_at VARCHAR(40)
);
`
	_, err = db.Exec(query)
	if err != nil {
		myLog.Log.Errorf(err.Error())
	}
	return &Postgres{
		Connection: db,
	}
}

func (db *Postgres) AddOrderStruct(order models.Order) (int, error) {
	myLog.Log.Debugf("Go to bd in Set")

	// id_delivery, err := db.AddDeliveryMan(order.DeliveryMan)
	// if err != nil {
	// 	myLog.Log.Errorf("Error CreateDelivery: %v", err.Error())
	// 	return 0, err
	// }
	id_payment, err := db.AddPayment(order.Payment)
	if err != nil {
		myLog.Log.Errorf("Error CreatePayment: %v", err.Error())
		return 0, err
	}

	id_item, err := db.AddItemsWithCategory(order.Items)
	if err != nil {
		myLog.Log.Errorf("Error CreateItems: %v", err.Error())
		return 0, err
	}
	// id_shipping_method, err := db.AddShippingMethod(order.ShippingMethod)
	// if err !=nil {
	// 	myLog.Log.Errorf("Error CreateShippingMethods: %v", err.Error())
	// 	return "", err
	// }

	id, err := db.AddOrder(order, id_item, id_payment)
	if err != nil {
		if err == sql.ErrNoRows {
			myLog.Log.Debugf("The entry was not added. No data returned: %+v", err.Error())
			//return 0, err
		} else {
			myLog.Log.Errorf("Error CreateOrder: %v", err.Error())
			return 0, err
		}

	}
	order.Id = id
	err = db.AddOrderStatus(id)
	if err != nil {
		myLog.Log.Errorf("Error Create Status Order: %v", err.Error())
		return 0, err
	}
	myLog.Log.Infof("Entry successfully added with ID: ", id)

	return id, nil
}

func (db *Postgres) AddDeliveryMan(delivery models.DeliveryMan) (int, error) {
	query_add_delivery :=
		`WITH insert_return AS (
		 INSERT INTO  delivery_man (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
        )
        SELECT id FROM insert_return`
	var id string
	err := db.Connection.QueryRow(query_add_delivery, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	id_, err := strconv.Atoi(id)
	return id_, nil
}

func (db *Postgres) AddDeliveryMach(order_id int, delivery_man_id int) error {
	var count_d, count_o int
	err := db.Connection.QueryRow("SELECT COUNT(*) FROM delivery_man WHERE id = $1", delivery_man_id).Scan(&count_d)
	if err != nil {
		myLog.Log.Errorf(err.Error())
	}

	err = db.Connection.QueryRow("SELECT COUNT(*) FROM orders WHERE id = $1", order_id).Scan(&count_o)
	if err != nil {
		myLog.Log.Errorf(err.Error())
	}
	fmt.Println(count_d, " ", count_o)
	// Проверяем, есть ли запись
	if count_o > 0 && count_d > 0 {
		fmt.Printf("Запись с id %d существует.\n", order_id)
		query_add_delivery :=
			`INSERT INTO  delivery (order_id, delivery_id, updated_at)
		VALUES ($1, $2, $3)`
		err = db.Connection.QueryRow(query_add_delivery, strconv.Itoa(order_id), strconv.Itoa(delivery_man_id), time.Now().Format("2006-01-02 15:04:05")).Err()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Запись с id %d не найдена.\n", order_id)
		return myErrors.ErrNotFoundOrder
	}

	err = db.UpdateStatus(order_id, "delivery")

	var count int
	var orderIDs []int
	err = db.Connection.QueryRow("SELECT COUNT(*) AS count, ARRAY_AGG(order_id) AS order_ids FROM delivery WHERE delivery_id = $1", delivery_man_id).Scan(&count, &orderIDs)
	if err != nil {
		myLog.Log.Errorf(err.Error())
	}

	if count_o >= 5 {
		for i := 0; i < len(orderIDs); i++ {
			db.UpdateStatus(orderIDs[i], "delivery")
			myLog.Log.Debugf("The order with ID has been sent for delivery: ", order_id)
		}

	}

	return nil
}

func (db *Postgres) AddPayment(payment models.Payment) (int, error) {
	query_add_payment := `
        WITH insert_return AS (
            INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, custom_fee)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	var id_payment string
	err := db.Connection.QueryRow(query_add_payment, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		strconv.Itoa(payment.Amount), strconv.Itoa(payment.PaymentDT), payment.Bank, strconv.Itoa(payment.DeliveryCost),
		strconv.Itoa(payment.CustomFee)).Scan(&id_payment)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(id_payment)
	if err != nil {
		myLog.Log.Errorf("Invalid id Payment")
	}
	return id, nil
}

func (db *Postgres) AddItemsWithCategory(items []models.Item) ([]int, error) {
	check_category := ` 
        SELECT id 
		FROM category_item 
		WHERE category = $1
		`

	query_add_items := `
        WITH insert_return AS (
            INSERT INTO item (track_number, category_id, price, name, size, total_price, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            RETURNING id
        )
        SELECT id FROM insert_return
    `
	query_add_cat := `
	WITH insert_return AS (
		INSERT INTO category_item (category)
		VALUES ($1)
		RETURNING id
	)
	SELECT id FROM insert_return
`
	var id_items []int
	for i := 0; i < len(items); i++ {
		var id_category int = 0

		err := db.Connection.QueryRow(check_category, items[i].Category.CategoryName).Scan(&id_category)

		if err == nil || id_category == 0 {
			myLog.Log.Debugf("This product category already exists, id: ", id_category)
		} else {
			err := db.Connection.QueryRow(query_add_cat, items[i].Category.CategoryName).Scan(&id_category)
			if err != nil {
				myLog.Log.Errorf("Error CreateCategory: %v", err.Error())
				return id_items, err
			}
		}

		var id int
		err = db.Connection.QueryRow(query_add_items, items[i].TrackNumber, id_category, strconv.Itoa(items[i].Price),
			items[i].Name, items[i].Size, strconv.Itoa(items[i].TotalPrice),
			items[i].Brand, strconv.Itoa(items[i].Status)).Scan(&id)
		if err != nil {
			myLog.Log.Errorf("Error CreateItems: %v", err.Error())
			return id_items, err
		}
		myLog.Log.Debugf("Insert items %v", id)
		//id_, err := strconv.Atoi(id)
		id_items = append(id_items, id)
	}
	return id_items, nil
}

func (db *Postgres) AddOrder(order models.Order, items []int, payment int) (int, error) {
	query_add_order := `
	WITH insert_return AS (
	INSERT INTO orders (payment_id, items_id, locale, delivery_service, date_created)
	VALUES ($1, $2::int[], $3, $4, $5)
	RETURNING id
        )
        SELECT id FROM insert_return
`
	var id int
	err := db.Connection.QueryRow(query_add_order, payment, pq.Array(items), order.Locale,
		order.DeliveryService, order.DateCreated).Scan(&id)
	return id, err
}

func (db *Postgres) AddOrderStatus(order_id int) error {
	query_add_order := `
	INSERT INTO order_status (order_id, status, updated_at)
	VALUES ($1, $2, $3)
	`
	err := db.Connection.QueryRow(query_add_order, strconv.Itoa(order_id), "create", time.Now().Format("2006-01-02 15:04:05")).Err()
	if err != nil {
		myLog.Log.Errorf("Error Create OrderStatus: %v", err.Error())
	}
	return err
}

func (db *Postgres) UpdateStatus(order_id int, status string) error {
	query_update_order_status := `
	UPDATE order_status 
	SET status = $1, updated_at = $2 
	WHERE order_id = $3
	`

	_, err := db.Connection.Exec(query_update_order_status, status, time.Now(), order_id)
	if err != nil {
		myLog.Log.Errorf("Error UpdateStatus: %+v", err.Error())
	}
	return err
}

// добавить категорию и что-то еще
func (db *Postgres) GetOrder(order_id string) (models.Order, error) {
	myLog.Log.Debugf("Go to db in GetOrder with id: %+v", order_id)
	var order models.Order
	var payment models.Payment
	query_get_order := `
    SELECT o.items, o.locale, o.delivery_service, o.date_created,
        p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.custom_fee
    FROM 
        orders o
    LEFT JOIN 
        payment p ON o.payment = p.id
    WHERE 
        o.order_uid = $1`
	rows, err := db.Connection.Query(query_get_order, order_id)
	if err != nil {
		myLog.Log.Errorf("Error GetOrder: %v", err.Error())
		return models.Order{}, err
	}
	defer rows.Close()
	//itemMap := make(map[string]models.Item)
	var id_items []string
	if !rows.Next() { // Проверка на наличие записи
		myLog.Log.Errorf("No order found with uuid: %v", order_id)
		return models.Order{}, myErrors.ErrNotFoundOrder
	}
	for rows.Next() {
		//var item models.Item
		err = rows.Scan(
			pq.Array(&id_items),
			&order.Locale,
			&order.DeliveryService,
			&order.DateCreated,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.CustomFee,
		)
		if err != nil {
			myLog.Log.Errorf("Error scanning row: %v", err.Error())
			return models.Order{}, err
		}
	}

	query_get_item :=
		`SELECT id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	WHERE id = $1`
	var item models.Item
	var items []models.Item
	for i := 0; i < len(order.Items); i++ {
		err = db.Connection.QueryRow(query_get_item, order.Items[i]).Scan(&item.Id, &item.TrackNumber, &item.Price, &item.Name, &item.Size,
			&item.TotalPrice, &item.Brand, &item.Status)
		if err != nil {
			myLog.Log.Errorf("Error GetItems: %v", err.Error())
			return models.Order{}, err
		}
		items = append(items, item)
	}
	order.Items = items
	//order.DeliveryMan = delivery
	order.Payment = payment
	return order, nil
}

// func (db *Postgres) GetAllOrders() (map[int]models.Order, error) {
// 	result := make(map[int]models.Order)

// 	query := `SELECT o.order_id,
//       	o.items, o.locale, o.delivery_service, o.date_created,
//         d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
//         p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.custom_fee

// 	FROM  orders o
//     LEFT JOIN
//         delivery d ON o.delivery = d.id
//     LEFT JOIN
//         payment p ON o.payment = p.id`

// 	rows, err := db.Connection.Query(query)
// 	if err != nil {
// 		return result, err
// 	}
// 	defer rows.Close()
// 	var id_items []string
// 	var order models.Order
// 	var delivery models.DeliveryMan
// 	var payment models.Payment
// 	for rows.Next() {
// 		err = rows.Scan(
// 			&order.Id,
// 			pq.Array(&id_items),
// 			&order.Locale,
// 			&order.DeliveryService,
// 			&order.DateCreated,
// 			&delivery.Name,
// 			&delivery.Phone,
// 			&delivery.Zip,
// 			&delivery.City,
// 			&delivery.Address,
// 			&delivery.Region,
// 			&delivery.Email,
// 			&payment.Transaction,
// 			&payment.RequestID,
// 			&payment.Currency,
// 			&payment.Provider,
// 			&payment.Amount,
// 			&payment.PaymentDT,
// 			&payment.Bank,
// 			&payment.DeliveryCost,
// 			&payment.CustomFee,
// 		)
// 		if err != nil {
// 			myLog.Log.Errorf("Error scanning row: %v", err.Error())
// 			return result, err
// 		}
// 		query_get_item :=
// 			`SELECT id, track_number, price, name, size, total_price, brand, status
// 	WHERE id = $1`
// 		var item models.Item
// 		var items []models.Item
// 		for i := 0; i < len(order.Items); i++ {
// 			err = db.Connection.QueryRow(query_get_item, order.Items[i]).Scan(&item.Id, &item.TrackNumber, &item.Price, &item.Name, &item.Size,
// 				&item.TotalPrice, &item.Brand, &item.Status)
// 			if err != nil {
// 				myLog.Log.Errorf("Error GetItems: %v", err.Error())
// 				return result, err
// 			}
// 			items = append(items, item)
// 		}
// 		order.Items = items
// 		order.Payment = payment
// 		result[order.Id] = order
// 		fmt.Println(order.Id)
// 	}
// 	return result, nil
// }
