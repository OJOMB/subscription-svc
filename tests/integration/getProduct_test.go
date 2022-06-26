package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProduct_successCase(t *testing.T) {
	testCases := []struct {
		name      string
		productID string
		expected  string
	}{
		{
			name:      "fetch product basic 30 days",
			productID: productIDBasic30Days,
			expected: `{
				"id": "Ft2GgLRgN3FbMveklTy-W",
				"name": "basic",
				"description": "Access basic content for 30 days",
				"duration_days": 30,
				"price": 20.00
			}`,
		},
		{
			name:      "fetch product basic 1 year",
			productID: productIDBasic1Year,
			expected: `{
				"id": "mApy9b9Fqpt_WjghgUkSY",
				"name": "basic",
				"description": "Access basic content for 1 year",
				"duration_days": 365,
				"price": 200.00
			}`,
		},
		{
			name:      "fetch product premium 30 days",
			productID: productIDPremium30Days,
			expected: `{
				"id": "7Yn_IvvYsfkeo7-ysixd7",
				"name": "premium",
				"description": "Access all content for 30 days",
				"duration_days": 30,
				"price": 30.00
			}`,
		},
		{
			name:      "fetch product premium 1 year",
			productID: productIDPremium1Year,
			expected: `{
				"id": "UMp3k41eV5mY_iOkiElGm",
				"name": "premium",
				"description": "Access all content for 1 year",
				"duration_days": 365,
				"price": 300.00
			}`,
		},
	}

	for idx, testcase := range testCases {
		// here we pin tc so test cases can run in parallel
		// normally i wouldn't parallelise tests within a package because of the unnecessary complexity it potentially invites
		// packages in Go are tested in parallel anyway so normally parallelising the tests within a package would be overkill
		// but just here for the sake of demo lets have some fun!
		tc := testcase
		t.Run(
			fmt.Sprintf("test case %d: %s", idx+1, tc.name),
			func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrlProducts, tc.productID), nil)
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

func TestGetProductWithNonExistentID_failureCase(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrlProducts, "i_dont_exist"), nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.JSONEq(t, `{"error": "product not found"}`, string(respBody))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
