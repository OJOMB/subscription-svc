package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type updateSubscriptionPlanStatusReq struct {
	NewStatus string `json:"new_status"`
}

func (app *App) handleUpdateSubscriptionPlanStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		planID := vars[urlVarSubscriptionPlanID]

		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("failed to read request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var updateData updateSubscriptionPlanStatusReq
		if err := json.Unmarshal(reqBodyBytes, &updateData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		if err := app.service.UpdateSubscriptionPlanStatus(r.Context(), planID, updateData.NewStatus); err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
