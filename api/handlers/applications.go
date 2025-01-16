package handlers

import (
	"chat-system/internal/database"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ApplicationHandlers struct {
	DBHandler *database.ApplicationsDatabaseHandler
}

func CreateApplicationHandlers() *ApplicationHandlers {
	dbHandler := database.NewApplicationsDatabaseHandler()
	return &ApplicationHandlers{DBHandler: dbHandler}
}

type createApplicationRequest struct {
	Name string `json:"name" validate:"required"`
}

type createApplicationResponse struct {
	Token string `json:"token"`
}

func (h *ApplicationHandlers) HandleCreateApplication(c echo.Context) error {
	request := new(createApplicationRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	token := uuid.New().String()

	err := h.DBHandler.InsertApplication(request.Name, token)
	if err != nil {
		log.Printf("error inserting application: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create application")
	}

	response := &Response[createApplicationResponse]{Data: createApplicationResponse{Token: token}}

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
