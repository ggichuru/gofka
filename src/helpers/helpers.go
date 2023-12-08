package helpers

import (
	"errors"
	"fmt"
	"gofka/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

/* Error Messages */
// helpers
var ErrUserNotFound = errors.New("user not found")
var ErrNoMessagesFound = errors.New("no messages found")

func FindUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func GetIDFromRequest(form_value string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(form_value))
	if err != nil {
		return 0, fmt.Errorf(
			"ID parsing failed. Form Value := %s => err : %w", form_value, err,
		)
	}
	return id, nil
}

func GetUserID(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}
