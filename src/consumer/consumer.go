/* SETUP CONSUMER ->
- {
	-> Listens to [notifications] topic,
	-> Provides endpoint to list notifications for a user
	}
*/

package main

import (
	"context"
	"fmt"
	"gofka/pkg/storage"
	"gofka/src/controllers"
	"gofka/src/utils"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

func main() {
	// TODO : Fix this bug

	// store := &storage.NotificationStore{
	// 	data: make(storage.UserNotifications),
	// }
	store := &storage.NotificationStore{}

	ctx, cancel := context.WithCancel(context.Background())
	go utils.SetupConsumerGroup(ctx, store, ConsumerGroup, KafkaServerAddress, ConsumerTopic)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		controllers.ViewMsgsHttpHandler(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
