package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) CheckUserExists(creatorUsername string) (bool, error) {
	const op = "storage.CheckUserExists"
	var Exists bool
	query := `SELECT EXISTS (SELECT 1 FROM employee WHERE username = $1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(creatorUsername).Scan(&Exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}
	return Exists, nil
}

func (s *Storage) CheckUserExistsByID(userID string) (bool, error) {
	const op = "storage.CheckUserExistsByID"
	var Exists bool
	query := `SELECT EXISTS (SELECT 1 FROM employee WHERE id = $1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(userID).Scan(&Exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}
	return Exists, nil
}

func (s *Storage) CheckOrganizationIdAndUserIDExists(organizationId, creatorUsername string) (bool, error) {
	const op = "storage.CheckOrganizationIdAndUserIDExists"
	var user_id string
	query := `SELECT id FROM employee WHERE username = $1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(creatorUsername).Scan(&user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	var exists bool
	query = `SELECT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = $1 AND user_id = $2)`

	stmt, err = s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(organizationId, user_id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}

	return exists, nil
}

func (s *Storage) GetUserIDUsername(username string) (string, error) {
	const op = "storage.CheckUserExistsByID"
	var UserdId string
	query := `SELECT id FROM employee WHERE username = $1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(username).Scan(&UserdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", fmt.Errorf("%s. Error executing query: %v", op, err)
		}
	}
	return UserdId, nil
}
