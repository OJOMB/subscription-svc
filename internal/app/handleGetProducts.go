package app

import (
	"encoding/json"
	"net/http"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
)

const (
	handleGetProducts   = "handleGetProducts"
	urlQueryVoucherCode = "voucher_code"
)

func (app *App) handleGetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		voucherCode := r.URL.Query().Get(urlQueryVoucherCode)

		var products []*domain.Product
		var err error
		if voucherCode != "" {
			products, err = app.service.GetProductsWithVoucher(r.Context(), voucherCode)
		} else {
			products, err = app.service.GetProducts(r.Context())
		}

		if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBytes, err := json.Marshal(products)
		if err != nil {
			app.logger.WithField(appHandler, handleGetProducts).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.Write(respBytes)
	}
}
