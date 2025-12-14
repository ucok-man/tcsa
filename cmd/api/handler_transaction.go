package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/cmd/api/dto"
	"github.com/ucok-man/tcsa/internal/data"
	"github.com/ucok-man/tcsa/internal/utility"
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

	return ctx.JSON(http.StatusCreated, envelope{
		"data": transaction,
	})
}

func (app *application) getByIdTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionParamIdDTO

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transaction, err := app.models.Transactions.GetById(dto.TransactionId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.ErrNotFound()
		default:
			return app.ErrInternalServer(err, "failed to get transaction by id", ctx.Request())
		}
	}

	return ctx.JSON(http.StatusOK, envelope{
		"data": transaction,
	})
}

func (app *application) removeByIdTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionParamIdDTO

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transaction, err := app.models.Transactions.GetById(dto.TransactionId)
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
		"data": transaction,
	})
}

func (app *application) updateByIdTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionUpdateDTO

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transaction, err := app.models.Transactions.GetById(dto.TransactionId)
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

	return ctx.JSON(http.StatusOK, envelope{
		"data": transaction,
	})
}

func (app *application) getAllTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionGetAllDTO

	// Set Default Value
	dto.Pagination.Page = utility.SetPtrValue(1)
	dto.Pagination.PageSize = utility.SetPtrValue(10)
	dto.Sort.Value = utility.SetPtrValue("id")

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	transactions, metadata, err := app.models.Transactions.GetAll(data.TransactionGetAllParam{
		Page:          *dto.Pagination.Page,
		PageSize:      *dto.Pagination.PageSize,
		PageOffset:    app.PageOffset(*dto.Pagination.Page, *dto.Pagination.PageSize),
		SortColumn:    app.SortColumn(*dto.Sort.Value),
		SortDirection: app.SortDirection(*dto.Sort.Value),
		FilterStatus:  utility.DerefOrDefault(dto.Filter.Status, ""),
		FilterUserId:  utility.DerefOrDefault(dto.Filter.UserId, 0),
	})
	if err != nil {
		return app.ErrInternalServer(err, "failed get all transactions", ctx.Request())
	}

	return ctx.JSON(http.StatusOK, envelope{
		"data":     transactions,
		"metadata": metadata,
	})
}
