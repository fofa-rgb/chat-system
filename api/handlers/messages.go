package handlers

import (
	"bytes"
	"chat-system/internal/database"
	"chat-system/internal/models"
	"context"
	"encoding/json"
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
		return echo.ErrBadRequest
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

	messageNum, err := h.MessagesDBHandler.InsertMessage(chatId, request.Body)
	if err != nil {
		log.Printf("error inserting message: %v", err)
		return echo.ErrInternalServerError
	}

	response := &response[createMessageResponse]{Data: createMessageResponse{MessageNumber: messageNum}}

	return c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) HandleGetAllMessagesForChat(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.ErrBadRequest
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
		return echo.ErrBadRequest
	}
	messageNumber, err := parseInt64Param("message_number", c)
	if err != nil {
		return echo.ErrBadRequest
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
		return echo.ErrBadRequest
	}
	messageNumber, err := parseInt64Param("message_number", c)
	if err != nil {
		return echo.ErrBadRequest
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

func (h *MessageHandler) HandleSearchMessages(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.ErrBadRequest
	}
	chatId, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}
	request := new(searchMessageRequest)
	if err := c.Bind(request); err != nil {
		log.Printf("error binding request: %v", err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		log.Printf("error validating request: %v", err)
		return echo.ErrBadRequest
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"chat_id": chatId,
						},
					},
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"body": map[string]interface{}{
								"value": "*" + request.Query + "*",
							},
						},
					},
				},
			},
		},
	}

	reqBody, _ := json.Marshal(searchQuery)
	res, err := database.ESClient.Search(
		database.ESClient.Search.WithContext(context.Background()),
		database.ESClient.Search.WithIndex("messages"),
		database.ESClient.Search.WithBody(bytes.NewReader(reqBody)),
		database.ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer res.Body.Close()

	if res.IsError() {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": res.String()})
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to parse response"})
	}

	// Extract the desired fields from the search results
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	filteredResults := []map[string]interface{}{}

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		filteredResults = append(filteredResults, map[string]interface{}{
			"number": source["number"],
			"body":   source["body"],
		})
	}

	return c.JSON(http.StatusOK, filteredResults)
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
