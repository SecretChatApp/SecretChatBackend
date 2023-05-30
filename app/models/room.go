package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID         string    `gorm:"size:36;not null;uniqueIndex;primary_key" json:"id"`
	ChatRoomID string    `gorm:"size:36;index" json:"chatroom_id"`
	MessageID  string    `gorm:"index"`
	Message    []Message `gorm:"" json:"messages"`
	Link       string    `gorm:"size:255;not null" json:"link"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (r *Room) RemoveRoom(id string, db *gorm.DB) error {
	err := db.Debug().Where("id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
