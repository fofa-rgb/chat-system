package handlers

import (
	"chat-system/internal/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatHandlers struct {
	ChatsDBHandler        *database.ChatsDatabaseHandler
	ApplicationsDBHandler *database.ApplicationsDatabaseHandler
}

type createChatRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type createChatResponse struct {
	ChatNumber int64 `json:"chatNumber" validate:"required"`
}

func CreateChatHandlers() *ChatHandlers {
	chatsDbHandler := database.NewChatsDatabaseHandler()
	applicationsDbHandler := database.NewChatsDatabaseHandler()

	return &ChatHandlers{ChatsDBHandler: chatsDbHandler, ApplicationsDBHandler: (*database.ApplicationsDatabaseHandler)(applicationsDbHandler)}
}

func (h *ChatHandlers) HandleCreateChat(c echo.Context) error {
	request := new(createChatRequest)
	token := c.Param("token")
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}

	applicationId, err := h.ApplicationsDBHandler.GetApplicationIdByToken(token)
	if err != nil {
		log.Printf("error getting application id: %v", err)
		return echo.ErrInternalServerError
	}

	chatNum, err := h.ChatsDBHandler.InsertChat(applicationId, request.Subject)
	if err != nil {
		log.Printf("error inserting application: %v", err)
		return echo.ErrInternalServerError
	}

	response := &Response[createChatResponse]{Data: createChatResponse{ChatNumber: chatNum}}

	return c.JSON(http.StatusOK, response)
}
