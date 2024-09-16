package rollbacktender

import (
	"fmt"
	"net/http"
	"tenderProject/backend/internal/handlers"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/api/versionvalidation"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type rollback interface {
	GetFullTender(oldTender *models.Tender, tenderId string) error
	SaveOldTender(oldTender models.Tender) error
	UpdateTender(tender models.Tender) (err error)
	GetOldTender(oldTender *models.Tender, version int) error
	handlers.DataServerChecks
}

func RollbackH(server rollback) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "internal.put.RollbackH"

		tenderId := chi.URLParam(request, "tenderId")
		versionRoll := chi.URLParam(request, "version")
		username := request.URL.Query().Get("username")

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

		versionRollInt, ok := versionvalidation.ValidateVersion(writer, request, op, versionRoll)
		if !ok {
			return
		}

		tender := models.Tender{}
		err = server.GetFullTender(&tender, tenderId)
		if err != nil {
			msgErr := fmt.Errorf("failed to retrieve data from the database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		if tender.Version < versionRollInt {
			msgErr := fmt.Errorf("requested version is too large")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		err = server.SaveOldTender(tender)
		if err != nil {
			msgErr := fmt.Errorf("failed to save data to database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		// реализовать получение старых данных
		tender.Version += 1

		err = server.GetOldTender(&tender, versionRollInt)
		if err != nil {
			msgErr := fmt.Errorf("cannot Get Old Tender Tender %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if tender.Id == "" {
			msgErr := fmt.Errorf("cannot find Old Tender Tender %w", err)
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}

		err = server.UpdateTender(tender)
		if err != nil {
			msgErr := fmt.Errorf("cannot update Tender %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		tender.CreatorUsername = ""
		tender.OrganizationId = ""

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, tender)
	}
}
