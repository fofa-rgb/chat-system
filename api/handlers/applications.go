package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApplicationHandlers struct {
}

func (h *ApplicationHandlers) HandleCreateApplication(c echo.Context) error {
	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}

func (h *ApplicationHandlers) HandleGetApplication(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}

func (h *ApplicationHandlers) HandleUpdateApplication(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}
