package utils

import (
	"encoding/json"
	"fmt"
	"gofka/pkg/models"
	"gofka/src/helpers"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// Setup a producer
func SetupProducer(broker_address string) (sarama.AsyncProducer, error) {
	/*
		STEP 1: setup parameters before connecting to kafka broker.
		STEP 2: ensure producer recieves acknowledgment of storage of message to topic.
		// STEP 3: initialize a new sync Kafka producer. [connect -> kafka broker]:{localhost:9092} ( from the Docker Compose Earlier. )
		StEP 3: Initialize a new async Kafka producer _> High  througput and efficiency.
	*/
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//// producer, err := sarama.NewSyncProducer([]string{broker_address}, config)
	// Initilialize a producer that connects to kafka broker
	producer, err := sarama.NewAsyncProducer([]string{broker_address}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer :-> %w", err)
	}
	return producer, nil
}

// Publish A Message to a Kafka Topic
func ProduceMessage(producer sarama.AsyncProducer, users []models.User, ctx *gin.Context, fromID, toID int, kafkaTopic string) error {
	/* [ SETTING UP A KAFKA HANDLER ] : { producing to a kafka topic }
	STEP 1: -> Get message from context. (in this case, get from post request)
	STEP 2: -> Get Both users by their ID's from [ db or store ]
	STEP 3: -> Encapsulate [content-info] into a json object. [json(topic struct)] - (notifications in this case - check models)
	STEP 4: -> Construct a Producer message. [topic, key, value ] {notifications, recepientID-enc, serialized_NotificationObject-enc}
	STEP 5: -> Send the message to Kafka topic  through channel
	*/

	// 1
	message := ctx.PostForm("message")

	// 2
	fromUser, err := helpers.FindUserByID(fromID, users)
	if err != nil {
		return err
	}
	toUser, err := helpers.FindUserByID(toID, users)
	if err != nil {
		return err
	}

	// 3
	// Initiailize notification struct
	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}
	// Convert it into a JSON object
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// 4
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	// 5
	producer.Input() <- msg

	// Read Form Success channel to prevent deadlock
	producer.Successes()

	return err
}
