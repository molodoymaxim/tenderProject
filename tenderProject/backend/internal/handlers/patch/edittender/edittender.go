package edittender

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"tenderProject/backend/internal/handlers"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/api/typecheck"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type editTender interface {
	GetFullTender(oldTender *models.Tender, tenderId string) error
	SaveOldTender(oldTender models.Tender) error
	UpdateTender(tender models.Tender) (err error)
	handlers.DataServerChecks
}

func EditTenderH(server editTender) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "internal.patch.EditTenderH"

		tenderId := chi.URLParam(request, "tenderId")
		username := request.URL.Query().Get("username")

		var newTender models.Tender

		err := render.DecodeJSON(request.Body, &newTender)
		if errors.Is(err, io.EOF) {
			msgErr := fmt.Errorf("handle an error if we receive a request with an empty body")
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}
		if err != nil {
			msgErr := fmt.Errorf("failed to deserialize the request. : %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		err = newTender.Validate()
		if err != nil {
			msgErr := fmt.Errorf("data validation failed. : %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
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

		tender := models.Tender{}
		err = server.GetFullTender(&tender, tenderId)
		if err != nil {
			msgErr := fmt.Errorf("failed to retrieve data from the database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		err = server.SaveOldTender(tender)
		if err != nil {
			msgErr := fmt.Errorf("failed to save data to database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		if newTender.Name != "" {
			tender.Name = newTender.Name
		}
		if newTender.Description != "" {
			tender.Description = newTender.Description
		}
		if newTender.ServiceType != "" {
			if typecheck.IsTenderServiceTypeIncorrect(newTender.ServiceType) {
				msgErr := fmt.Errorf("incorrectly specified service type")
				response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
				return
			}
			tender.ServiceType = newTender.ServiceType
		}
		tender.Version += 1

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
