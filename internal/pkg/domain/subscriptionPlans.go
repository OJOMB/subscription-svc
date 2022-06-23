package domain

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type SubscriptionStatus int

const (
	SubscriptionPlanActive SubscriptionStatus = iota
	SubscriptionPlanExpired
	SubscriptionPlanPaused
	SubscriptionPlanCancelled
)

var subscriptionPlanStatuses = [4]string{
	"ACTIVE",
	"EXPIRED",
	"PAUSED",
	"CANCELLED",
}

func NewSubscriptionStatus(status string) (SubscriptionStatus, error) {
	switch status {
	case "ACTIVE":
		return SubscriptionPlanActive, nil
	case "EXPIRED":
		return SubscriptionPlanExpired, nil
	case "PAUSED":
		return SubscriptionPlanPaused, nil
	case "CANCELLED":
		return SubscriptionPlanCancelled, nil
	default:
		return -1, fmt.Errorf("invalid status")
	}
}

func (ss SubscriptionStatus) String() string {
	return subscriptionPlanStatuses[ss]
}

func (ss *SubscriptionStatus) Scan(value interface{}) error {
	var err error
	*ss, err = NewSubscriptionStatus(string(value.([]byte)))
	return err
}

func (u SubscriptionStatus) Value() (driver.Value, error) {
	return u.String(), nil
}

func (ss SubscriptionStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(subscriptionPlanStatuses[ss])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (ss *SubscriptionStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*ss, err = NewSubscriptionStatus(j)
	return err
}

type SubscriptionPlan struct {
	ID          string             `json:"id"`
	UserID      string             `json:"user_id"`
	ProductID   string             `json:"product_id"`
	Status      SubscriptionStatus `json:"status"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     time.Time          `json:"end_date"`
	NetPrice    decimal.Decimal    `json:"net_price"`
	GrossPrice  decimal.Decimal    `json:"gross_price"`
	Tax         decimal.Decimal    `json:"tax"`
	Discount    decimal.Decimal    `json:"discount"`
	VoucherCode *string            `json:"voucher_code"`
}

type SubscriptionPlanPause struct {
	ID                 string    `json:"id"`
	SubscriptionPlanID string    `json:"subscription_plan_id"`
	PauseDate          time.Time `json:"pause_date"`
	EndDateAtPause     time.Time `json:"end_date_at_pause"`
	Resumed            bool      `json:"resumed"`
}
