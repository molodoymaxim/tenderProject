package models

import (
	"fmt"
	"time"
)

type Bid struct {
	Id          string    `json:"id,omitempty,omitempty"`
	TenderId    string    `json:"tenderId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	AuthorType  string    `json:"authorType,omitempty"`
	AuthorId    string    `json:"authorId,omitempty"`
	Status      string    `json:"status,omitempty"`
	Version     int       `json:"version,omitempty"`
	CreatedAt   time.Time `json:"CreatedAt,omitempty"`
}

func (b Bid) Validate() error {
	if len([]rune(b.Name)) > 100 {
		return fmt.Errorf("name len more than 100")
	} else if len([]rune(b.Description)) > 500 {
		return fmt.Errorf("description len more than 500")
	} else if len([]rune(b.AuthorId)) > 100 {
		return fmt.Errorf("description len more than 500")
	} else {
		return nil
	}
}
