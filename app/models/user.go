package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string     `gorm:"size:36;not null;uniqueIndex;primary_key" json:"id"`
	ChatRoom      []ChatRoom `gorm:"foreignKey:UserID" json:"chat_room"`
	Name          string     `gorm:"size:50;not null" json:"name"`
	Email         string     `gorm:"size:50;not null" json:"email"`
	Password      string     `gorm:"size:255;not null" json:"password"`
	RememberToken string     `gorm:"size:255;not null" json:"remember_token"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (u *User) CreateUser(db *gorm.DB) (string, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

func (u *User) GerUser(db *gorm.DB, email string) error {
	err := db.Debug().Model(&User{}).Where("email = ?", email).First(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserID(db *gorm.DB, email string) (string, error) {
	err := db.Debug().Select("id").Model(&User{}).Where("email = ?", email).First(&u).Error
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

func (u *User) GetAllChatRoom(db *gorm.DB, email string) error {
	err := db.Debug().Model(&User{}).Preload("ChatRoom").Where("email = ?", email).First(&u).Error
	if err != nil {
		return err
	}

	return nil
}
