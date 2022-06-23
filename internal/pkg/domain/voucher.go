package domain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type VoucherType int

const (
	voucherTypeFixedAmount VoucherType = iota
	voucherTypePercentage
)

var voucherTypes = [2]string{"FIXED_AMOUNT", "PERCENTAGE"}

func NewVoucherType(voucherType string) (VoucherType, error) {
	switch voucherType {
	case "FIXED_AMOUNT":
		return voucherTypeFixedAmount, nil
	case "PERCENTAGE":
		return voucherTypePercentage, nil
	default:
		return -1, fmt.Errorf("invalid voucher type")
	}
}

func (vt VoucherType) String() string {
	return voucherTypes[vt]
}

func (vt *VoucherType) Scan(value interface{}) error {
	var err error
	*vt, err = NewVoucherType(string(value.([]byte)))
	return err
}

func (vt VoucherType) Value() (driver.Value, error) {
	return vt.String(), nil
}

type Voucher struct {
	Code       string          `json:"code"`
	Type       VoucherType     `json:"type"`
	Value      decimal.Decimal `json:"value"`
	MaxUses    *int            `json:"max_uses"`
	ValidFrom  time.Time       `json:"valid_from"`
	ValidUntil time.Time       `json:"valid_until"`
}

type VoucherProduct struct {
	VoucherCode string `json:"voucher_code"`
	ProductID   string `json:"product_id"`
}
