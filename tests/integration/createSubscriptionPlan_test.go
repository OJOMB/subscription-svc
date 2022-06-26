package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription_successCase(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		baseUrlSubscriptionPlans,
		strings.NewReader(
			fmt.Sprintf(
				`{
					"user_id": "%s",
					"product_id": "%s"
				}`,
				userIDWithoutExistingPlan1, productIDBasic30Days,
			),
		),
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var plan domain.SubscriptionPlan
	err = json.Unmarshal(respBody, &plan)
	assert.NoError(t, err)

	assert.Equal(t, plan.UserID, userIDWithoutExistingPlan1)
	assert.Equal(t, productIDBasic30Days, plan.ProductID)
	assert.Equal(t, decimal.NewFromInt(0), plan.Discount)
	assert.Equal(t, decimal.NewFromInt(20), plan.NetPrice)
	assert.Equal(t, decimal.NewFromInt(20), plan.GrossPrice)
	assert.Equal(t, decimal.NewFromInt(1), plan.Tax)

	// nanosecond differences mean I have to do this hack
	// if I had time i would set up and inject a dummy clock that would output a fixed time into the service to avoid this
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Year(), plan.EndDate.Year())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Month(), plan.EndDate.Month())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Day(), plan.EndDate.Day())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Hour(), plan.EndDate.Hour())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Minute(), plan.EndDate.Minute())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Second(), plan.EndDate.Second())
}

func TestCreateSubscriptionWithVoucher_successCase(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		baseUrlSubscriptionPlans,
		strings.NewReader(
			fmt.Sprintf(
				`{
					"user_id": "%s",
					"product_id": "%s",
					"voucher_code": "%s"
				}`,
				userIDWithoutExistingPlan2, productIDBasic30Days, voucherCode10Off,
			),
		),
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var plan domain.SubscriptionPlan
	err = json.Unmarshal(respBody, &plan)
	assert.NoError(t, err)

	assert.Equal(t, plan.UserID, userIDWithoutExistingPlan2)
	assert.Equal(t, productIDBasic30Days, plan.ProductID)
	assert.Equal(t, decimal.NewFromInt(10), plan.Discount)
	assert.Equal(t, decimal.NewFromInt(10), plan.NetPrice)
	assert.Equal(t, decimal.NewFromInt(20), plan.GrossPrice)
	assert.Equal(t, decimal.NewFromInt(1), plan.Tax)

	// nanosecond differences mean I have to do this hack
	// if I had time i would set up and inject a dummy clock that would output a fixed time into the service to avoid this
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Year(), plan.EndDate.Year())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Month(), plan.EndDate.Month())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Day(), plan.EndDate.Day())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Hour(), plan.EndDate.Hour())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Minute(), plan.EndDate.Minute())
	assert.Equal(t, plan.StartDate.AddDate(0, 0, 30).Second(), plan.EndDate.Second())
}
