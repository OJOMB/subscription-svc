package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	componentService = "Service"
	svcMethod        = "method"
)

type Service struct {
	logger  *logrus.Entry
	repo    Repo
	idGen   IDGenerator
	taxRate decimal.Decimal
}

func NewService(logger *logrus.Logger, repo Repo, idGen IDGenerator, taxRate decimal.Decimal) *Service {
	return &Service{
		logger:  logger.WithField("component", componentService),
		repo:    repo,
		idGen:   idGen,
		taxRate: taxRate,
	}
}

func (svc *Service) CreateProduct(ctx context.Context, p Product) (*Product, error) {
	return nil, nil
}

func (svc *Service) GetProducts(ctx context.Context) ([]*Product, error) {
	products, err := svc.repo.GetProducts(ctx)
	if err != nil {
		svc.logger.WithField(svcMethod, "GetProducts").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to retrieve products")
	}

	return products, nil
}

func (svc *Service) GetProductsWithVoucher(ctx context.Context, voucherCode string) ([]*Product, error) {
	voucher, err := svc.repo.GetVoucher(ctx, voucherCode)
	if voucher == nil {
		return nil, fmt.Errorf("resource not found - no such voucher with code '%s' exists", voucherCode)
	} else if err != nil {
		svc.logger.WithField(svcMethod, "GetProductsWithVoucher").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to retrieve voucher with code '%s'", voucherCode)
	}

	if err := svc.validateVoucher(voucher); err != nil {
		return nil, fmt.Errorf("bad input data - %v", err)
	}

	products, err := svc.repo.GetProductsWithVoucher(ctx, voucher.Code)
	if err != nil {
		svc.logger.WithField(svcMethod, "GetProductsWithVoucher").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to retrieve products applicable to voucher with code '%s'", voucherCode)
	}

	return products, nil
}

func (svc *Service) GetProduct(ctx context.Context, productID string) (*Product, error) {
	if productID == "" {
		return nil, fmt.Errorf("bad input data - productID must not be empty")
	}

	product, err := svc.repo.GetProduct(ctx, productID)
	if err != nil {
		svc.logger.WithField(svcMethod, "GetProduct").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to retrieve product with ID '%s'", productID)
	}

	return product, nil
}

func (svc *Service) CreateSubscriptionPlan(ctx context.Context, userID, productID, voucherCode string) (*SubscriptionPlan, error) {
	if userID == "" || productID == "" {
		return nil, fmt.Errorf("bad input data - valid userID and productID must be provided")
	}

	// user can only have one plan at a time so need to check there isn't an already existing valid plan on record
	existingPlan, err := svc.repo.GetCurrentSubscriptionPlan(ctx, userID)
	if existingPlan != nil {
		return nil, fmt.Errorf("bad input data - user already has a subscription plan in place with ID '%s'", existingPlan.ID)
	} else if err != nil {
		svc.logger.WithField(svcMethod, "CreateSubscriptionPlan").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to check if user with ID '%s' has an existing plan", userID)
	}

	product, err := svc.repo.GetProduct(ctx, productID)
	if product == nil {
		return nil, fmt.Errorf("resource not found - no such product with ID '%s' exists", productID)
	} else if err != nil {
		return nil, fmt.Errorf("encountered error - failed to retrieve product with ID '%s'", productID)
	}

	var discount = decimal.NewFromInt(0)
	if voucherCode != "" {
		voucher, err := svc.repo.GetVoucher(ctx, voucherCode)
		if err != nil {
			svc.logger.WithField(svcMethod, "CreateSubscriptionPlan").WithError(err).Error("repo error")
			return nil, fmt.Errorf("encountered error - failed to retrieve voucher with code '%s'", voucherCode)
		}

		if err := svc.validateVoucher(voucher); err != nil {
			return nil, fmt.Errorf("bad input data - %v", err)
		}

		isApplicable, err := svc.checkVoucherAppliesToProduct(ctx, voucherCode, productID)
		if !isApplicable {
			return nil, fmt.Errorf("bad input data - voucher with code '%s' is not applicable to product with ID '%s'", voucherCode, productID)
		} else if err != nil {
			return nil, err
		}

		if isApplicable {
			discount = svc.calcSubscriptionPlanDiscount(product.Price, voucher)
		}
	}

	tax := product.Price.Mul(svc.taxRate)
	netPrice := product.Price.Sub(discount)
	if netPrice.IsNegative() {
		netPrice = decimal.NewFromInt(0)
	}

	planID, err := svc.idGen.New()
	if err != nil {
		svc.logger.WithField(svcMethod, "CreateSubscriptionPlan").WithError(err).Error("ID Generator error")
		return nil, fmt.Errorf("encountered error - failed to create valid ID for plan")
	}

	plan := SubscriptionPlan{
		ID:          planID,
		UserID:      userID,
		ProductID:   productID,
		Status:      SubscriptionPlanActive,
		StartDate:   time.Now().UTC(),
		EndDate:     time.Now().AddDate(0, 0, product.DurationDays).UTC(),
		NetPrice:    netPrice,
		GrossPrice:  product.Price,
		Tax:         tax,
		Discount:    discount,
		VoucherCode: &voucherCode,
	}

	if err := svc.repo.CreateSubscriptionPlan(ctx, plan); err != nil {
		svc.logger.WithField(svcMethod, "CreateSubscriptionPlan").WithError(err).Error("repo error")
		return nil, fmt.Errorf("encountered error - failed to create subscription plan")
	}

	return &plan, nil
}

func (svc *Service) validateVoucher(voucher *Voucher) error {
	if voucher.ValidFrom.After(time.Now()) || voucher.ValidUntil.Before(time.Now()) {
		return fmt.Errorf("voucher is out of valid time range")
	}

	if voucher.MaxUses != nil && *voucher.MaxUses < 1 {
		return fmt.Errorf("voucher has already exceeded max uses")
	}

	return nil
}

func (svc *Service) checkVoucherAppliesToProduct(ctx context.Context, voucherCode, productID string) (bool, error) {
	voucherProductEntry, err := svc.repo.GetVoucherProduct(ctx, voucherCode, productID)
	if err != nil {
		svc.logger.WithField(svcMethod, "checkVoucherAppliesToProduct").WithError(err).Error("repo error")
		return false, fmt.Errorf("encountered error - failed to check if voucher with code '%s' is applicable to product with ID '%s'", voucherCode, productID)
	}

	return voucherProductEntry != nil, nil
}

func (svc *Service) calcSubscriptionPlanDiscount(price decimal.Decimal, voucher *Voucher) decimal.Decimal {
	if voucher.Type == voucherTypeFixedAmount {
		return voucher.Value.Round(2)
	}

	return price.Mul(voucher.Value.Div(decimal.NewFromInt(100))).Round(2)
}

func (svc *Service) UpdateSubscriptionPlanStatus(ctx context.Context, planID, newStatus string) error {
	if planID == "" || newStatus == "" {
		return fmt.Errorf("bad input data - valid subscription plan ID and new status must be provided")
	}

	newPlanStatus, err := NewSubscriptionStatus(newStatus)
	if err != nil {
		return fmt.Errorf("bad input data - given status '%s' is invalid, must be one of ['ACTIVE', 'EXPIRED', 'PAUSED', 'CANCELLED']", newStatus)
	}

	plan, err := svc.repo.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		svc.logger.WithField(svcMethod, "UpdateSubscriptionPlanStatus").WithError(err).Error("repo error")
		return fmt.Errorf("encountered error - failed to retrieve subscription plan with ID '%s'", planID)
	}

	if newPlanStatus == SubscriptionPlanActive {
		return svc.resumeSubscriptionPlan(ctx, plan)
	} else if newPlanStatus == SubscriptionPlanPaused {
		return svc.pauseSubscriptionPlan(ctx, plan)
	} else if newPlanStatus == SubscriptionPlanCancelled {
		return svc.cancelSubscriptionPlan(ctx, plan)
	}

	return nil
}

func (svc *Service) resumeSubscriptionPlan(ctx context.Context, plan *SubscriptionPlan) error {
	switch plan.Status {
	case SubscriptionPlanPaused:
		if err := svc.repo.ResumeSubscriptionPlan(ctx, plan.ID); err != nil {
			svc.logger.WithField(svcMethod, "resumeSubscriptionPlan").WithError(err).Error("repo error")
			return fmt.Errorf("encountered error - failed to resume subscription plan with ID '%s'", plan.ID)
		}

		return nil
	case SubscriptionPlanActive:
		return fmt.Errorf("resource state conflict - subscription plan is already active")
	default:
		return fmt.Errorf("resource state conflict - cannot activate subscription plan from status '%s'", plan.Status)
	}
}

func (svc *Service) pauseSubscriptionPlan(ctx context.Context, plan *SubscriptionPlan) error {
	switch plan.Status {
	case SubscriptionPlanActive:
		newPauseID, err := svc.idGen.New()
		if err != nil {
			svc.logger.WithField(svcMethod, "pauseSubscriptionPlan").WithError(err).Error("ID Generator error")
			return fmt.Errorf("encountered error - failed to generate valid ID for new pause record")
		}

		if err := svc.repo.PauseSubscriptionPlan(ctx, plan.ID, newPauseID); err != nil {
			svc.logger.WithField(svcMethod, "pauseSubscriptionPlan").WithError(err).Error("repo error")
			return fmt.Errorf("encountered error - failed to resume subscription plan with ID '%s'", plan.ID)
		}

		return nil
	case SubscriptionPlanPaused:
		return fmt.Errorf("resource state conflict - subscription plan is already paused")
	default:
		return fmt.Errorf("resource state conflict - cannot pause subscription plan from status '%s'", plan.Status)
	}
}

func (svc *Service) cancelSubscriptionPlan(ctx context.Context, plan *SubscriptionPlan) error {
	switch plan.Status {
	case SubscriptionPlanCancelled:
		return fmt.Errorf("resource state conflict - subscription plan is already cancelled")
	case SubscriptionPlanExpired:
		return fmt.Errorf("resource state conflict - subscription plan is already expired")
	default:
		if err := svc.repo.CancelSubscriptionPlan(ctx, plan.ID); err != nil {
			svc.logger.WithField(svcMethod, "cancelSubscriptionPlan").WithError(err).Error("repo error")
			return fmt.Errorf("encountered error - failed to cancel subscription plan with ID '%s'", plan.ID)
		}
	}

	return nil
}

func (svc *Service) CreateVoucher(ctx context.Context, voucher Voucher, products []Product) error {
	return nil
}
