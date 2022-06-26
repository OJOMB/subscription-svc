package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetSubscriptionPlan = "handleGetSubscriptionPlan"

func (app *App) handleGetSubscriptionPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionPlanID := vars[urlVarSubscriptionPlanID]

		plan, err := app.service.GetSubscriptionPlan(r.Context(), subscriptionPlanID)
		if plan == nil {
			apperr := newAppErr("subscription plan not found", http.StatusNotFound)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		} else if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBytes, err := json.Marshal(plan)
		if err != nil {
			app.logger.WithField(appHandler, handleGetSubscriptionPlan).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.Write(respBytes)
	}
}
