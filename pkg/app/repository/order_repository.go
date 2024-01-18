package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/models"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/consts"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/instrumentation"
	"gorm.io/gorm"
)

type OrderRepository struct {
	database *gorm.DB
}

func NewOrderRepository(database *gorm.DB) models.OrderRepository {
	return &OrderRepository{
		database: database,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, Order *models.Order) (string, error) {
	_, _ = instrumentation.TraceAPMRequest(ctx, "CreateOrder", consts.SpanTypeQueryExecution)
	localTime := time.Now()
	order := models.Order{
		ProductName: Order.ProductName,
		Quantity:    Order.Quantity,
		UserName:    Order.UserName,
		CreatedAt:   localTime,
		UpdatedAt:   localTime,
	}

	if err := r.database.Create(&order).Error; err != nil {
		// Check if err is of type *pgconn.PgError and error code is 23505, which is the error code for unique_violation
		if err, ok := err.(*pgconn.PgError); ok && err.Code == consts.UniqueViolation {
			return "", cerr.NewCustomErrorWithCodeAndOrigin("Order already exists", cerr.InvalidRequestErrorCode, err)
		}
		return "", err
	}

	return order.UUID.String(), nil
}

func (r *OrderRepository) GetOrderByOrderUserName(ctx context.Context, OrderUserName string) (*models.Order, error) {
	_, _ = instrumentation.TraceAPMRequest(ctx, "GetOrderByOrderUserName", consts.SpanTypeQueryExecution)
	Order := &models.Order{}
	if err := r.database.Where("user_name = ?", OrderUserName).First(&Order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cerr.NewCustomErrorWithCodeAndOrigin("Order not found", cerr.NotFoundErrorCode, err)
		}
		return nil, err
	}
	return Order, nil
}
