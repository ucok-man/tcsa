package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/internal/serializer"
	"github.com/ucok-man/tcsa/internal/validator"
)

func (app *application) routes() http.Handler {
	ec := echo.New()
	ec.JSONSerializer = serializer.New()
	ec.Validator = validator.New()
	ec.Logger = app.logger
	ec.HTTPErrorHandler = app.HTTPErrorHandler

	ec.Use(app.withRecover())
	ec.Use(app.withRequestLogger())

	ec.GET("/healthcheck", app.healthcheckHandler)

	transactions := ec.Group("/transactions")
	{
		transactions.POST("", app.createTransactionHandler)
		transactions.GET("/:id", app.getByIdTransactionHandler)
		transactions.DELETE("/:id", app.removeByIdTransactionHandler)

	}

	return ec
}
