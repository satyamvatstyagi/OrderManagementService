package usecase

import (
	"context"
	"html"
	"strings"

	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/domain"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/models"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
)

type OrderUsecase struct {
	OrderRepository models.OrderRepository
}

func NewOrderUsecase(OrderRepository models.OrderRepository) domain.OrderUsecase {
	return &OrderUsecase{
		OrderRepository: OrderRepository,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, createOrderRequest *domain.CreateOrderRequest) (*domain.CreateOrderResponse, error) {
	// Check if Order already exists
	_, err := u.OrderRepository.GetOrderByOrderUserName(ctx, html.EscapeString(strings.TrimSpace(createOrderRequest.UserName)))
	if err == nil {
		return nil, cerr.NewCustomErrorWithCodeAndOrigin("Order already exists", cerr.InvalidRequestErrorCode, err)
	}

	// Create the Order
	Order := &models.Order{
		ProductName: html.EscapeString(strings.TrimSpace(createOrderRequest.ProductName)),
		Quantity:    createOrderRequest.Quantity,
		UserName:    html.EscapeString(strings.TrimSpace(createOrderRequest.UserName)),
	}

	OrderID, err := u.OrderRepository.CreateOrder(ctx, Order)
	if err != nil {
		return nil, err
	}

	return &domain.CreateOrderResponse{
		OrderID: OrderID,
	}, nil
}

func (u *OrderUsecase) GetOrderByOrderUserName(ctx context.Context, getOrderByOrderUserNameRequest *domain.GetOrderByOrderUserNameRequest) (*domain.GetOrderByOrderUserNameResponse, error) {
	Order, err := u.OrderRepository.GetOrderByOrderUserName(ctx, html.EscapeString(strings.TrimSpace(getOrderByOrderUserNameRequest.UserName)))
	if err != nil {
		return nil, err
	}

	return &domain.GetOrderByOrderUserNameResponse{
		OrderID:     Order.UUID.String(),
		ProductName: Order.ProductName,
		Quantity:    Order.Quantity,
	}, nil
}
