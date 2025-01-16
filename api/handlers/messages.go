package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MessageHandlers struct {
}

func (h *MessageHandlers) HandleCreateMessage(c echo.Context) error {
	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandlers) HandleGetMessages(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandlers) HandleGetMessage(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandlers) HandleSearchMessages(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}
func (h *MessageHandlers) HandleUpdateMessage(c echo.Context) error {

	response := &Response[string]{Data: "place holder"}
	return c.JSON(http.StatusOK, response)
}
