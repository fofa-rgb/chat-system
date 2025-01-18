package handlers

import (
	"chat-system/internal/database"
	"chat-system/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	MessagesDBHandler     *database.MessagesDatabaseHandler
	ChatsDBHandler        *database.ChatsDatabaseHandler
	ApplicationsDBHandler *database.ApplicationsDatabaseHandler
}

func CreateMessageHandlers() *MessageHandler {
	messagesDbHandler := database.NewMessagesDatabaseHandler()
	chatsDbHandler := database.NewChatsDatabaseHandler()
	applicationsDbHandler := database.NewApplicationsDatabaseHandler()

	return &MessageHandler{MessagesDBHandler: messagesDbHandler, ChatsDBHandler: chatsDbHandler, ApplicationsDBHandler: applicationsDbHandler}
}

func (h *MessageHandler) HandleCreateMessage(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	request := new(createMessageRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}

	chatId, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}

	chatNum, err := h.MessagesDBHandler.InsertMessage(chatId, request.Body)
	if err != nil {
		log.Printf("error inserting message: %v", err)
		return echo.ErrInternalServerError
	}

	response := &response[createChatResponse]{Data: createChatResponse{ChatNumber: chatNum}}

	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) HandleGetAllMessagesForChat(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	chatId, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}
	messages, err := h.MessagesDBHandler.GetAllMessagesForAChat(chatId)
	if err != nil {
		log.Printf("error getting messages: %v", err)
		return echo.ErrInternalServerError
	}
	var userExposedMessages []models.UserExposedMessage
	for _, message := range messages {
		userExposedMessages = append(userExposedMessages, models.UserExposedMessage{
			Number: message.Number,
			Body:   message.Body,
		})
	}

	response := &response[[]models.UserExposedMessage]{Data: userExposedMessages}

	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) HandleGetMessage(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	messageNumber, err := parseInt64Param("message_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	chatId, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}

	message, err := h.MessagesDBHandler.GetMessageByChatIdAndMessageNumber(chatId, messageNumber)
	if err != nil {
		log.Printf("error getting message: %v", err)
		return echo.ErrInternalServerError
	}

	userExposedMessage := models.UserExposedMessage{
		Number: message.Number,
		Body:   message.Body,
	}

	response := &response[models.UserExposedMessage]{Data: userExposedMessage}

	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) HandleUpdateMessageBody(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	messageNumber, err := parseInt64Param("message_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	request := new(updateMessageRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}

	chatId, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}

	updatedMessage, err := h.MessagesDBHandler.UpdateMessageBody(chatId, messageNumber, request.NewBody)
	if err != nil {
		log.Printf("error updating message: %v", err)
		return echo.ErrInternalServerError
	}

	userExposedMessage := models.UserExposedMessage{
		Body:   updatedMessage.Body,
		Number: updatedMessage.Number,
	}

	response := &response[models.UserExposedMessage]{Data: userExposedMessage}

	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) getChatIdFromAppTokenAndChatNumber(token string, chatNumber int64) (int64, error) {
	applicationId, err := h.ApplicationsDBHandler.GetApplicationIdByToken(token)
	if err != nil {
		log.Printf("error getting app id: %v", err)
		return 0, err
	}
	chatId, err := h.ChatsDBHandler.GetChatIdByAppIdAndChatNumber(applicationId, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return 0, err
	}
	return chatId, nil
}
