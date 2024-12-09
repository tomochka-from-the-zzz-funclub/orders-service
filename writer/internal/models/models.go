package models

type Item struct {
	Id          int          `json:"id"`
	TrackNumber string       `json:"track_number"`
	Category    CategoryItem `json:"category"`
	Price       int          `json:"price"`
	Name        string       `json:"name"`
	Size        string       `json:"size"`
	TotalPrice  int          `json:"total_price"`
	Brand       string       `json:"brand"`
	Status      int          `json:"status"`
}

type DeliveryMan struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	CustomFee    int    `json:"custom_fee"`
}

type Order struct {
	Id int `json:"order_id"`
	//DeliveryMan     DeliveryMan `json:"delivery"`
	Payment         Payment `json:"payment"`
	Items           []Item  `json:"items"`
	Locale          string  `json:"locale"`
	DeliveryService string  `json:"delivery_service"`
	DateCreated     string  `json:"date_created"`
}

type CategoryItem struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type OrderStatus struct {
	Order      Order  `json:"order"`
	Status     string `json:"status"`
	Updated_at string `json:"updated_at"`
}

// type Caustom struct {
// 	Id         int
// 	FirstName  string
// 	LastName   string
// 	Email      string
// 	Phone      string
// 	Address    string
// 	City       string
// 	Region     string
// 	Zip        int
// 	Created_at string
// }
