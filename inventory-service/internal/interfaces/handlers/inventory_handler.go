package handlers

import (
	"context"

	"inventory-service/internal/application"
	"inventory-service/internal/domain"
	"proto/inventory"
)

type InventoryHandler struct {
	inventory.UnimplementedInventoryServiceServer
	productUseCase  *application.ProductUseCase
	categoryUseCase *application.CategoryUseCase
}

func NewInventoryHandler(productUseCase *application.ProductUseCase, categoryUseCase *application.CategoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		productUseCase:  productUseCase,
		categoryUseCase: categoryUseCase,
	}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *inventory.ProductRequest) (*inventory.ProductResponse, error) {
	product := domain.NewProduct(
		req.Product.Name,
		req.Product.Description,
		float64(req.Product.Price),
		int(req.Product.Stock),
		req.Product.CategoryId,
	)

	created, err := h.productUseCase.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          created.ID,
			Name:        created.Name,
			Description: created.Description,
			Price:       float32(created.Price),
			Stock:       int32(created.Stock),
			CategoryId:  created.CategoryID,
		},
	}, nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, req *inventory.ProductID) (*inventory.ProductResponse, error) {
	product, err := h.productUseCase.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
			CategoryId:  product.CategoryID,
		},
	}, nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *inventory.ProductRequest) (*inventory.ProductResponse, error) {
	product := &domain.Product{
		ID:          req.Product.Id,
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       float64(req.Product.Price),
		Stock:       int(req.Product.Stock),
		CategoryID:  req.Product.CategoryId,
	}

	updated, err := h.productUseCase.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, nil
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          updated.ID,
			Name:        updated.Name,
			Description: updated.Description,
			Price:       float32(updated.Price),
			Stock:       int32(updated.Stock),
			CategoryId:  updated.CategoryID,
		},
	}, nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, req *inventory.ProductID) (*inventory.Empty, error) {
	err := h.productUseCase.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &inventory.Empty{}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, req *inventory.ProductListRequest) (*inventory.ProductListResponse, error) {
	products, total, err := h.productUseCase.ListProducts(ctx, req.CategoryId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, err
	}

	var protoProducts []*inventory.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		})
	}

	return &inventory.ProductListResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

func (h *InventoryHandler) CreateCategory(ctx context.Context, req *inventory.CategoryRequest) (*inventory.CategoryResponse, error) {
	category := domain.NewCategory(
		req.Category.Name,
		req.Category.Description,
	)

	created, err := h.categoryUseCase.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return &inventory.CategoryResponse{
		Category: &inventory.Category{
			Id:          created.ID,
			Name:        created.Name,
			Description: created.Description,
		},
	}, nil
}

func (h *InventoryHandler) GetCategory(ctx context.Context, req *inventory.CategoryID) (*inventory.CategoryResponse, error) {
	category, err := h.categoryUseCase.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	return &inventory.CategoryResponse{
		Category: &inventory.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}, nil
}

func (h *InventoryHandler) UpdateCategory(ctx context.Context, req *inventory.CategoryRequest) (*inventory.CategoryResponse, error) {
	category := &domain.Category{
		ID:          req.Category.Id,
		Name:        req.Category.Name,
		Description: req.Category.Description,
	}

	updated, err := h.categoryUseCase.UpdateCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, nil
	}

	return &inventory.CategoryResponse{
		Category: &inventory.Category{
			Id:          updated.ID,
			Name:        updated.Name,
			Description: updated.Description,
		},
	}, nil
}

func (h *InventoryHandler) DeleteCategory(ctx context.Context, req *inventory.CategoryID) (*inventory.Empty, error) {
	err := h.categoryUseCase.DeleteCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &inventory.Empty{}, nil
}

func (h *InventoryHandler) ListCategories(ctx context.Context, req *inventory.Empty) (*inventory.CategoryListResponse, error) {
	categories, err := h.categoryUseCase.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var protoCategories []*inventory.Category
	for _, c := range categories {
		protoCategories = append(protoCategories, &inventory.Category{
			Id:          c.ID,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return &inventory.CategoryListResponse{
		Categories: protoCategories,
	}, nil
}
