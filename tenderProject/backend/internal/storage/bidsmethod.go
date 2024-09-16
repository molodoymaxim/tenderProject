package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tenderProject/backend/internal/lib/models"
)

func (s *Storage) CreateBidsDB() error {

	createType := `CREATE TYPE bids_status AS ENUM ('Created', 'Published', 'Canceled');`
	_, err := s.db.Exec(createType)
	if err != nil {
		msgErr := fmt.Errorf("bids_status have already been created", err)
		log.Println(msgErr)
	}
	createType = `CREATE TYPE author_types AS ENUM ('Organization', 'User');`
	_, err = s.db.Exec(createType)
	if err != nil {
		msgErr := fmt.Errorf("author_types have already been created", err)
		log.Println(msgErr)
	}

	createTableSQL := ` CREATE TABLE IF NOT EXISTS bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenderID UUID REFERENCES tenders(id) ON DELETE CASCADE, 
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    authorType  author_types NOT NULL,
    authorId TEXT, 
    status bids_status NOT NULL,
    version INT,
	createdAt TIMESTAMP NOT NULL 
);`
	_, err = s.db.Exec(createTableSQL)
	if err != nil {
		msgErr := fmt.Errorf("Error creating bids table:", err)
		log.Println(msgErr)
		return msgErr
	}
	return nil
}

func (s *Storage) CreateNewBid(bid models.Bid) (Id string, err error) {
	const op = "storage.CreateNewBid"

	insertQuery := `
		INSERT INTO bids (tenderID, name, description, authorType, authorId, status, version, createdAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8)
		RETURNING id;`

	stmt, err := s.db.Prepare(insertQuery)

	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	err = stmt.QueryRow(bid.TenderId, bid.Name, bid.Description, bid.AuthorType, bid.AuthorId, bid.Status, bid.Version, bid.CreatedAt).Scan(&Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	defer stmt.Close()
	return "", nil
}

// GetBids method to retrieve all user suggestions.
func (s *Storage) GetBids(bids *[]models.Bid, limit, offset int, username, companyID, tenderID string, bidSearchType int) error {
	const op = "storage.GetBids"
	var query string

	switch bidSearchType {

	case 0:
		query = "SELECT id, name, authorType, authorId, status, version, createdAt FROM bids where authorId=$1 ORDER BY name LIMIT $2 OFFSET $3"
	case 1:
		query = "SELECT id, name, authorType, authorId, status, version, createdAt FROM bids WHERE authorId in ($1, $2) ORDER BY name LIMIT $3 OFFSET $4"
	case 2:
		query = "SELECT id, name, authorType, authorId, status, version, createdAt FROM bids WHERE tenderID=$1 ORDER BY name LIMIT $2 OFFSET $3"
	default:
		return fmt.Errorf("unknown serchingType: %d", bidSearchType)
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	defer stmt.Close()

	var rows *sql.Rows
	switch bidSearchType {
	case 0:
		rows, err = stmt.Query(username, limit, offset)
	case 1:
		rows, err = stmt.Query(username, companyID, limit, offset)
	case 2:
		rows, err = stmt.Query(tenderID, limit, offset)
	}

	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		bid := models.Bid{}
		err = rows.Scan(&bid.Id, &bid.Name, &bid.AuthorType, &bid.AuthorId,
			&bid.Status, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return fmt.Errorf("%s. failed scan from database: %v", op, err)
		}
		*bids = append(*bids, bid)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("%s. rows.Next() contains errors: %v", op, err)
	}

	return nil
}

func (s *Storage) GetBid(bid *models.Bid, bidID string) (successfulRequest bool, err error) {
	const op = "storage.GetBid"

	query := "SELECT * FROM bids WHERE id=$1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(bidID).Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.AuthorType, &bid.AuthorId, &bid.Status, &bid.Version, &bid.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	return true, nil
}

func (s *Storage) CheckBidExists(bidID string) (bool, error) {
	const op = "storage.CheckBidExists"

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM bids WHERE id=$1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(bidID).Scan(&exists)
	if err != nil {
		return false, nil
	}

	return exists, nil
}

func (s *Storage) UpdateBidStatus(bidID, status string) error {
	const op = "storage.UpdateBidStatus"

	query := `UPDATE bids SET status = $1  WHERE id = $2`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(status, bidID)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

// UpdateBid a method to update the bids storage
func (s *Storage) UpdateBid(bid models.Bid) error {
	const op = "storage.UpdateBid"
	query := "UPDATE bids SET name = $1, description = $2, version = $3 WHERE id = $4"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(bid.Name, bid.Description, bid.Version, bid.Id)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}
