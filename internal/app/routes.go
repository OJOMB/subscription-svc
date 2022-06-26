package app

import (
	"fmt"
	"net/http"

	"github.com/OJOMB/subscription-svc/internal/app/middleware"
)

const (
	urlVarProductID          = "productID"
	urlVarSubscriptionPlanID = "planID"
	urlVarVoucherCode        = "voucherCode"
)

func (app *App) routes() {
	if app.docsEnabled {
		fs := http.FileServer(http.Dir("./api/OpenAPI/"))
		app.router.PathPrefix("/docs").Handler(http.StripPrefix("/docs", fs))
	}

	// users
	app.router.HandleFunc("/api/v1/users", app.handleCreateUser()).Methods(http.MethodPost)

	// products
	app.router.HandleFunc("/api/v1/products", app.handleGetProducts()).Methods(http.MethodGet)
	app.router.HandleFunc(fmt.Sprintf("/api/v1/products/{%s}", urlVarProductID), app.handleGetProduct()).Methods(http.MethodGet)
	app.router.HandleFunc("/api/v1/products", app.handleCreateProduct()).Methods(http.MethodPost)

	// subscription-plans
	app.router.HandleFunc("/api/v1/subscription-plans", app.handleCreateSubscriptionPlan()).Methods(http.MethodPost)
	app.router.HandleFunc(fmt.Sprintf("/api/v1/subscription-plans/{%s}", urlVarSubscriptionPlanID), app.handleGetSubscriptionPlan()).Methods(http.MethodGet)
	app.router.HandleFunc(
		fmt.Sprintf("/api/v1/subscription-plans/{%s}/status", urlVarSubscriptionPlanID),
		app.handleUpdateSubscriptionPlanStatus(),
	).Methods(http.MethodPut)

	// vouchers
	app.router.HandleFunc(fmt.Sprintf("/api/v1/vouchers/{%s}", urlVarVoucherCode), app.handleCreateVoucher()).Methods(http.MethodPut)

	app.router.Use(middleware.NewRequestResponseLogger(app.logger).Middleware)
}
