package cron

import (
	"chat-system/internal/database"
	"log"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	applicationDBHandler *database.ApplicationsDatabaseHandler
	chatsDBHandler       *database.ChatsDatabaseHandler
}

func NewCronJob() *CronJob {
	appDBHandler := database.NewApplicationsDatabaseHandler()
	chatDBHandler := database.NewChatsDatabaseHandler()
	return &CronJob{applicationDBHandler: appDBHandler, chatsDBHandler: chatDBHandler}
}

func (cj *CronJob) Start() {
	c := cron.New()

	_, err := c.AddFunc("@every 50m", func() {
		log.Println("Cron job started: Updating stats...")
		err := cj.applicationDBHandler.UpdateChatsCount()
		if err != nil {
			log.Printf("Error updating chats_count: %v\n", err)
		}
		err = cj.chatsDBHandler.UpdateMessagesCount()
		if err != nil {
			log.Printf("Error updating messages_count: %v\n", err)
		}
		log.Println("Cron job completed.")
	})
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v\n", err)
	}

	c.Start()
	log.Println("Cron scheduler started.")
}
