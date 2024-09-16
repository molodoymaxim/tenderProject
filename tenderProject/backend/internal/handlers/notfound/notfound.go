package notfound

import (
	"errors"
	"net/http"
	"tenderProject/backend/internal/lib/api/response"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	const op = "backend.handlers.notFound"
	err := errors.New("invalid URL request")
	response.AnswerError(w, r, op, http.StatusNotFound, err)
	return
}
