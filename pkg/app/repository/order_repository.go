package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/models"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/consts"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/instrumentation"
	"go.elastic.co/apm/v2"
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

	localTime := time.Now()
	order := models.Order{
		ProductName: Order.ProductName,
		Quantity:    Order.Quantity,
		UserName:    Order.UserName,
		CreatedAt:   localTime,
		UpdatedAt:   localTime,
	}

	statement := r.database.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&order)
	})
	instrument := instrumentation.InitGormAPM(ctx, "postgresql", statement)
	defer instrument.GetSpan().End()

	if err := r.database.Create(&order).Error; err != nil {
		// Check if err is of type *pgconn.PgError and error code is 23505, which is the error code for unique_violation
		if err, ok := err.(*pgconn.PgError); ok && err.Code == consts.UniqueViolation {
			apm.CaptureError(ctx, fmt.Errorf("db error: %s", err.Error())).Send()
			return "", cerr.NewCustomErrorWithCodeAndOrigin("Order already exists", cerr.InvalidRequestErrorCode, err)
		}
		apm.CaptureError(ctx, fmt.Errorf("db error: %s", err.Error())).Send()
		return "", err
	}

	return order.UUID.String(), nil
}

func (r *OrderRepository) GetOrderByOrderUserName(ctx context.Context, OrderUserName string) (*models.Order, error) {

	Order := &models.Order{}

	statement := r.database.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("user_name = ?", OrderUserName).First(&Order)
	})
	instrument := instrumentation.InitGormAPM(ctx, "postgresql", statement)
	defer instrument.GetSpan().End()

	if err := r.database.Where("user_name = ?", OrderUserName).First(&Order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			apm.CaptureError(ctx, fmt.Errorf("db error: %s", err.Error())).Send()
			return nil, cerr.NewCustomErrorWithCodeAndOrigin("Order not found", cerr.NotFoundErrorCode, err)
		}
		apm.CaptureError(ctx, fmt.Errorf("db error: %s", err.Error())).Send()
		return nil, err
	}
	return Order, nil
}
