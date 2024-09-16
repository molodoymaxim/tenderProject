package response

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type response struct {
	Message string `json:"Message,omitempty"`
	Reason  string `json:"error,omitempty"`
}

func AnswerError(w http.ResponseWriter, r *http.Request, op string, httpCode int, err error) {
	reqID := middleware.GetReqID(r.Context())
	logMSG := fmt.Sprintf("op - %s. unsuccessful request. err = %s. requestID=%s", op, err, reqID)
	log.Printf(logMSG)
	msgErr := fmt.Sprintf("error in  %s. Error %s", op, err)
	response := response{"", msgErr}
	render.Status(r, httpCode)
	render.JSON(w, r, response)
}

func AnswerSuccess(w http.ResponseWriter, r *http.Request, httpCode int, msg string) {
	reqID := middleware.GetReqID(r.Context())
	logMSG := fmt.Sprintf("successful request. request msg = %s requestID=%", msg, reqID)
	log.Printf(logMSG)
	response := response{msg, ""}
	render.Status(r, httpCode)
	render.JSON(w, r, response)
}
