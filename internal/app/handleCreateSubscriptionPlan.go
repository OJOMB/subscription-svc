package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const handleCreateSubscriptionPlan = "handleCreateSubscriptionPlan"

type createSubscriptionPlanreq struct {
	UserID      string `json:"user_id"`
	ProductID   string `json:"product_id"`
	VoucherCode string `json:"voucher_code"`
}

func (app *App) handleCreateSubscriptionPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("failed to read request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var createData createSubscriptionPlanreq
		if err := json.Unmarshal(reqBodyBytes, &createData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		CreatedSubscriptionPlan, err := app.service.CreateSubscriptionPlan(
			r.Context(), createData.UserID, createData.ProductID, createData.VoucherCode,
		)
		if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(CreatedSubscriptionPlan)
		if err != nil {
			app.logger.WithField(appHandler, handleCreateSubscriptionPlan).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(respBodyBytes)
	}
}
