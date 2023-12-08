package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"gofka/pkg/models"
	"gofka/pkg/storage"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	store *storage.NotificationStore // Reference to notification store struct.
}

/*
	 Consumer Methods
		: required : - satisfy the Sarama ConsumerGroupHandler Interface.
		: use : - { initialization & cleanups  }[message consumption] // ( just placeholders for now )
*/
func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	/* Listen to new messages on the topic

	For each message
		- STEP 1: Fetch the UserID which is the message's Key.
		- STEP 2: Unmarshal the message into a notification struct.
		- STEP 3: Add notification to store
	*/
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}

		consumer.store.Add(userID, notification)
		session.MarkMessage(msg, "")
	}

	return nil
}

func SetupConsumerGroup(ctx context.Context, store *storage.NotificationStore, consumer_group, broker_address, consumer_topic string) {
	/*  setup the kafka-consumer_group

	STEP 1: Initiliaze Consumer Group
	STEP 2: Define A consumer
	STEP 3: Consume incoming messages from kafka_topic
			- Process errors arising.
	*/

	consumerGroup, err := initiailizeConsumerGroup(broker_address, consumer_group)
	if err != nil {
		log.Printf("error from consumer \t: err : -> \n%v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}

	for {
		err = consumerGroup.Consume(ctx, []string{consumer_topic}, consumer)
		if err != nil {
			log.Printf("error from consumer \t: err : -> \n%v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func initiailizeConsumerGroup(broker_address, consumer_group string) (sarama.ConsumerGroup, error) {
	/*
		STEP 1: Initialize a new default configuration.
		STEP 2: Create a new Kafka Consumer group
					- Connect to Kafka Broker
	*/
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{broker_address}, consumer_group, config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group [%s] \n: err : -> \t%v", consumer_group, err)
	}

	return consumerGroup, nil
}
