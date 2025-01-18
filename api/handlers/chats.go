package handlers

import (
	"chat-system/internal/database"
	"chat-system/internal/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ChatHandlers struct {
	ChatsDBHandler        *database.ChatsDatabaseHandler
	ApplicationsDBHandler *database.ApplicationsDatabaseHandler
	WriteQueue            chan ChatWriteRequest
	TaskStatusMap         map[string]*ChatTaskStatus
}
type ChatTaskStatus struct {
	Status     string // "Pending", "Completed", "Error"
	ChatNumber int64
	Error      string
}

type ChatWriteRequest struct {
	TaskID        string
	ApplicationID int64
	Subject       string
}

func CreateChatHandlers() *ChatHandlers {
	chatsDbHandler := database.NewChatsDatabaseHandler()
	applicationsDbHandler := database.NewApplicationsDatabaseHandler()

	handler := &ChatHandlers{
		ChatsDBHandler:        chatsDbHandler,
		ApplicationsDBHandler: applicationsDbHandler,
		WriteQueue:            make(chan ChatWriteRequest, 1000), // Adjust buffer size as needed
		TaskStatusMap:         make(map[string]*ChatTaskStatus),
	}

	// Start the background worker
	go handler.startWorker()

	return handler
}

func (h *ChatHandlers) startWorker() {
	for req := range h.WriteQueue {
		status := h.TaskStatusMap[req.TaskID]
		chatNum, err := h.ChatsDBHandler.InsertChat(req.ApplicationID, req.Subject)
		if err != nil {
			log.Printf("Error inserting chat: %v", err)
			status.Status = "Error"
			status.Error = "Failed to create chat"
		} else {
			status.Status = "Completed"
			status.ChatNumber = chatNum
		}
	}
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

	taskID := uuid.New().String()

	h.TaskStatusMap[taskID] = &ChatTaskStatus{
		Status: "Pending",
	}

	// Push the request to the queue
	h.WriteQueue <- ChatWriteRequest{
		TaskID:        taskID,
		ApplicationID: applicationId,
		Subject:       request.Subject,
	}

	// Respond with the updated status-check URL
	statusURL := c.Scheme() + "://" + c.Request().Host + "/chats/status/" + taskID
	return c.JSON(http.StatusAccepted, map[string]string{
		"status_url": statusURL,
	})
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
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

func (h *ChatHandlers) HandleUpdateChatSubject(c echo.Context) error {
	token := c.Param("token")
	chatNumber, err := parseInt64Param("chat_number", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	request := new(updateChatRequest)
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
		log.Printf("error getting app id: %v", err)
		return echo.ErrInternalServerError
	}

	chat, err := h.ChatsDBHandler.UpdateChatSubject(applicationId, chatNumber, request.NewSubject)
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

func (h *ChatHandlers) HandleGetStatus(c echo.Context) error {
	taskID := c.Param("taskID")

	status, exists := h.TaskStatusMap[taskID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      status.Status,
		"chat_number": status.ChatNumber,
		"error":       status.Error,
	})
}
