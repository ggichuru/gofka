package controllers

import (
	"errors"
	"gofka/pkg/models"
	"gofka/pkg/storage"
	"gofka/src/helpers"
	"gofka/src/utils"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// Route Controller to handle Http Request to send a message -> { /msgsend }
func SendMsgHttpHandler(producer sarama.AsyncProducer, users []models.User, kafkaTopic string) gin.HandlerFunc {
	/* Process incoming requests { while listening to http responses }
	STEP 1: ensure valid sender and recipient ID's are provided
	STEP 2: Produce message to a kafka topic.
			a -: If no err, send headers for Success.
	*/
	return func(ctx *gin.Context) {
		// 1
		fromID, err := helpers.GetIDFromRequest("fromID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		toID, err := helpers.GetIDFromRequest("toID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// 2
		err = utils.ProduceMessage(producer, users, ctx, fromID, toID, kafkaTopic)
		if errors.Is(err, helpers.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found!"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// 2-:a
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

func ViewMsgsHttpHandler(ctx *gin.Context, store *storage.NotificationStore) {
	userID, err := helpers.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message":       "No notifications found for user",
				"notifications": []models.Notification{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}
