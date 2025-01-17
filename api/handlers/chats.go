package handlers

import (
	"chat-system/internal/database"
	"chat-system/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ChatHandlers struct {
	ChatsDBHandler        *database.ChatsDatabaseHandler
	ApplicationsDBHandler *database.ApplicationsDatabaseHandler
}

func CreateChatHandlers() *ChatHandlers {
	chatsDbHandler := database.NewChatsDatabaseHandler()
	applicationsDbHandler := database.NewChatsDatabaseHandler()

	return &ChatHandlers{ChatsDBHandler: chatsDbHandler, ApplicationsDBHandler: (*database.ApplicationsDatabaseHandler)(applicationsDbHandler)}
}

func (h *ChatHandlers) HandleCreateChat(c echo.Context) error {
	token := c.Param("token")
	request := new(createChatRequest)
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
		log.Printf("error inserting chat: %v", err)
		return echo.ErrInternalServerError
	}

	response := &response[createChatResponse]{Data: createChatResponse{ChatNumber: chatNum}}

	return c.JSON(http.StatusOK, response)
}

func (h *ChatHandlers) HandleGetAllChatsForApplication(c echo.Context) error {
	token := c.Param("token")
	applicationId, err := h.ApplicationsDBHandler.GetApplicationIdByToken(token)
	if err != nil {
		log.Printf("error getting app id: %v", err)
		return echo.ErrInternalServerError
	}

	chats, err := h.ChatsDBHandler.GetAllChatsForAnApp(applicationId)
	if err != nil {
		log.Printf("error getting chat: %v", err)
		return echo.ErrInternalServerError
	}
	var userExposedChats []models.UserExposedChat
	for _, chat := range chats {
		userExposedChats = append(userExposedChats, models.UserExposedChat{
			Subject:       chat.Subject,
			Number:        chat.Number,
			MessagesCount: chat.MessagesCount,
		})
	}

	response := &response[[]models.UserExposedChat]{Data: userExposedChats}

	return c.JSON(http.StatusOK, response)
}

func (h *ChatHandlers) HandleGetChat(c echo.Context) error {
	token := c.Param("token")
	chatNumberStr := c.Param("number")

	// Convert chatNumber to int64
	chatNumber, err := strconv.ParseInt(chatNumberStr, 10, 64)
	if err != nil {
		log.Printf("error parsing chat number: %v", err)
		return echo.ErrBadRequest
	}

	applicationId, err := h.ApplicationsDBHandler.GetApplicationIdByToken(token)
	if err != nil {
		log.Printf("error getting app id: %v", err)
		return echo.ErrInternalServerError
	}

	chat, err := h.ChatsDBHandler.GetChatByApplicationIdAndChatNumber(applicationId, chatNumber)
	if err != nil {
		log.Printf("error getting chat: %v", err)
		return echo.ErrInternalServerError
	}

	userExposedChat := models.UserExposedChat{
		Subject:       chat.Subject,
		Number:        chat.Number,
		MessagesCount: chat.MessagesCount,
	}

	response := &response[models.UserExposedChat]{Data: userExposedChat}

	return c.JSON(http.StatusOK, response)
}
