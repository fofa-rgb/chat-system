package handlers

import (
	"bytes"
	"chat-system/internal/database"
	"chat-system/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MessageHandlers struct {
	MessagesDBHandler     *database.MessagesDatabaseHandler
	ChatsDBHandler        *database.ChatsDatabaseHandler
	ApplicationsDBHandler *database.ApplicationsDatabaseHandler
	WriteQueue            chan MessageWriteRequest
	UpdateQueue           chan MessageUpdateRequest
	TaskStatusMap         map[string]*MessageTaskStatus
}

type MessageWriteRequest struct {
	TaskID      string
	ChatID      int64
	MessageBody string
}
type MessageUpdateRequest struct {
	TaskID        string
	MessageNumber int64
	ChatID        int64
	NewBody       string
}

type MessageTaskStatus struct {
	Status string // "Pending", "Completed", "Error"
	models.UserExposedMessage
	Error string
}

func CreateMessageHandlers() *MessageHandlers {
	messagesDbHandler := database.NewMessagesDatabaseHandler()
	chatsDbHandler := database.NewChatsDatabaseHandler()
	applicationsDbHandler := database.NewApplicationsDatabaseHandler()

	handler := &MessageHandlers{
		MessagesDBHandler:     messagesDbHandler,
		ChatsDBHandler:        chatsDbHandler,
		ApplicationsDBHandler: applicationsDbHandler,
		WriteQueue:            make(chan MessageWriteRequest, 1000),
		UpdateQueue:           make(chan MessageUpdateRequest, 1000),
		TaskStatusMap:         make(map[string]*MessageTaskStatus),
	}

	go handler.startWorker()

	return handler
}

func (h *MessageHandlers) startWorker() {
	for {
		select {
		case createReq := <-h.WriteQueue:
			status := h.TaskStatusMap[createReq.TaskID]
			messageNum, err := h.MessagesDBHandler.InsertMessage(createReq.ChatID, createReq.MessageBody)
			if err != nil {
				log.Printf("Error inserting message: %v", err)
				status.Status = "Error"
				status.Error = "Failed to create message"
			} else {
				status.Status = "Completed"
				status.Number = messageNum
				status.Body = createReq.MessageBody
			}
		case updateReq := <-h.UpdateQueue:
			status := h.TaskStatusMap[updateReq.TaskID]
			newMessage, err := h.MessagesDBHandler.UpdateMessageBody(updateReq.ChatID, updateReq.MessageNumber, updateReq.NewBody)
			if err != nil {
				log.Printf("Error updating chat: %v", err)
				status.Status = "Error"
				status.Error = "Failed to update chat subject"
			} else {
				status.Status = "Completed"
				status.Number = newMessage.Number
				status.Body = newMessage.Body
			}
		}
	}
}

func (h *MessageHandlers) HandleCreateMessage(c echo.Context) error {
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

	chatID, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}

	// Generate a unique task ID
	taskID := uuid.New().String()

	h.TaskStatusMap[taskID] = &MessageTaskStatus{
		Status: "Pending",
	}

	// Push the request to the queue
	h.WriteQueue <- MessageWriteRequest{
		TaskID:      taskID,
		ChatID:      chatID,
		MessageBody: request.Body,
	}

	// Respond with the status-check URL
	statusURL := c.Scheme() + "://" + c.Request().Host + "/messages/status/" + taskID
	return c.JSON(http.StatusAccepted, map[string]string{
		"status_url": statusURL,
	})
}

func (h *MessageHandlers) HandleGetMessageStatus(c echo.Context) error {
	taskID := c.Param("taskID")

	status, exists := h.TaskStatusMap[taskID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	return c.JSON(http.StatusOK, status)
}

func (h *MessageHandlers) HandleGetAllMessagesForChat(c echo.Context) error {
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

func (h *MessageHandlers) HandleGetMessage(c echo.Context) error {
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

func (h *MessageHandlers) HandleUpdateMessageBody(c echo.Context) error {
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

	chatID, err := h.getChatIdFromAppTokenAndChatNumber(token, chatNumber)
	if err != nil {
		log.Printf("error getting chat id: %v", err)
		return echo.ErrInternalServerError
	}

	taskID := uuid.New().String()

	h.TaskStatusMap[taskID] = &MessageTaskStatus{
		Status: "Pending",
	}

	// Push the update request to the UpdateQueue
	h.UpdateQueue <- MessageUpdateRequest{
		TaskID:        taskID,
		ChatID:        chatID,
		MessageNumber: messageNumber,
		NewBody:       request.NewBody,
	}

	// Respond with the status-check URL
	statusURL := c.Scheme() + "://" + c.Request().Host + "/messages/status/" + taskID
	return c.JSON(http.StatusAccepted, map[string]string{
		"status_url": statusURL,
	})
}

func (h *MessageHandlers) HandleSearchMessages(c echo.Context) error {
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
								"value": "*" + request.Query + "*", // Matching any part of the string
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

func (h *MessageHandlers) getChatIdFromAppTokenAndChatNumber(token string, chatNumber int64) (int64, error) {
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
