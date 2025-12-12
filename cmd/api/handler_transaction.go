package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/cmd/api/dto"
	"github.com/ucok-man/tcsa/internal/data"
)

func (app *application) createTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionCreateDTO

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transaction := data.Transaction{
		UserId:    dto.UserId,
		Amount:    dto.Amount,
		Status:    data.TransactionStatusPending,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := app.models.Transactions.Insert(&transaction)
	if err != nil {
		return app.ErrInternalServer(err, "failed insert transaction", ctx.Request())
	}

	return ctx.JSON(http.StatusCreated, envelope{"transaction": transaction})
}

func (app *application) getByIdTransactionHandler(ctx echo.Context) error {
	transactionId, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrBadRequest(err.Error())
	}

	transaction, err := app.models.Transactions.GetById(transactionId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.ErrNotFound()
		default:
			return app.ErrInternalServer(err, "failed to get transaction by id", ctx.Request())
		}
	}

	return ctx.JSON(http.StatusOK, envelope{"transaction": transaction})
}

func (app *application) removeByIdTransactionHandler(ctx echo.Context) error {
	transactionId, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrBadRequest(err.Error())
	}

	transaction, err := app.models.Transactions.GetById(transactionId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.ErrNotFound()
		default:
			return app.ErrInternalServer(err, "failed to get transaction by id", ctx.Request())
		}
	}

	err = app.models.Transactions.DeleteOne(transaction.ID)
	if err != nil {
		return app.ErrInternalServer(err, "failed to delete transaction", ctx.Request())
	}

	return ctx.JSON(http.StatusOK, envelope{
		"message": fmt.Sprintf("transaction record with id %v succesfully dedelte", transaction.ID),
	})
}

func (app *application) updateByIdTransactionHandler(ctx echo.Context) error {
	transactionId, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrBadRequest(err.Error())
	}

	var dto dto.TransactionUpdateDTO

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transaction, err := app.models.Transactions.GetById(transactionId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.ErrNotFound()
		default:
			return app.ErrInternalServer(err, "failed to get transaction by id", ctx.Request())
		}
	}

	if dto.Amount != nil {
		transaction.Amount = *dto.Amount
	}
	if dto.Status != nil {
		transaction.Status = data.TransactionStatus(*dto.Status)
	}
	transaction.UpdatedAt = time.Now()

	err = app.models.Transactions.Update(transaction)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return app.ErrEditConflict()
		default:
			return app.ErrInternalServer(err, "failed to update transaction", ctx.Request())
		}
	}

	return ctx.JSON(http.StatusOK, envelope{"transaction": transaction})
}
