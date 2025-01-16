package main

import (
	"chat-system/api/handlers"
	"chat-system/internal/database"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func placeHolderHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Not yet implemented")
}

func main() {
	// Database setup
	database.InitDB()
	database.ESClientConnection()
	database.ESCreateIndexIfNotExist()

	e := echo.New()
	appHandlers := &handlers.ApplicationHandlers{}

	// Root route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, please accept me :D")
	})

	// Applications routes
	e.POST("/applications", appHandlers.HandleCreateApplication)
	e.GET("/applications", placeHolderHandler)
	e.GET("/applications/:token", placeHolderHandler)
	e.PUT("/applications/:token", placeHolderHandler)

	// Chats routes
	e.POST("/applications/:token/chats", placeHolderHandler)
	e.GET("/applications/:token/chats", placeHolderHandler)
	e.GET("/applications/:token/chats/:number", placeHolderHandler)
	e.PUT("/applications/:token/chats/:number", placeHolderHandler)

	// Messages routes
	e.POST("/applications/:token/chats/:chat_number/messages", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/messages", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/:number", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/messages/search", placeHolderHandler)
	e.PUT("/applications/:token/chats/:chat_number/:number", placeHolderHandler)
	e.POST("/applications/:token/chats/:chat_number/messages/index", placeHolderHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
