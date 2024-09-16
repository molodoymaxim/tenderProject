package storage

import (
	"fmt"
	"log"
)

func (s *Storage) CreateFeedbackStorage() error {

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS feedback_storage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bidId UUID REFERENCES bids(id) ON DELETE CASCADE, 
    name VARCHAR(50) NOT NULL,
	status bids_status NOT NULL,
	authorType author_types NOT NULL,
	authorId VARCHAR(50) NOT NULL,
	bidFeedback VARCHAR(1000) NOT NULL,
	username VARCHAR(50) NOT NULL,
    version INT NOT NULL,
	createdAt TIMESTAMP NOT NULL 
);
`
	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		msgErr := fmt.Errorf("Error creating bids table:", err)
		log.Println(msgErr)
		return msgErr
	}
	return nil
}

//func (s *Storage) NewFeedback(feedback models.FeedBack) (string, error) {
//	const op = "storage.NewFeedback"
//
//	insertQuery := `
//		INSERT INTO feedback_storage (name, description, serviceType, status, organizationId, creatorUsername,version, createdAt)
//		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//		RETURNING id;
//	`
//
//	stmt, err := s.db.Prepare(insertQuery)
//	defer stmt.Close()
//	if err != nil {
//		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
//	}
//
//	var ID string
//	err = stmt.QueryRow(tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatorUsername, tender.Version, tender.CreatedAt).Scan(&ID)
//	if err != nil {
//		return "", fmt.Errorf("%s. Error executing query: %v", op, err)
//	}
//
//	return ID, nil
//}
