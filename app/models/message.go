package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         string `gorm:"size:36;not null; uniqueIndex;primary_key" json:"id"`
	ChatRoomID string `gorm:"not null" json:"chatroom_id"`
	Text       string `gorm:"type:text" json:"text"`
	Sender     string `gorm:"size:255" json:"sender"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (m *Message) CreateMessage(db *gorm.DB) error {
	err := db.Debug().Create(&m).Error
	if err != nil {
		return err
	}

	return nil
}
