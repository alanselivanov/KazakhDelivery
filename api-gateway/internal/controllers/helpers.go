package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
	c.Abort()
}

func ParseUintParam(c *gin.Context, param string) (uint, error) {
	val := c.Param(param)
	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func ParsePaginationParams(c *gin.Context) (page int, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ = strconv.Atoi(pageStr)
	limit, _ = strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}
