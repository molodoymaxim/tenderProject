package storage

import "fmt"

func (s *Storage) GetCompanyIDbyUser(username string) (companyId string, err error) {
	const op = "storage.GetCompanyIDbyUser"
	query := `
	SELECT o.id AS organization_id 
	FROM employee e
	JOIN organization_responsible orr ON e.id = orr.user_id
	JOIN organization o ON orr.organization_id = o.id
	WHERE e.username = $1;`

	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(username).Scan(&companyId)
	if err != nil {
		return "", fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return companyId, nil
}

//const op = "storage.SaveOldTender"
//
//insertQuery := `
//		INSERT INTO tender_history (tender_id ,name, description, serviceType,version)
//		VALUES ($1, $2, $3, $4, $5)
//	`
//
//stmt, err := s.db.Prepare(insertQuery)
//defer stmt.Close()
//if err != nil {
//return fmt.Errorf("%s. Error preparing statement: %v", op, err)
//}
//
//_, err = stmt.Exec(oldTender.Id, oldTender.Name, oldTender.Description, oldTender.ServiceType, oldTender.Version)
//if err != nil {
//return fmt.Errorf("%s. Error executing query: %v", op, err)
//}
//
//return nil
