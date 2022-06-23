package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetProduct = "handleGetProduct"

func (app *App) handleGetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID := vars[urlVarProductID]

		product, err := app.service.GetProduct(r.Context(), productID)
		if product == nil {
			apperr := newAppErr("product not found", http.StatusNotFound)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		} else if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBytes, err := json.Marshal(product)
		if err != nil {
			app.logger.WithField(appHandler, handleGetProduct).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.Write(respBytes)
	}
}
