package main

import (
	"chat-system/api/cron"
	"chat-system/api/handlers"
	"chat-system/internal/database"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func placeHolderHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Not yet implemented")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Database setup
	database.InitDB()
	database.ESClientConnection()
	database.ESCreateIndexIfNotExist()

	// Start the cron job in a Goroutine
	go func() {
		cronJob := cron.NewCronJob()
		cronJob.Start()
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	v := validator.New()
	e.Validator = &CustomValidator{validator: v}

	appHandlers := handlers.CreateApplicationHandlers()
	chatHandlers := handlers.CreateChatHandlers()
	messageHandlers := handlers.CreateMessageHandlers()

	// Root route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, please accept me :D")
	})

	// Applications routes
	e.POST("/applications", appHandlers.HandleCreateApplication)
	e.GET("/applications", appHandlers.HandleGetAllApplications)
	e.GET("/applications/:token", appHandlers.HandleGetApplicationByToken)
	e.PATCH("/applications/:token", appHandlers.HandleUpdateApplicationName)

	// Chats routes
	e.POST("/applications/:token/chats", chatHandlers.HandleCreateChat)
	e.GET("/applications/:token/chats", chatHandlers.HandleGetAllChatsForApplication)
	e.GET("/applications/:token/chats/:chat_number", chatHandlers.HandleGetChat)
	e.PATCH("/applications/:token/chats/:chat_number", chatHandlers.HandleUpdateChatSubject)

	// Messages routes
	e.POST("/applications/:token/chats/:chat_number/messages", messageHandlers.HandleCreateMessage)
	e.GET("/applications/:token/chats/:chat_number/messages", messageHandlers.HandleGetAllMessagesForChat)
	e.GET("/applications/:token/chats/:chat_number/messages/:message_number", messageHandlers.HandleGetMessage)
	e.PATCH("/applications/:token/chats/:chat_number/messages/:message_number", messageHandlers.HandleUpdateMessageBody)

	//message queue status routes
	e.GET("/chats/status/:taskID", chatHandlers.HandleGetStatus)
	e.GET("/messages/status/:taskID", messageHandlers.HandleGetMessageStatus)

	//elastic search messages
	e.GET("/applications/:token/chats/:chat_number/messages/search", messageHandlers.HandleSearchMessages)
	e.POST("/applications/:token/chats/:chat_number/messages/index", placeHolderHandler)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
