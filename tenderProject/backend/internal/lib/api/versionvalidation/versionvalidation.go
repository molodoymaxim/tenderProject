package versionvalidation

import (
	"fmt"
	"net/http"
	"strconv"
	"tenderProject/backend/internal/lib/api/response"
)

func ValidateVersion(writer http.ResponseWriter, request *http.Request, op, versionRoll string) (int, bool) {
	var versionRollInt int
	var err error
	if versionRoll == "" {
		msgErr := fmt.Errorf("version is empty")
		response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
		return 0, false
	} else {
		versionRollInt, err = strconv.Atoi(versionRoll)
		if err != nil {
			msgErr := fmt.Errorf("failed to convert the version to a number")
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return 0, false
		}
		if versionRollInt < 0 {
			msgErr := fmt.Errorf("version less than zero")
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return 0, false
		}
	}
	return versionRollInt, true
}
