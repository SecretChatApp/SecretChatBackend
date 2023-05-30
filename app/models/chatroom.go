package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	ID        string `gorm:"size:36;not null;uniqueIndex;primary_key" json:"id"`
	UserID    string `gorm:"not null"`
	User      *User
	Message   []Message `gorm:"foregnKey:ChatRoomID" json:"messages"`
	Title     string    `gorm:"size:50" json:"title"`
	Subject   string    `gorm:"size:255" json:"subject"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (c *ChatRoom) CreateChatRoom(db *gorm.DB) (string, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return "", err
	}

	return c.Title, nil
}

func (c *ChatRoom) GetRoomInformation(id string, db *gorm.DB) error {
	err := db.Debug().Where("chat_room.id = ?", id).Joins("User").Preload("Message", func(db *gorm.DB) *gorm.DB {
		return db.Select([]string{"ChatRoomID", "Text", "Sender", "CreatedAt"}).Order("messages.created_at ASC")
	}).First(&c).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRoom) GetRoomById(id string, db *gorm.DB) error {
	err := db.Debug().Where("id = ?", id).First(&c).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRoom) UpdateRoom(title string, subject string, db *gorm.DB) error {
	err := db.Debug().Model(&c).Updates(map[string]interface{}{"title": title, "subject": subject}).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRoom) RemoveChatRoom(id string, db *gorm.DB) error {
	err := db.Debug().Where("id = ?", id).Delete(&ChatRoom{}).Error
	if err != nil {
		return err
	}

	return nil
}
