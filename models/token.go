package models

import "time"

type Token struct {
	SessionToken string     `json:"session_token" gorm:"primaryKey"`
	AccessToken  string     `json:"access_token" gorm:"unique; not null"`
	RefreshToken string     `json:"refresh_token" gorm:"unique"`
	TimeToLive   *time.Time `json:"time_to_live"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
