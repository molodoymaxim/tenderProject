package typecheck

import (
	"tenderProject/backend/internal/lib/models"
)

func IsTenderServiceTypeIncorrect(serviceType string) bool {
	for _, enm := range models.TenderServiceTypeEnum {
		if enm == serviceType {
			return false
		}
	}
	return true
}

func IsTenderStatusIncorrect(status string) bool {
	for _, enm := range models.TenderStatusEnum {
		if enm == status {
			return false
		}
	}
	return true
}

func IsAuthorTypeEnumIncorrect(status string) bool {
	for _, enm := range models.AuthorTypeEnum {
		if enm == status {
			return false
		}
	}
	return true
}

func IsBidsStatusEmumIncorrect(status string) bool {
	for _, enm := range models.BidsStatusEmum {
		if enm == status {
			return false
		}
	}
	return true
}
