package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	order "proto/order"
)

type OrderController struct {
	client order.OrderServiceClient
}

func NewOrderController(serviceAddr string) *OrderController {
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &OrderController{
		client: order.NewOrderServiceClient(conn),
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req order.OrderRequest
	if err := ctx.ShouldBindJSON(&req.Order); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.client.CreateOrder(ctx, &req)
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, res.Order)
}

func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.client.GetOrder(ctx, &order.OrderID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if res.Order == nil {
		RespondWithError(ctx, http.StatusNotFound, "order not found")
		return
	}

	ctx.JSON(http.StatusOK, res.Order)
}

func (c *OrderController) UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var req order.OrderRequest
	if err := ctx.ShouldBindJSON(&req.Order); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.Order.Id = id

	res, err := c.client.UpdateOrder(ctx, &req)
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Order)
}

func (c *OrderController) ListOrders(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		RespondWithError(ctx, http.StatusBadRequest, "user_id is required")
		return
	}

	res, err := c.client.ListOrders(ctx, &order.UserID{Id: userID})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Orders)
}
