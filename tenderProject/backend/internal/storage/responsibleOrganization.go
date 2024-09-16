package storage

import (
	"fmt"
	"log"
)

func (s *Storage) CreateRelation() {
	createRelation := `
  ALTER TABLE organization_responsible
  ADD CONSTRAINT unique_organization_responsible UNIQUE (organization_id);
);
`
	_, err := s.db.Exec(createRelation)
	if err != nil {
		msgErr := fmt.Errorf("Error create relation:", err)
		log.Println(msgErr)
	}
	log.Println("managed to create an relationship")
}
