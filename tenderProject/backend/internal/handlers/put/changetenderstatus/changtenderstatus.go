package changetenderstatus

import (
	"fmt"
	"net/http"

	"tenderProject/backend/internal/handlers"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/api/typecheck"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UpdateTenderStatus interface {
	UpdateTenderStatus(tenderId string, status string) (err error)
	GetTender(tender *models.Tender, tenderId string) error
	handlers.DataServerChecks
}

// ChangeTenderStatus Handler  to change the status of the tender
func ChangeTenderStatus(server UpdateTenderStatus) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "internal.put.ChangeTenderStatus"
		tenderId := chi.URLParam(request, "tenderId")

		username := request.URL.Query().Get("username")
		status := request.URL.Query().Get("status")

		if typecheck.IsTenderStatusIncorrect(status) {
			msgErr := fmt.Errorf("incorrect status")
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}

		сhecking, err := server.CheckUserExists(username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check User Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !сhecking {
			msgErr := fmt.Errorf("the user does not exist or is incorrect.")
			response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
			return
		}

		сhecking, err = server.CheckTenderExists(tenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check and Organization User Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !сhecking {
			msgErr := fmt.Errorf("The tender number was not found")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		organizationId, err := server.GetCompanyID(tenderId)
		if err != nil {
			msgErr := fmt.Errorf("Tender not found.")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		сhecking, err = server.CheckOrganizationIdAndUserIDExists(organizationId, username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check and Organization User Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !сhecking {
			msgErr := fmt.Errorf("the link between the employee and the organization could not be verified")
			response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
			return
		}

		err = server.UpdateTenderStatus(tenderId, status)
		if err != nil {
			msgErr := fmt.Errorf("cannot update Tender Status %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		tender := models.Tender{}
		err = server.GetTender(&tender, tenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot recive info from DB %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, tender)
	}
}
