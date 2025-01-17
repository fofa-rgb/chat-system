package handlers

import (
	"chat-system/internal/database"
	"chat-system/internal/models"
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

func (h *ApplicationHandlers) HandleCreateApplication(c echo.Context) error {
	request := new(createApplicationRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}

	token := uuid.New().String()

	err := h.DBHandler.InsertApplication(request.Name, token)
	if err != nil {
		log.Printf("error inserting application: %v", err)
		return echo.ErrInternalServerError
	}

	response := &response[createApplicationResponse]{Data: createApplicationResponse{Token: token}}

	return c.JSON(http.StatusOK, response)
}

func (h *ApplicationHandlers) HandleGetApplicationByToken(c echo.Context) error {
	token := c.Param("token")
	app, err := h.DBHandler.GetApplicationByToken(token)
	if err != nil {
		log.Printf("error getting application: %v", err)
		return echo.ErrInternalServerError
	}
	userApp := models.UserExposedApplication{
		Name:       app.Name,
		Token:      app.Token,
		ChatsCount: app.ChatsCount,
	}
	response := &response[models.UserExposedApplication]{Data: userApp}
	return c.JSON(http.StatusOK, response)
}
func (h *ApplicationHandlers) HandleGetAllApplications(c echo.Context) error {
	allApps, err := h.DBHandler.GetAllApplications()
	if err != nil {
		log.Printf("error getting applications: %v", err)
		return echo.ErrInternalServerError
	}
	var userExposedApps []models.UserExposedApplication
	for _, app := range allApps {
		userExposedApps = append(userExposedApps, models.UserExposedApplication{
			Name:       app.Name,
			Token:      app.Token,
			ChatsCount: app.ChatsCount,
		})
	}
	response := &response[[]models.UserExposedApplication]{Data: userExposedApps}

	return c.JSON(http.StatusOK, response)
}

func (h *ApplicationHandlers) HandleUpdateApplicationName(c echo.Context) error {
	token := c.Param("token")
	request := new(updateApplicationNameRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}
	err := h.DBHandler.UpdateApplicationName(token, request.NewName)
	if err != nil {
		log.Printf("error updating application: %v", err)
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusAccepted)
}
