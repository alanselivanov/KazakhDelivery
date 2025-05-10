package messaging

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderCreatedEvent struct {
	OrderID   string      `json:"order_id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Timestamp int64       `json:"timestamp"`
}

type StockUpdateResult struct {
	ProductID string `json:"product_id"`
	Success   bool   `json:"success"`
	Quantity  int    `json:"quantity"`
	Error     string `json:"error,omitempty"`
}

const (
	SubjectOrderCreated = "order.created"
	SubjectDeadLetter   = "dead.letter.queue"
)
