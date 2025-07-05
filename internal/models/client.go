package models

import "time"

type Client struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	ChannelID   int        `json:"channel_id"`
	Channel     string     `json:"channel,omitempty"`
	Bonus       int        `json:"bonus"`
	Visits      int        `json:"visits"`
        Income      int        `json:"income"`
        Status      string     `json:"status"`
        CreatedAt   time.Time  `json:"created_at"`
        UpdatedAt   time.Time  `json:"updated_at"`
}
