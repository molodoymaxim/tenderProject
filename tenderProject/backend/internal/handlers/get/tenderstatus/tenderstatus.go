package tenderstatus

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tenderProject/backend/internal/handlers/get"
	"tenderProject/backend/internal/lib/api/response"
)

func TenderStatus(server get.ServerGet) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "internal.post.tenderstatus"
		tenderId := chi.URLParam(request, "tenderId")

		username := request.URL.Query().Get("username")

		cheking, err := server.CheckUserExists(username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check User Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !cheking {
			msgErr := fmt.Errorf("the user does not exist or is incorrect.")
			response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
			return
		}

		cheking, err = server.CheckTenderExists(tenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot check tender Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !cheking {
			msgErr := fmt.Errorf("tender number not found")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		organizationId, err := server.GetCompanyID(tenderId)
		if err != nil {
			msgErr := fmt.Errorf("company not found.")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		cheking, err = server.CheckOrganizationIdAndUserIDExists(organizationId, username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check and Organization User Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !cheking {
			msgErr := fmt.Errorf("the link between the employee and the organization could not be verified")
			response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
			return
		}

		status, err := server.GetTenderStatus(tenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot recive info from DB %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if status == "" {
			msgErr := fmt.Errorf("Tender not found.")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, status)
	}
}
