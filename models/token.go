package models

import "time"

type TokenMap struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	SessionToken string     `json:"session_token" gorm:"index" gorm:"unique"`
	AccessToken  string     `json:"access_token" gorm:"unique"`
	CreatedAt    *time.Time `json:"created_at"`
}
