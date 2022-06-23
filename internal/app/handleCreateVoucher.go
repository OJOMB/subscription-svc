package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
)

type createVoucherReq struct {
	Voucher  domain.Voucher   `json:"voucher"`
	Products []domain.Product `json:"products"`
}

func (app *App) handleCreateVoucher() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("failed to read request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var createVoucherData createVoucherReq
		if err := json.Unmarshal(reqBodyBytes, &createVoucherData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		if err := app.service.CreateVoucher(r.Context(), createVoucherData.Voucher, createVoucherData.Products); err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
