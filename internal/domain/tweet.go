package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
}

func NewTweet(userID, content string) *Tweet {
	return &Tweet{
		ID:        uuid.NewString(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
}
