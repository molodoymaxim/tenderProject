package get

import (
	"tenderProject/backend/internal/handlers"
	"tenderProject/backend/internal/lib/models"
)

const ServiceTypeIsEmpty = 0
const ServiceTypeNotEmpty = 1
const UsernameSearch = 2

const BidSearchByUser = 0
const BidSearchByUserAncCompanyID = 1
const BidSearchTenderId = 2

type ServerGet interface {
	GetTenders(tenders *[]models.Tender, limit, offset int, searchInfo string, serchingType int) error
	GetTenderStatus(id string) (status string, err error)
	handlers.DataServerChecks
}
