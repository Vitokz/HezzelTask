package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Log struct {
	LogText   string    `json:"log_text"`
	EventName string    `json:"event_name"`
	EventDate time.Time `json:"event_date"`
}
