package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Root route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, please accept me :D")
	})

	// Applications routes
	e.POST("/applications", placeHolderHandler)
	e.GET("/applications", placeHolderHandler)
	e.GET("/applications/:token", placeHolderHandler)
	e.PUT("/applications/:token", placeHolderHandler)

	// Chats routes
	e.POST("/applications/:token/chats", placeHolderHandler)
	//get all chats for an app
	e.GET("/applications/:token/chats", placeHolderHandler)
	//get a specific chat for an app
	e.GET("/applications/:token/chats/:number", placeHolderHandler)
	e.PUT("/applications/:token/chats/:number", placeHolderHandler)

	// Messages routes
	e.POST("/applications/:token/chats/:chat_number/messages", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/messages", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/:number", placeHolderHandler)
	e.GET("/applications/:token/chats/:chat_number/messages/search", placeHolderHandler)
	e.PUT("/applications/:token/chats/:chat_number/:number", placeHolderHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

func placeHolderHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Not yet implemented")
}
