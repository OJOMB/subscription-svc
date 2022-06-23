package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/subscription-svc/internal/pkg/domain"
)

const handleCreateProduct = "handleCreateProduct"

func (app *App) handleCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("failed to read request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var product domain.Product
		if err := json.Unmarshal(reqBodyBytes, &product); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		CreatedProduct, err := app.service.CreateProduct(r.Context(), product)
		if err != nil {
			apperr := app.newAppErrFromSvcErr(err)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(CreatedProduct)
		if err != nil {
			app.logger.WithField(appHandler, handleCreateProduct).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(respBodyBytes)
	}
}
