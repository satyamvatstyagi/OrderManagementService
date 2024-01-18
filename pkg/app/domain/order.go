package domain

import "context"

type OrderUsecase interface {
	CreateOrder(ctx context.Context, createOrderRequest *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrderByOrderUserName(ctx context.Context, getOrderByOrderUserNameRequest *GetOrderByOrderUserNameRequest) (*GetOrderByOrderUserNameResponse, error)
}

type CreateOrderRequest struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UserName    string `json:"user_name"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}

type GetOrderByOrderUserNameRequest struct {
	UserName string `uri:"username" binding:"required"`
}

type GetOrderByOrderUserNameResponse struct {
	OrderID     string `json:"order_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
