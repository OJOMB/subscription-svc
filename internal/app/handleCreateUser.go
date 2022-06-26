package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type createUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`
}

func (app *App) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("failed to read request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var userData createUserReq
		if err := json.Unmarshal(reqBodyBytes, &userData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		user, err := app.service.CreateUser(r.Context(), userData.FirstName, userData.LastName, userData.DOB, userData.Email)
		if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(user)
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
