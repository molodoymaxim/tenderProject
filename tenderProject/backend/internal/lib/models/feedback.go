package models

import (
	"time"
)

type FeedBack struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"bidId,omitempty"`
	Status      string    `json:"status,omitempty"`
	BidId       string    `json:"bidId,omitempty"`
	BidFeedback string    `json:"BidFeedback,omitempty"`
	UserName    string    `json:"UserName,omitempty"`
	Version     int       `json:"version,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}
