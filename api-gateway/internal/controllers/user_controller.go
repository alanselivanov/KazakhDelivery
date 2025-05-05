package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user "proto/user"
)

type UserController struct {
	client user.UserServiceClient
}

func NewUserController(serviceAddr string) *UserController {
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &UserController{
		client: user.NewUserServiceClient(conn),
	}
}

func (c *UserController) RegisterUser(ctx *gin.Context) {
	var req user.UserRequest
	if err := ctx.ShouldBindJSON(&req.User); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.client.RegisterUser(ctx, &req)
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, res.User)
}

func (c *UserController) AuthenticateUser(ctx *gin.Context) {
	var authReq user.AuthRequest
	if err := ctx.ShouldBindJSON(&authReq); err != nil {
		RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.client.AuthenticateUser(ctx, &authReq)
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if !res.Success {
		RespondWithError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   res.Token,
		"success": true,
	})
}

func (c *UserController) GetUserProfile(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.client.GetUserProfile(ctx, &user.UserID{Id: id})
	if err != nil {
		RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}
