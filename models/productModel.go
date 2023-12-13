package models

import "time"

type Product struct {
	ID        uint   `gorm:"primarykey"`
	UserId    uint   `gorm:"column:user_id"`
	Name      string `gorm:"type:varchar(255)"`
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"foreignKey:user_id;references:id"`
}
