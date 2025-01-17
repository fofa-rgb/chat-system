package handlers

//general
type response[T any] struct {
	Data T `json:"data"`
}

//applications
type createApplicationRequest struct {
	Name string `json:"name" validate:"required"`
}

type createApplicationResponse struct {
	Token string `json:"token"`
}
type updateApplicationNameRequest struct {
	NewName string `json:"newName"`
}

//chats
type createChatRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type createChatResponse struct {
	ChatNumber int64 `json:"chatNumber" validate:"required"`
}
