package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/models"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/cerr"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/common/consts"
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

func (r *OrderRepository) CreateOrder(Order *models.Order) (string, error) {
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

func (r *OrderRepository) GetOrderByOrderUserName(OrderUserName string) (*models.Order, error) {
	Order := &models.Order{}
	if err := r.database.Where("user_name = ?", OrderUserName).First(&Order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cerr.NewCustomErrorWithCodeAndOrigin("Order not found", cerr.NotFoundErrorCode, err)
		}
		return nil, err
	}
	return Order, nil
}
