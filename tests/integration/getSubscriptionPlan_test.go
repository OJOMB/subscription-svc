package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubscriptionPlan_successCase(t *testing.T) {
	testCases := []struct {
		name               string
		subscriptionPlanID string
		expected           string
	}{
		{
			name:               "fetch subscription plan with ID 'KlRuHEGEseQBogzXpc8ns'",
			subscriptionPlanID: "KlRuHEGEseQBogzXpc8ns",
			expected: `{
				"id" : "KlRuHEGEseQBogzXpc8ns",
				"user_id" : "Yh4WkMFE0MTZnUjD9aYFy",
				"product_id" : "Ft2GgLRgN3FbMveklTy-W",
				"status" : "ACTIVE",
				"start_date" : "2022-06-23T00:00:00Z",
				"end_date" : "2022-07-23T00:00:00Z",
				"net_price" : 18.90,
				"gross_price" : 20.00,
				"tax" : 1.00,
				"discount" : 2.00,
				"voucher_code" : "10-percent-off"
			}`,
		},
		{
			name:               "fetch subscription plan with ID 'LwBhyg4FcWwUH1ORGVXDU'",
			subscriptionPlanID: "LwBhyg4FcWwUH1ORGVXDU",
			expected: `{
				"id" : "LwBhyg4FcWwUH1ORGVXDU",
				"user_id" : "521Qlk96BbJoGKseV1nPZ",
				"product_id" : "mApy9b9Fqpt_WjghgUkSY",
				"status" : "ACTIVE",
				"start_date" : "2022-01-01T00:00:00Z",
				"end_date" : "2023-01-01T00:00:00Z",
				"net_price" : 315.00,
				"gross_price" : 300.00,
				"tax" : 15.00,
				"discount" : 0.00,
				"voucher_code" : null
			}`,
		},
		{
			name:               "fetch subscription plan with ID 'rhmeplLbg8bxWsqLzZQ6i'",
			subscriptionPlanID: "rhmeplLbg8bxWsqLzZQ6i",
			expected: `{
				"id": "rhmeplLbg8bxWsqLzZQ6i",
				"user_id": "b3lFU5zF9zB37DxKk-zCC",
				"product_id": "7Yn_IvvYsfkeo7-ysixd7",
				"status": "ACTIVE",
				"start_date": "2022-06-01T00:00:00Z",
				"end_date": "2022-07-01T00:00:00Z",
				"net_price": 31.5,
				"gross_price": 30,
				"tax": 1.5,
				"discount": 0,
				"voucher_code": null
			}`,
		},
		{
			name:               "fetch subscription plan with ID 'sFF_eBjQgBcTZAKPCSzo5'",
			subscriptionPlanID: "sFF_eBjQgBcTZAKPCSzo5",
			expected: `{
				"id" : "sFF_eBjQgBcTZAKPCSzo5",
				"user_id" : "FUQQzY_-4Tv_p7SFeHJEI",
				"product_id" : "Ft2GgLRgN3FbMveklTy-W",
				"status" : "CANCELLED",
				"start_date" : "2022-06-01T00:00:00Z",
				"end_date" : "2022-07-01T00:08:55Z",
				"net_price" : 21.00,
				"gross_price" : 20.00,
				"tax" : 1.00,
				"discount" : 0.00,
				"voucher_code" : null
			}`,
		},
	}

	for idx, tc := range testCases {
		t.Run(
			fmt.Sprintf("test case %d: %s", idx+1, tc.name),
			func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrlSubscriptionPlans, tc.subscriptionPlanID), nil)
				assert.NoError(t, err)

				resp, err := http.DefaultClient.Do(req)
				assert.NoError(t, err)

				respBody, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.JSONEq(t, tc.expected, string(respBody))
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			},
		)
	}
}
