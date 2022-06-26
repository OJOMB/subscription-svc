package app

import (
	"context"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
)

type Service interface {
	CreateUser(ctx context.Context, firstName, lastName, dob, email string) (*domain.User, error)
	CreateProduct(ctx context.Context, product domain.Product) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]*domain.Product, error)
	GetProductsWithVoucher(ctx context.Context, voucherCode string) ([]*domain.Product, error)
	GetProduct(ctx context.Context, productID string) (*domain.Product, error)
	CreateSubscriptionPlan(ctx context.Context, userID, productID, voucherCode string) (*domain.SubscriptionPlan, error)
	GetSubscriptionPlan(ctx context.Context, subscriptionPlanID string) (*domain.SubscriptionPlan, error)
	UpdateSubscriptionPlanStatus(ctx context.Context, planID, newStatus string) error
	CreateVoucher(ctx context.Context, voucher domain.Voucher, products []domain.Product) error
}
