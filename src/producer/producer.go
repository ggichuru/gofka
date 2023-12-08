/*
	- Simulate users sending notifications.
	-> Show case how the [ notifications ] are produced to a [ kafka_topic ]
*/

package main

import (
	"fmt"
	"gofka/pkg/models"
	"gofka/src/controllers"
	"gofka/src/utils"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

/* Functions */

func main() {
	/*
		STEP 1: setup producer.
		STEP 2: create a Router ( `Gin` router )
		STEP 3: define { POST } endpoint { /msgsend }
	*/

	users := []models.User{
		{ID: 1, Name: "One"},
		{ID: 2, Name: "Two"},
		{ID: 3, Name: "Three"},
		{ID: 4, Name: "Four"},
	}

	producer, err := utils.SetupProducer(KafkaServerAddress)
	if err != nil {
		log.Fatalf("Failed to intiliaize producer: %v", err)
	}

	defer producer.Close() // avoid leaks and memory loss

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/msgsend", controllers.SendMsgHttpHandler(producer, users, KafkaTopic)) //handle notifications

	fmt.Printf("[topic: %s] PRODUCER started at http://localhost%s\n", KafkaTopic, ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server : %v", err)
	}

}
