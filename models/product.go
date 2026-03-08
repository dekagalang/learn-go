package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
