package domain

import "context"

type Repo interface {
	GetProduct(ctx context.Context, productID string) (*Product, error)
	GetProducts(ctx context.Context) ([]*Product, error)
	GetProductsWithVoucher(ctx context.Context, voucherCode string) ([]*Product, error)
	CreateSubscriptionPlan(ctx context.Context, plan SubscriptionPlan) error
	GetSubscriptionPlan(ctx context.Context, planID string) (*SubscriptionPlan, error)
	PauseSubscriptionPlan(ctx context.Context, planID, pauseID string) error
	ResumeSubscriptionPlan(ctx context.Context, planID string) error
	CancelSubscriptionPlan(ctx context.Context, planID string) error
	GetCurrentSubscriptionPlan(ctx context.Context, userID string) (*SubscriptionPlan, error)
	GetVoucher(ctx context.Context, voucherCode string) (*Voucher, error)
	GetVoucherProduct(ctx context.Context, voucherCode, productID string) (*VoucherProduct, error)
}
