package handlers

type DataServerChecks interface {
	CheckUserExists(creatorUsername string) (bool, error)
	//CheckUserExistsByID(creatorUsername string) (bool, error)
	CheckOrganizationIdAndUserIDExists(organizationId, creatorUsername string) (bool, error)
	GetCompanyID(id string) (organizationId string, err error)
	CheckOrganizationExists(organizationid string) (bool, error)
	CheckTenderExists(tenderId string) (bool, error)
}
