package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         string    `gorm:"size:36;not null; uniqueIndex;primary_key" json:"id"`
	ChatRoomID string    `gorm:"not null" json:"chatroom_id"`
	Text       string    `gorm:"type:text" json:"text"`
	Sender     string    `gorm:"size:255" json:"sender"`
	CreatedAt  time.Time `json:"created_at"`
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

func (m *Message) GetAllMessagesByChatRoomId(db *gorm.DB, chatRoomId string) ([]Message, error) {
	var messages []Message
	err := db.Debug().Where("chat_room_id = ?", chatRoomId).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return messages, err
	}

	return messages, nil
}
