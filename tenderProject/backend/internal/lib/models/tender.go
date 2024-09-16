package models

import (
	"fmt"
	"time"
)

type Tender struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	ServiceType     string    `json:"serviceType,omitempty"`
	Status          string    `json:"status,omitempty"`
	OrganizationId  string    `json:"organizationId,omitempty"`
	CreatorUsername string    `json:"creatorUsername,omitempty"`
	Version         int       `json:"verstion, omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
}

func (е Tender) Validate() error {
	if len([]rune(е.Name)) > 100 {
		return fmt.Errorf("name len more than 100")
	} else if len([]rune(е.Description)) > 500 {
		return fmt.Errorf("description len more than 500")
	} else if len([]rune(е.CreatorUsername)) > 100 {
		return fmt.Errorf("description len more than 500")
	} else {
		return nil
	}
}
