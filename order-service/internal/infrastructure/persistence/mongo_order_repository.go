package persistence

import (
	"context"
	"errors"
	"time"

	"order-service/internal/domain"
	"order-service/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoOrderRepository struct {
	db *database.MongoDBConnector
}

func NewMongoOrderRepository(db *database.MongoDBConnector) *mongoOrderRepository {
	return &mongoOrderRepository{db: db}
}

func (r *mongoOrderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	itemDTOs := make([]database.OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		itemDTOs[i] = database.OrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	orderDTO := &database.OrderDTO{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     itemDTOs,
		Total:     order.Total,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

	_, err := r.db.OrderCollection().InsertOne(ctx, orderDTO)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *mongoOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var orderDTO database.OrderDTO

	filter := bson.M{"_id": id}
	err := r.db.OrderCollection().FindOne(ctx, filter).Decode(&orderDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	orderItems := make([]domain.OrderItem, len(orderDTO.Items))
	for i, item := range orderDTO.Items {
		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &domain.Order{
		ID:        orderDTO.ID,
		UserID:    orderDTO.UserID,
		Items:     orderItems,
		Total:     orderDTO.Total,
		Status:    domain.OrderStatus(orderDTO.Status),
		CreatedAt: orderDTO.CreatedAt,
		UpdatedAt: orderDTO.UpdatedAt,
	}, nil
}

func (r *mongoOrderRepository) Update(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	itemDTOs := make([]database.OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		itemDTOs[i] = database.OrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	filter := bson.M{"_id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"user_id":    order.UserID,
			"items":      itemDTOs,
			"total":      order.Total,
			"status":     string(order.Status),
			"updated_at": time.Now(),
		},
	}

	result, err := r.db.OrderCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, nil
	}

	return order, nil
}

func (r *mongoOrderRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.db.OrderCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orderDTOs []database.OrderDTO
	if err := cursor.All(ctx, &orderDTOs); err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(orderDTOs))
	for i, dto := range orderDTOs {
		orderItems := make([]domain.OrderItem, len(dto.Items))
		for j, item := range dto.Items {
			orderItems[j] = domain.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}

		orders[i] = &domain.Order{
			ID:        dto.ID,
			UserID:    dto.UserID,
			Items:     orderItems,
			Total:     dto.Total,
			Status:    domain.OrderStatus(dto.Status),
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		}
	}

	return orders, nil
}
