package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tenderProject/backend/internal/lib/models"
)

// CreateBidStory IF NOT EXISTS
func (s Storage) CreateBidStory() error {

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS bids_version (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    budID UUID REFERENCES bids(id) ON DELETE CASCADE, 
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    version INT NOT NULL
);
`
	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		msgErr := fmt.Errorf("Error creating bids story table:", err)
		log.Println(msgErr)
		return msgErr
	}
	return nil
}

func (s *Storage) SaveOldBid(bid models.Bid) error {
	const op = "storage.SaveOldBid"

	insertQuery := `
		INSERT INTO bids_version (budID ,name, description, version)
		VALUES ($1, $2, $3, $4)
	`

	stmt, err := s.db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(bid.Id, bid.Name, bid.Description, bid.Version)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

func (s Storage) GetOldBid(bid *models.Bid, version int) error {
	const op = "storage.GetOldBid"

	insertQuery := `
		SELECT name, description FROM tender_history
		WHERE budID = $1 AND version = $2
	`

	stmt, err := s.db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(bid.Id, version).Scan(&bid.Name, &bid.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	return nil
}
