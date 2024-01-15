package domain

type OrderUsecase interface {
	CreateOrder(createOrderRequest *CreateOrderRequest) (*CreateOrderResponse, error)
}

type CreateOrderRequest struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UserName    string `json:"user_name"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}
