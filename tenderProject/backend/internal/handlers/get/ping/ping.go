package ping

import (
	"net/http"
	"tenderProject/backend/internal/lib/api/response"
)

func Ping(writer http.ResponseWriter, request *http.Request) {
	response.AnswerSuccess(writer, request, http.StatusOK, "PING CALLED")
	return
}
