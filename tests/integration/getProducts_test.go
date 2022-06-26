package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts_happyCase(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, baseUrlProducts, nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.JSONEq(
		t,
		`[
			{
				"id": "7Yn_IvvYsfkeo7-ysixd7",
				"name": "premium",
				"description": "Access all content for 30 days",
				"duration_days": 30,
				"price": 30
			},
			{
				"id": "Ft2GgLRgN3FbMveklTy-W",
				"name": "basic",
				"description": "Access basic content for 30 days",
				"duration_days": 30,
				"price": 20
			},
			{
				"id": "mApy9b9Fqpt_WjghgUkSY",
				"name": "basic",
				"description": "Access basic content for 1 year",
				"duration_days": 365,
				"price": 200
			},
			{
				"id": "UMp3k41eV5mY_iOkiElGm",
				"name": "premium",
				"description": "Access all content for 1 year",
				"duration_days": 365,
				"price": 300
			}
		]`,
		string(respBody),
	)
}

func TestGetProductsWithVoucherCode_successCase(t *testing.T) {
	// the dummy data is setup in such a way that the voucher codes are only applicable to the 30 day products
	// so we would expect only the 30 day products to be returned here
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?voucher_code=%s", baseUrlProducts, voucherCode20Off), nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.JSONEq(
		t,
		`[
			{
				"id": "7Yn_IvvYsfkeo7-ysixd7",
				"name": "premium",
				"description": "Access all content for 30 days",
				"duration_days": 30,
				"price": 30
			},
			{
				"id": "Ft2GgLRgN3FbMveklTy-W",
				"name": "basic",
				"description": "Access basic content for 30 days",
				"duration_days": 30,
				"price": 20
			}
		]`,
		string(respBody),
	)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetProductsWithNonExistentVoucherCode_failureCase(t *testing.T) {
	// the dummy data is setup in such a way that the voucher codes are only applicable to the 30 day products
	// so we would expect only the 30 day products to be returned here
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?voucher_code=%s", baseUrlProducts, voucherCodeNonExistent), nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.JSONEq(t, `{"error": "resource not found - no such voucher with code 'i_dont_exist' exists"}`, string(respBody))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
