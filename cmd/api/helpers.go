package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type envelope map[string]any

func (app *application) getParamId(ctx echo.Context) (int, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return int(0), fmt.Errorf("invalid id parameter")
	}

	return int(id), nil
}

func (app *application) SortColumn(value string) string {
	return strings.TrimPrefix(value, "-")
}

func (app *application) SortDirection(value string) string {
	if strings.HasPrefix(value, "-") {
		return "DESC"
	}

	return "ASC"
}

func (app *application) PageOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
