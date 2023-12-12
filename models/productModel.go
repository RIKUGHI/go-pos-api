package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserId    uint      `json:"user_id" gorm:"column:user_id"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:user_id;references:id"`
}
