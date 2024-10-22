package models

import "time"

type User struct {
	ID       uint    `gorm:"primary_key;auto_increment" json:"id"`
	Name     string  `gorm:"size:100" json:"name"`
	Email    string  `gorm:"unique;size:100" json:"email"`
	Password string  `gorm:"size:255" json:"password"`
	Balance  float64 `gorm:"type:decimal(10,2)" json:"balance"`
}

type Transaction struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	FromUserID uint      `json:"from_user_id"`
	ToUserID   uint      `json:"to_user_id"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
