package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetSegments returns the saved segments (a JSON array of {name, query}).
func (a *App) GetSegments(c echo.Context) error {
	out, err := a.core.GetSegments()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{out})
}

// UpdateSegments replaces the full saved-segments list.
func (a *App) UpdateSegments(c echo.Context) error {
	var b json.RawMessage
	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := a.core.SaveSegments(b); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{true})
}
