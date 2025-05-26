package orderservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Order struct {
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	Status       OrderStatus  `json:"status"`
	Items        []OrderItem  `json:"items"`
	TotalAmount  float64      `json:"total_amount"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	PaymentInfo  PaymentInfo  `json:"payment_info"`
	ShippingInfo ShippingInfo `json:"shipping_info"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type PaymentInfo struct {
	Method        string    `json:"method"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
}

type ShippingInfo struct {
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	TrackingID  string    `json:"tracking_id,omitempty"`
	ShippedAt   time.Time `json:"shipped_at,omitempty"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	FindByID(ctx context.Context, id string) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) ([]*Order, error)
}

type InventoryService interface {
	CheckStock(ctx context.Context, productID string, quantity int) (bool, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
	ReleaseStock(ctx context.Context, productID string, quantity int) error
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64, method string) (string, error)
	RefundPayment(ctx context.Context, transactionID string) error
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *Order) (*Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id string) (*Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *Order) (*Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Order), args.Error(1)
}

type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	args := m.Called(ctx, productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockInventoryService) ReserveStock(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) ReleaseStock(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, orderID string, amount float64, method string) (string, error) {
	args := m.Called(ctx, orderID, amount, method)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentService) RefundPayment(ctx context.Context, transactionID string) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}

type OrderProcessor struct {
	orderRepo     OrderRepository
	inventoryServ InventoryService
	paymentServ   PaymentService
}

func NewOrderProcessor(repo OrderRepository, invServ InventoryService, payServ PaymentService) *OrderProcessor {
	return &OrderProcessor{
		orderRepo:     repo,
		inventoryServ: invServ,
		paymentServ:   payServ,
	}
}

func (p *OrderProcessor) ProcessOrder(ctx context.Context, order *Order) (*Order, error) {
	if order == nil || order.UserID == "" || len(order.Items) == 0 {
		return nil, errors.New("invalid order data")
	}

	reservedProducts := make(map[string]int)

	for _, item := range order.Items {
		inStock, err := p.inventoryServ.CheckStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			for prodID, qty := range reservedProducts {
				_ = p.inventoryServ.ReleaseStock(ctx, prodID, qty)
			}
			return nil, errors.New("failed to check inventory")
		}

		if !inStock {
			for prodID, qty := range reservedProducts {
				_ = p.inventoryServ.ReleaseStock(ctx, prodID, qty)
			}
			return nil, errors.New("insufficient inventory for product: " + item.ProductID)
		}

		err = p.inventoryServ.ReserveStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			for prodID, qty := range reservedProducts {
				_ = p.inventoryServ.ReleaseStock(ctx, prodID, qty)
			}
			return nil, errors.New("failed to reserve inventory")
		}

		reservedProducts[item.ProductID] = item.Quantity
	}

	var totalAmount float64
	for _, item := range order.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = totalAmount

	order.Status = OrderStatusPending

	savedOrder, err := p.orderRepo.Create(ctx, order)
	if err != nil {
		for prodID, qty := range reservedProducts {
			_ = p.inventoryServ.ReleaseStock(ctx, prodID, qty)
		}
		return nil, errors.New("failed to save order")
	}

	transactionID, err := p.paymentServ.ProcessPayment(ctx, savedOrder.ID, savedOrder.TotalAmount, savedOrder.PaymentInfo.Method)
	if err != nil {
		savedOrder.Status = OrderStatusCancelled
		_, _ = p.orderRepo.Update(ctx, savedOrder)
		for prodID, qty := range reservedProducts {
			_ = p.inventoryServ.ReleaseStock(ctx, prodID, qty)
		}
		return nil, errors.New("payment processing failed")
	}

	savedOrder.PaymentInfo.TransactionID = transactionID
	savedOrder.PaymentInfo.Status = "completed"
	savedOrder.PaymentInfo.PaidAt = time.Now()

	savedOrder.Status = OrderStatusProcessing
	updatedOrder, err := p.orderRepo.Update(ctx, savedOrder)
	if err != nil {
		return nil, errors.New("failed to update order status")
	}

	return updatedOrder, nil
}

func (p *OrderProcessor) CancelOrder(ctx context.Context, orderID string) error {
	order, err := p.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status == OrderStatusShipped || order.Status == OrderStatusDelivered {
		return errors.New("cannot cancel order that has been shipped or delivered")
	}

	if order.PaymentInfo.TransactionID != "" {
		err = p.paymentServ.RefundPayment(ctx, order.PaymentInfo.TransactionID)
		if err != nil {
			return errors.New("failed to refund payment")
		}
	}

	for _, item := range order.Items {
		err = p.inventoryServ.ReleaseStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return errors.New("failed to release inventory")
		}
	}

	order.Status = OrderStatusCancelled
	_, err = p.orderRepo.Update(ctx, order)
	if err != nil {
		return errors.New("failed to update order status")
	}

	return nil
}

func TestOrderProcessor_ProcessOrder_Success(t *testing.T) {
	ctx := context.Background()
	mockOrderRepo := new(MockOrderRepository)
	mockInventoryServ := new(MockInventoryService)
	mockPaymentServ := new(MockPaymentService)

	processor := NewOrderProcessor(mockOrderRepo, mockInventoryServ, mockPaymentServ)

	order := &Order{
		ID:     "order-123",
		UserID: "user-123",
		Items: []OrderItem{
			{
				ProductID: "prod-1",
				Quantity:  2,
				Price:     10.00,
			},
			{
				ProductID: "prod-2",
				Quantity:  1,
				Price:     15.00,
			},
		},
		PaymentInfo: PaymentInfo{
			Method: "credit_card",
		},
		ShippingInfo: ShippingInfo{
			Address:    "123 Test St",
			City:       "Test City",
			PostalCode: "12345",
			Country:    "Test Country",
		},
	}

	expectedTotalAmount := 35.00
	expectedTransactionID := "tx-123456"

	mockInventoryServ.On("CheckStock", ctx, "prod-1", 2).Return(true, nil)
	mockInventoryServ.On("CheckStock", ctx, "prod-2", 1).Return(true, nil)

	mockInventoryServ.On("ReserveStock", ctx, "prod-1", 2).Return(nil)
	mockInventoryServ.On("ReserveStock", ctx, "prod-2", 1).Return(nil)

	mockOrderRepo.On("Create", ctx, mock.AnythingOfType("*orderservice_test.Order")).
		Run(func(args mock.Arguments) {
			savedOrder := args.Get(1).(*Order)
			assert.Equal(t, OrderStatusPending, savedOrder.Status)
			assert.Equal(t, expectedTotalAmount, savedOrder.TotalAmount)
		}).
		Return(order, nil)

	mockPaymentServ.On("ProcessPayment", ctx, order.ID, expectedTotalAmount, "credit_card").
		Return(expectedTransactionID, nil)

	mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*orderservice_test.Order")).
		Run(func(args mock.Arguments) {
			updatedOrder := args.Get(1).(*Order)
			assert.Equal(t, OrderStatusProcessing, updatedOrder.Status)
			assert.Equal(t, expectedTransactionID, updatedOrder.PaymentInfo.TransactionID)
			assert.Equal(t, "completed", updatedOrder.PaymentInfo.Status)
			assert.NotZero(t, updatedOrder.PaymentInfo.PaidAt)
		}).
		Return(order, nil)

	result, err := processor.ProcessOrder(ctx, order)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, order.ID, result.ID)
	mockOrderRepo.AssertExpectations(t)
	mockInventoryServ.AssertExpectations(t)
	mockPaymentServ.AssertExpectations(t)
}

func TestOrderProcessor_ProcessOrder_InsufficientInventory(t *testing.T) {
	ctx := context.Background()
	mockOrderRepo := new(MockOrderRepository)
	mockInventoryServ := new(MockInventoryService)
	mockPaymentServ := new(MockPaymentService)

	processor := NewOrderProcessor(mockOrderRepo, mockInventoryServ, mockPaymentServ)

	order := &Order{
		ID:     "order-123",
		UserID: "user-123",
		Items: []OrderItem{
			{
				ProductID: "prod-1",
				Quantity:  2,
				Price:     10.00,
			},
			{
				ProductID: "prod-out-of-stock",
				Quantity:  1,
				Price:     15.00,
			},
		},
	}

	mockInventoryServ.On("CheckStock", ctx, "prod-1", 2).Return(true, nil)
	mockInventoryServ.On("ReserveStock", ctx, "prod-1", 2).Return(nil)

	mockInventoryServ.On("CheckStock", ctx, "prod-out-of-stock", 1).Return(false, nil)

	mockInventoryServ.On("ReleaseStock", ctx, "prod-1", 2).Return(nil)

	result, err := processor.ProcessOrder(ctx, order)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "insufficient inventory")
	mockOrderRepo.AssertNotCalled(t, "Create")
	mockPaymentServ.AssertNotCalled(t, "ProcessPayment")
	mockInventoryServ.AssertExpectations(t)
}

func TestOrderProcessor_CancelOrder_Success(t *testing.T) {
	ctx := context.Background()
	mockOrderRepo := new(MockOrderRepository)
	mockInventoryServ := new(MockInventoryService)
	mockPaymentServ := new(MockPaymentService)

	processor := NewOrderProcessor(mockOrderRepo, mockInventoryServ, mockPaymentServ)

	order := &Order{
		ID:     "order-to-cancel",
		UserID: "user-123",
		Status: OrderStatusProcessing,
		Items: []OrderItem{
			{
				ProductID: "prod-1",
				Quantity:  2,
				Price:     10.00,
			},
		},
		PaymentInfo: PaymentInfo{
			TransactionID: "tx-123456",
			Method:        "credit_card",
			Status:        "completed",
		},
	}

	mockOrderRepo.On("FindByID", ctx, "order-to-cancel").Return(order, nil)
	mockPaymentServ.On("RefundPayment", ctx, "tx-123456").Return(nil)
	mockInventoryServ.On("ReleaseStock", ctx, "prod-1", 2).Return(nil)
	mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*orderservice_test.Order")).
		Run(func(args mock.Arguments) {
			updatedOrder := args.Get(1).(*Order)
			assert.Equal(t, OrderStatusCancelled, updatedOrder.Status)
		}).
		Return(order, nil)

	err := processor.CancelOrder(ctx, "order-to-cancel")

	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockInventoryServ.AssertExpectations(t)
	mockPaymentServ.AssertExpectations(t)
}

func TestOrderProcessor_CancelOrder_AlreadyShipped(t *testing.T) {
	ctx := context.Background()
	mockOrderRepo := new(MockOrderRepository)
	mockInventoryServ := new(MockInventoryService)
	mockPaymentServ := new(MockPaymentService)

	processor := NewOrderProcessor(mockOrderRepo, mockInventoryServ, mockPaymentServ)

	order := &Order{
		ID:     "shipped-order",
		UserID: "user-123",
		Status: OrderStatusShipped,
		Items: []OrderItem{
			{
				ProductID: "prod-1",
				Quantity:  2,
				Price:     10.00,
			},
		},
	}

	mockOrderRepo.On("FindByID", ctx, "shipped-order").Return(order, nil)

	err := processor.CancelOrder(ctx, "shipped-order")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel order that has been shipped")
	mockPaymentServ.AssertNotCalled(t, "RefundPayment")
	mockInventoryServ.AssertNotCalled(t, "ReleaseStock")
	mockOrderRepo.AssertNotCalled(t, "Update")
}
