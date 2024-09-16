package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tenderProject/backend/internal/lib/models"
)

// NewVersionStorage method creates a table to store the history of tenders
func (s *Storage) NewVersionStorage() error {
	createTableHistory := `
	   CREATE TABLE IF NOT EXISTS tender_history (
	   id BIGSERIAL PRIMARY KEY, 
	   tender_id UUID REFERENCES tenders(id) ON DELETE CASCADE,
	   name VARCHAR(100) NOT NULL,
	   description VARCHAR(500) NOT NULL,
	   serviceType service_type NOT NULL,
		version INT NOT NULL
	);`

	_, err := s.db.Exec(createTableHistory)
	if err != nil {
		msgErr := fmt.Errorf("Error creating table history", err)
		log.Println(msgErr)
		return msgErr
	}

	createHash := "CREATE INDEX idx_tender_id_hash ON tender_history USING hash (tender_id);"
	_, err = s.db.Exec(createHash)
	if err != nil {
		msgErr := fmt.Errorf("Error creating create нash in tender_history  or нash was created:", err)
		log.Println(msgErr)
	}

	return nil
}

// SaveOldTender save the tender version in the History table
func (s *Storage) SaveOldTender(oldTender models.Tender) error {
	const op = "storage.SaveOldTender"

	insertQuery := `
		INSERT INTO tender_history (tender_id ,name, description, serviceType,version)
		VALUES ($1, $2, $3, $4, $5)
	`

	stmt, err := s.db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(oldTender.Id, oldTender.Name, oldTender.Description, oldTender.ServiceType, oldTender.Version)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

// SaveOldTender save the tender version in the History table
func (s *Storage) GetOldTender(tender *models.Tender, version int) error {
	const op = "storage.GetOldTender"

	insertQuery := `
		SELECT name, description, serviceType FROM tender_history
		WHERE tender_id = $1 AND version = $2
	
	`

	stmt, err := s.db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(tender.Id, version).Scan(&tender.Name, &tender.Description, &tender.ServiceType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	return nil
}
