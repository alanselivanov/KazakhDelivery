package handlers

import (
	"context"

	"inventory-service/internal/application"
	"inventory-service/internal/domain"
	"proto/inventory"
)

type CategoryHandler struct {
	inventory.UnimplementedInventoryServiceServer
	categoryUseCase *application.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase *application.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *inventory.CategoryRequest) (*inventory.CategoryResponse, error) {
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

func (h *CategoryHandler) GetCategory(ctx context.Context, req *inventory.CategoryID) (*inventory.CategoryResponse, error) {
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

func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *inventory.CategoryRequest) (*inventory.CategoryResponse, error) {
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

func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *inventory.CategoryID) (*inventory.Empty, error) {
	err := h.categoryUseCase.DeleteCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &inventory.Empty{}, nil
}

func (h *CategoryHandler) ListCategories(ctx context.Context, req *inventory.Empty) (*inventory.CategoryListResponse, error) {
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
