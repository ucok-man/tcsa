package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/internal/tlog"
)

func (app *application) HTTPErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details,omitempty"`
	}

	if he, ok := err.(*echo.HTTPError); ok {
		response.Code = http.StatusText(he.Code)
		switch he.Code {
		case http.StatusUnprocessableEntity:
			response.Message = "unable to proccess request because some malformed input"
			response.Details = he.Message

		// Change default notfound and method not allowed.
		case http.StatusNotFound:
			response.Message = "the requested resource could not be found"
		case http.StatusMethodNotAllowed:
			response.Message = fmt.Sprintf("the %s method is not supported for this resource", ctx.Request().Method)
		default:
			msg, ok := he.Message.(string)
			if !ok {
				msg = fmt.Sprintf("%v", he.Message)
			}

			// Extract just the actual error message if it's in Echo's format
			if idx := strings.Index(msg, "message="); idx != -1 {
				msg = msg[idx+8:] // Skip "message="
				// Remove internal error suffix if present
				if commaIdx := strings.Index(msg, ", internal="); commaIdx != -1 {
					msg = msg[:commaIdx]
				}
			}
			response.Message = msg
		}

		err = ctx.JSON(he.Code, envelope{"error": response})
		if err != nil {
			app.logger.Errorj(tlog.JSON{
				"message": "error sending json response",
				"error":   err,
			})
			ctx.Response().WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Uncaught Error
	app.logger.Errorj(tlog.JSON{
		"message": "unhandled error occured",
		"error":   err,
	})
	response.Code = http.StatusText(http.StatusInternalServerError)
	response.Message = "the server encountered a problem and could not process your request"
	err = ctx.JSON(http.StatusInternalServerError, envelope{"error": response})
	if err != nil {
		app.logger.Errorj(tlog.JSON{
			"message": "error sending json response",
			"error":   err,
		})
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) ErrInternalServer(err error, message string, req *http.Request) error {
	app.logger.Errorj(tlog.JSON{
		"message": message,
		"path":    req.URL,
		"method":  req.Method,
		"error":   err,
	})
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		"the server encountered a problem and could not process your request",
	)
}

func (app *application) ErrNotFound(customMsg ...string) error {
	message := "the requested resource could not be found"
	if len(customMsg) > 0 && customMsg[0] != "" {
		message = customMsg[0]
	}
	return echo.NewHTTPError(http.StatusNotFound, message)
}

func (app *application) ErrMethodNotAllowed(method string) error {
	return echo.NewHTTPError(
		http.StatusMethodNotAllowed,
		fmt.Sprintf("the %s method is not supported for this resource", method),
	)
}

func (app *application) ErrBadRequest(message string) error {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func (app *application) ErrFailedValidation(errmap any) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, errmap)
}

func (app *application) ErrEditConflict() error {
	return echo.NewHTTPError(
		http.StatusConflict,
		"unable to update the record due to an edit conflict, please try again",
	)
}

func (app *application) ErrRateLimitExceeded() error {
	return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
}

func (app *application) ErrForbidden(message ...string) error {
	msg := "forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return echo.NewHTTPError(http.StatusForbidden, msg)
}
