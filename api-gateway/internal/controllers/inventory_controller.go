package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventory "proto/inventory"
)

type InventoryController struct {
	client inventory.InventoryServiceClient
}

func NewInventoryController(serviceAddr string) *InventoryController {
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &InventoryController{
		client: inventory.NewInventoryServiceClient(conn),
	}
}

func (c *InventoryController) CreateProduct(ctx *gin.Context) {
	var product inventory.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.client.CreateProduct(ctx, &inventory.ProductRequest{Product: &product})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, res.Product)
}

func (c *InventoryController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.client.GetProduct(ctx, &inventory.ProductID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if res.Product == nil {
		RespondWithError(ctx, http.StatusNotFound, "product not found")
		return
	}

	ctx.JSON(http.StatusOK, res.Product)
}

func (c *InventoryController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product inventory.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	product.Id = id

	res, err := c.client.UpdateProduct(ctx, &inventory.ProductRequest{Product: &product})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Product)
}

func (c *InventoryController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := c.client.DeleteProduct(ctx, &inventory.ProductID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *InventoryController) ListProducts(ctx *gin.Context) {
	categoryID := ctx.Query("category_id")
	page, limit := ParsePaginationParams(ctx)

	res, err := c.client.ListProducts(ctx, &inventory.ProductListRequest{
		CategoryId: categoryID,
		Page:       int32(page),
		Limit:      int32(limit),
	})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": res.Products,
		"total":    res.Total,
	})
}

func (c *InventoryController) CreateCategory(ctx *gin.Context) {
	var category inventory.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.client.CreateCategory(ctx, &inventory.CategoryRequest{Category: &category})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, res.Category)
}

func (c *InventoryController) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.client.GetCategory(ctx, &inventory.CategoryID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if res.Category == nil {
		RespondWithError(ctx, http.StatusNotFound, "category not found")
		return
	}

	ctx.JSON(http.StatusOK, res.Category)
}

func (c *InventoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category inventory.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	category.Id = id

	res, err := c.client.UpdateCategory(ctx, &inventory.CategoryRequest{Category: &category})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Category)
}

func (c *InventoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := c.client.DeleteCategory(ctx, &inventory.CategoryID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *InventoryController) ListCategories(ctx *gin.Context) {
	res, err := c.client.ListCategories(ctx, &inventory.Empty{})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Categories)
}
