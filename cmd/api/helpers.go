package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type envelope map[string]any

func (app *application) GetParamId(ctx echo.Context) (int, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return int(0), fmt.Errorf("invalid id parameter")
	}

	return int(id), nil
}

func (app *application) SortColumn(value string) string {
	column := strings.TrimPrefix(value, "-")
	return column
}

func (app *application) SortDirection(value string) string {
	var direction string

	if strings.HasPrefix(value, "-") {
		direction = "DESC"
	} else {
		direction = "ASC"
	}
	return direction

}

func (app *application) PageOffset(page, pageSize int) int {
	offset := (page - 1) * pageSize
	return offset
}
