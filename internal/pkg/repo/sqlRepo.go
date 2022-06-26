package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const componentSQLRepo = "SQLRepo"

type SQlRepo struct {
	db     *sql.DB
	logger *logrus.Entry
}

func NewSQLRepo(db *sql.DB, logger *logrus.Logger) *SQlRepo {
	return &SQlRepo{
		db:     db,
		logger: logger.WithField("component", componentSQLRepo),
	}
}

func (r *SQlRepo) GetProduct(ctx context.Context, productID string) (*domain.Product, error) {
	var product domain.Product

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, duration_days, price FROM products WHERE id = ?;`, productID,
	).Scan(&product.ID, &product.Name, &product.Description, &product.DurationDays, &product.Price)
	if err != nil && err == sql.ErrNoRows {
		r.logger.WithField("method", "GetProduct").Infof("found no product with ID %s", productID)
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetProduct").Error("failed to retrieve product by ID")
		return nil, err
	}

	return &product, nil
}

func (r *SQlRepo) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, description, duration_days, price FROM products;`,
	)
	if err != nil && err == sql.ErrNoRows {
		r.logger.WithField("method", "GetProducts").Info("found no products")
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetProducts").Error("failed to retrieve products")
		return nil, err
	}

	products := make([]*domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.DurationDays,
			&product.Price,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetProducts").Error("failed to scan retrieved products")
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (r *SQlRepo) GetProductsWithVoucher(ctx context.Context, voucherCode string) ([]*domain.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT p.id, p.name, p.description, p.duration_days, p.price FROM products_vouchers AS pv
		JOIN products AS p
		WHERE pv.voucher_code = ? AND p.id = pv.product_id;`,
		&voucherCode,
	)
	if err != nil && err == sql.ErrNoRows {
		r.logger.WithField("method", "GetProductsWithVoucher").Infof("found no products with voucher code: %s", voucherCode)
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetProductsWithVoucher").Error("failed to retrieve products with voucher")
		return nil, err
	}

	products := make([]*domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.DurationDays,
			&product.Price,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetProductsWithVoucher").Error("failed to scan retrieved products")
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (r *SQlRepo) CreateSubscriptionPlan(ctx context.Context, plan domain.SubscriptionPlan) error {
	startDate := plan.StartDate.UTC()
	endDate := plan.EndDate.UTC()

	if _, err := r.db.ExecContext(
		ctx,
		`INSERT INTO subscription_plans
		(id, user_id, product_id, start_date, end_date, net_price, gross_price, tax, discount, voucher_code)
		VALUES (?,?,?,?,?,?,?,?,?,?);`,
		plan.ID, plan.UserID, plan.ProductID, &startDate, &endDate, plan.NetPrice, plan.GrossPrice, plan.Tax, plan.Discount, plan.VoucherCode,
	); err != nil {
		r.logger.WithError(err).WithField("method", "CreateSubscriptionPlan").Error("failed to insert subscription plan")
		return err
	}

	return nil
}

func (r *SQlRepo) GetSubscriptionPlan(ctx context.Context, planID string) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, product_id, status, start_date, end_date, net_price, gross_price, tax, discount, voucher_code
		FROM subscription_plans WHERE id = ?`,
		planID,
	).Scan(
		&plan.ID, &plan.UserID, &plan.ProductID, &plan.Status,
		&plan.StartDate, &plan.EndDate,
		&plan.NetPrice, &plan.GrossPrice, &plan.Tax, &plan.Discount, &plan.VoucherCode,
	); err != nil && err == sql.ErrNoRows {
		r.logger.WithField("method", "GetSubscriptionPlan").Infof("found no such subscription plan with ID %s", planID)
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetSubscriptionPlan").Errorf("failed to retrieve subscription plan with ID %s", planID)
		return nil, err
	}

	return &plan, nil
}

func (r *SQlRepo) PauseSubscriptionPlan(ctx context.Context, planID, pauseID string, pauseTime time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start tx: %v", err)
	}

	// for data consistency we need to check the status is pausable (again)
	// skipping doing so here could lead to double pauses
	var originalStatusStr string
	var originalEndDate time.Time
	if err := tx.QueryRowContext(ctx, `SELECT status, end_date FROM subscription_plans WHERE id = ?`, &planID).
		Scan(&originalStatusStr, &originalEndDate); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to retrieve data on plan with ID %s: %v", planID, err)
	}

	oldStatus, err := domain.NewSubscriptionStatus(originalStatusStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to retrieve data on plan with ID %s: %v", planID, err)
	}

	if oldStatus != domain.SubscriptionPlanActive {
		tx.Rollback()
		return fmt.Errorf("cannot pause plan with status %s", oldStatus.String())
	}

	originalEndDateStr := originalEndDate.UTC()
	if _, err = tx.ExecContext(
		ctx,
		`INSERT INTO subscription_plan_pauses (id, subscription_plan_id, pause_date, end_date_at_pause)
		 VALUES (?, ?, ?, ?);`,
		&pauseID, &planID, &pauseTime, &originalEndDateStr,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create subscription plan pause record for plan with ID %s: %v", planID, err)
	}

	statusPaused := domain.SubscriptionPlanPaused.String()
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE subscription_plans SET status = ? WHERE id = ?;`,
		&statusPaused, &planID,
	); err != nil {
		tx.Rollback()
		r.logger.WithError(err).WithField("method", "UpdateSubscriptionPlanStatus").Error("failed to update subscription plan status")
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *SQlRepo) ResumeSubscriptionPlan(ctx context.Context, planID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start tx: %v", err)
	}

	var currentStatusStr string
	if err := tx.QueryRowContext(
		ctx,
		`SELECT status FROM subscription_plans WHERE id = ?;`,
		&planID,
	).Scan(&currentStatusStr); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to retrieve data on subscription plan with ID %s: %v", planID, err)
	}

	currentStatus, err := domain.NewSubscriptionStatus(currentStatusStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to convert current Subscription Status '%s': %v", currentStatusStr, err)
	}

	if currentStatus != domain.SubscriptionPlanPaused {
		tx.Rollback()
		return fmt.Errorf("cannot resume subscription plan with status %s", currentStatus.String())
	}

	var pause domain.SubscriptionPlanPause
	if err := tx.QueryRowContext(
		ctx,
		`SELECT id, pause_date, end_date_at_pause
		FROM subscription_plan_pauses
		WHERE subscription_plan_id = ? AND resumed IS NULL;`,
		&planID,
	).Scan(&pause.ID, &pause.PauseDate, &pause.EndDateAtPause); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to retrieve data on current pause for subscription plan with ID %s: %v", planID, err)
	}

	now := time.Now()
	newEndDate := now.Add(pause.EndDateAtPause.Sub(pause.PauseDate))
	statusActive := domain.SubscriptionPlanActive.String()
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE subscription_plans SET status = ?, end_date = ? WHERE id = ?;`,
		&statusActive, &newEndDate, &planID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "ResumeSubscriptionPlan").Error("failed to update subscription plan with ID %s to status ACTIVE", planID)
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE subscription_plan_pauses SET resumed = ? WHERE id = ?;`,
		now, &pause.ID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "ResumeSubscriptionPlan").Error("failed to add resumed timestamp to pause record with ID %s", pause.ID)
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *SQlRepo) CancelSubscriptionPlan(ctx context.Context, planID string) error {
	statusCancelledStr := domain.SubscriptionPlanCancelled.String()

	if _, err := r.db.ExecContext(
		ctx,
		`UPDATE subscription_plans SET status = ? WHERE id = ?;`,
		&statusCancelledStr, &planID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "UpdateSubscriptionPlanStatus").Error("failed to update subscription plan status")
		return err
	}

	return nil
}

func (r *SQlRepo) GetVoucher(ctx context.Context, voucherCode string) (*domain.Voucher, error) {
	var voucherValue float64

	var voucher domain.Voucher
	err := r.db.QueryRowContext(
		ctx,
		`SELECT code, type, value, max_uses, valid_from, valid_until FROM vouchers WHERE code = ?;`,
		voucherCode,
	).Scan(&voucher.Code, &voucher.Type, &voucherValue, &voucher.MaxUses, &voucher.ValidFrom, &voucher.ValidUntil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "GetVoucher").Error("failed to retrieve voucher by ID")
		return nil, err
	}

	voucher.Value = decimal.NewFromFloat(voucherValue)

	return &voucher, nil
}

func (r *SQlRepo) GetVoucherProduct(ctx context.Context, voucherCode, productID string) (*domain.VoucherProduct, error) {
	var voucherProductRecord domain.VoucherProduct
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT voucher_code, product_id FROM products_vouchers WHERE voucher_code = ? AND product_id = ?;`,
		voucherCode, productID,
	).Scan(&voucherProductRecord.VoucherCode, &voucherProductRecord.ProductID); err != nil && err == sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "GetVoucher").Error("found no such voucher product record")
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetVoucher").Error("failed to retrieve voucher by ID")
		return nil, err
	}

	return &voucherProductRecord, nil
}

func (r *SQlRepo) GetCurrentSubscriptionPlan(ctx context.Context, userID string) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, product_id, status, start_date, end_date, net_price, gross_price, tax, discount, voucher_code
		FROM subscription_plans WHERE user_id = ? AND (status = 'ACTIVE' OR status = 'PAUSED');`,
		userID,
	).Scan(
		&plan.ID, &plan.UserID, &plan.ProductID, &plan.Status,
		&plan.StartDate, &plan.EndDate,
		&plan.NetPrice, &plan.GrossPrice, &plan.Tax, &plan.Discount, &plan.VoucherCode,
	); err != nil && err == sql.ErrNoRows {
		r.logger.WithField("method", "GetCurrentUserSubscriptionPlan").Infof("found no current subscription plan for user with ID %s", userID)
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetCurrentUserSubscriptionPlan").Errorf("failed to retrieve current subscription plan for user with ID %s", userID)
		return nil, err
	}

	return &plan, nil
}
