package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/gorm"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/models"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/consts"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/instrumentation"
	"go.elastic.co/apm/module/apmgorm/v2"
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
	span, ctx := instrumentation.TraceAPMRequest(ctx, "CreateOrder", consts.SpanTypeQueryExecution)
	defer span.End()
	db := apmgorm.WithContext(ctx, r.database)
	localTime := time.Now()
	order := models.Order{
		ProductName: Order.ProductName,
		Quantity:    Order.Quantity,
		UserName:    Order.UserName,
		CreatedAt:   localTime,
		UpdatedAt:   localTime,
	}

	if err := db.Create(&order).Error; err != nil {
		// Check if err is of type *pgconn.PgError and error code is 23505, which is the error code for unique_violation
		if err, ok := err.(*pgconn.PgError); ok && err.Code == consts.UniqueViolation {
			return "", cerr.NewCustomErrorWithCodeAndOrigin("Order already exists", cerr.InvalidRequestErrorCode, err)
		}
		return "", err
	}

	return order.UUID.String(), nil
}

func (r *OrderRepository) GetOrderByOrderUserName(ctx context.Context, OrderUserName string) (*models.Order, error) {
	span, ctx := instrumentation.TraceAPMRequest(ctx, "GetOrderByOrderUserName", consts.SpanTypeQueryExecution)
	defer span.End()
	db := apmgorm.WithContext(ctx, r.database)
	Order := &models.Order{}
	if err := db.Where("user_name = ?", OrderUserName).First(&Order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cerr.NewCustomErrorWithCodeAndOrigin("Order not found", cerr.NotFoundErrorCode, err)
		}
		return nil, err
	}
	return Order, nil
}
