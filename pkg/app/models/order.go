package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UUID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();unique"`
	ProductName string    `gorm:"not null;"`
	Quantity    int       `gorm:"not null;"`
	UserName    string    `gorm:"not null;"`
	CreatedAt   time.Time `gorm:"not null;"`
	UpdatedAt   time.Time `gorm:"not null;"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, Order *Order) (string, error)
	GetOrderByOrderUserName(ctx context.Context, OrderUserName string) (*Order, error)
}
