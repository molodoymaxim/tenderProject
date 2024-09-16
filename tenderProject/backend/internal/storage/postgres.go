package storage

import (
	"database/sql"
	"fmt"
	"log"
	"tenderProject/backend/internal/config"
	"tenderProject/backend/internal/lib/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func ConnectToStorage(config *config.Config, isLocal bool) *Storage {
	var connStr string
	if isLocal {
		fmt.Println("Try local")
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			config.PorstgresUserName, config.PorstgresPassword, config.PorstgresDatabase)
	} else {
		connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
			config.PorstgresUserName, config.PorstgresPassword, config.PorstgresHost, config.PorstgresPort, config.PorstgresDatabase)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error during connection verification: %v", err)
	}

	log.Println("Connection to the database was successful")
	return &Storage{db: db}
}

// CreateTables method for adding the required databases
func (s *Storage) CreateTables() {
	err := s.NewTenderStorage()
	if err != nil {
		log.Fatal("failed to create tender storage")
	}
	err = s.NewVersionStorage()
	if err != nil {
		log.Fatal("failed to create version storage")
	}
	err = s.CreateBidsDB()
	if err != nil {
		log.Fatal("failed to create version storage")
	}

	err = s.CreateBidStory()
	if err != nil {
		log.Fatal("failed to create bid story storage")
	}

	//storageData.CreateRelation()

}

// NewTenderStorage create new tender Storage
func (s *Storage) NewTenderStorage() error {

	createType := `
	CREATE TYPE service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
`
	_, err := s.db.Exec(createType)
	if err != nil {
		msgErr := fmt.Errorf("service_type have already been created", err)
		log.Println(msgErr)
	}

	createType = `
	CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
`
	_, err = s.db.Exec(createType)
	if err != nil {
		msgErr := fmt.Errorf("tender_status have already been created", err)
		log.Println(msgErr)
	}

	createTableSQL := `
  	CREATE TABLE IF NOT EXISTS tenders (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		description VARCHAR(500) NOT NULL,
		serviceType service_type NOT NULL,
		status tender_status NOT NULL,
	    organizationId UUID REFERENCES organization(id) ON DELETE CASCADE,
		creatorUsername VARCHAR(50) REFERENCES employee(username),
		version INT DEFAULT 1 NOT NULL,
		createdAt TIMESTAMP NOT NULL
	 );`

	_, err = s.db.Exec(createTableSQL)
	if err != nil {
		msgErr := fmt.Errorf("Error creating table:", err)
		log.Println(msgErr)
		return msgErr
	}

	return nil
}

// Close connection to BD
func (s *Storage) Close() {
	s.db.Close()
}

// GetTenders Receipt of all published tenders
func (s *Storage) GetTenders(tenders *[]models.Tender, limit, offset int, searchInfo string, serchingType int) error {
	const op = "storage.GetTenders"
	status := "Published"
	var query string
	switch serchingType {
	case 0:
		query = "SELECT id, name, description, serviceType, status, version, createdAt FROM tenders where status=$1 ORDER BY name LIMIT $2 OFFSET $3"
	case 1:
		query = "SELECT id, name, description, serviceType, status, version, createdAt FROM tenders WHERE status=$1 and serviceType = $2 ORDER BY name LIMIT $3 OFFSET $4"
	case 2:
		query = "SELECT id, name, description, serviceType, status, version, createdAt FROM tenders WHERE  creatorUsername = $1 ORDER BY name LIMIT $2 OFFSET $3"
	default:
		return fmt.Errorf("unknown serchingType: %d", serchingType)
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	defer stmt.Close()

	var rows *sql.Rows
	switch serchingType {
	case 0:
		rows, err = stmt.Query(status, limit, offset)
	case 1:
		rows, err = stmt.Query(status, searchInfo, limit, offset)
	case 2:
		rows, err = stmt.Query(searchInfo, limit, offset)
	}

	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType,
			&tender.Status, &tender.Version, &tender.CreatedAt)
		if err != nil {
			return fmt.Errorf("%s. failed scan from database: %v", op, err)
		}
		*tenders = append(*tenders, tender)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("%s. rows.Next() contains errors: %v", op, err)
	}

	return nil
}

// CreateTender Creating a new tender
func (s *Storage) CreateTender(tender models.Tender) (string, error) {
	const op = "storage.CreateTender"

	insertQuery := `
		INSERT INTO tenders (name, description, serviceType, status, organizationId, creatorUsername,version, createdAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	stmt, err := s.db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	var ID string
	err = stmt.QueryRow(tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatorUsername, tender.Version, tender.CreatedAt).Scan(&ID)
	if err != nil {
		return "", fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return ID, nil
}

// CheckOrganizationExists check if the Organizatio exists in the Organizatio database. Return TRUE if User EXISTS
func (s *Storage) CheckOrganizationExists(organizationid string) (bool, error) {
	const op = "storage.CheckCompanyExists"
	var Exists bool
	query := `SELECT EXISTS (SELECT 1 FROM organization WHERE id = $1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(organizationid).Scan(&Exists)
	if err != nil {
		return false, nil
	}
	return Exists, nil
}

// CheckTenderExists check if the tender exists in the tenders database. Return TRUE if User EXISTS
func (s *Storage) CheckTenderExists(tenderId string) (bool, error) {
	const op = "storage.CheckTenderExists"
	var Exists bool
	query := `SELECT EXISTS (SELECT 1 FROM tenders WHERE id = $1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(tenderId).Scan(&Exists)
	if err != nil {
		return false, nil
	}
	return Exists, nil
}

// GetTenderStatus return Tender Status from tenders DB
func (s *Storage) GetTenderStatus(id string) (status string, err error) {
	const op = "storage.GetTenderStatusAndCompanyID"

	query := "SELECT status FROM tenders WHERE id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(id).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("%s. Error executing query: %v", op, err)
	}
	return status, nil
}

// UpdateTenderStatus update Tender Status from tenders DB
func (s *Storage) UpdateTenderStatus(tenderId string, status string) (err error) {
	const op = "storage.UpdateTenderStatus"

	query := "UPDATE tenders SET status = $1 WHERE id = $2"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(status, tenderId)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

// UpdateTender update Tender  in tenders DB
func (s *Storage) UpdateTender(tender models.Tender) (err error) {
	const op = "storage.UpdateTender"

	query := "UPDATE tenders  SET name = $1, description = $2, serviceType = $3, version = $4 WHERE id = $5"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	_, err = stmt.Exec(tender.Name, tender.Description, tender.ServiceType, tender.Version, tender.Id)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

// GetCompanyID return organizationId from tenders DB
func (s *Storage) GetCompanyID(id string) (organizationId string, err error) {
	const op = "storage.GetCompanyID"

	query := "SELECT organizationId FROM tenders WHERE id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}

	err = stmt.QueryRow(id).Scan(&organizationId)
	if err != nil {
		return "", nil
	}

	return organizationId, nil
}

// GetTender returns one tender by id
func (s *Storage) GetTender(tender *models.Tender, tenderId string) error {
	const op = "storage.GetTender"

	var query = "SELECT id, name, description, serviceType, status, version, createdAt FROM tenders where id=$1 "

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(tenderId).Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.Version, &tender.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}

// GetFullTender returns one tender(name, description, serviceType, version) by id
func (s *Storage) GetFullTender(tender *models.Tender, tenderId string) error {
	const op = "storage.GetTender"

	var query = "SELECT * FROM tenders where id=$1 "

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s. Error preparing statement: %v", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(tenderId).Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatorUsername, &tender.Version, &tender.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s. Error executing query: %v", op, err)
	}

	return nil
}
