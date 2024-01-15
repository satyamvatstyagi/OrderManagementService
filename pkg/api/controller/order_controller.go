package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/domain"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
)

type OrderController struct {
	OrderUsecase domain.OrderUsecase
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var createOrderRequest domain.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&createOrderRequest); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, cerr.NewCustomErrorWithCodeAndOrigin("Invalid request", cerr.InvalidRequestErrorCode, err))
		return
	}

	createOrderResponse, err := c.OrderUsecase.CreateOrder(&createOrderRequest)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, cerr.NewCustomErrorWithCodeAndOrigin("Internal server error", cerr.InternalServerErrorCode, err))
		return
	}

	ctx.JSON(http.StatusCreated, createOrderResponse)
}

func (c *OrderController) GetOrderByOrderUserName(ctx *gin.Context) {
	var getOrderByOrderUserNameRequest domain.GetOrderByOrderUserNameRequest
	if err := ctx.ShouldBindUri(&getOrderByOrderUserNameRequest); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, cerr.NewCustomErrorWithCodeAndOrigin("Invalid request", cerr.InvalidRequestErrorCode, err))
		return
	}

	getOrderByOrderUserNameResponse, err := c.OrderUsecase.GetOrderByOrderUserName(&getOrderByOrderUserNameRequest)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, cerr.NewCustomErrorWithCodeAndOrigin("Internal server error", cerr.InternalServerErrorCode, err))
		return
	}

	ctx.JSON(http.StatusOK, getOrderByOrderUserNameResponse)
}
