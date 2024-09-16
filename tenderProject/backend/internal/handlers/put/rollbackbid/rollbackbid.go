package rollbackbid

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"tenderProject/backend/internal/handlers/get"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/api/versionvalidation"
	"tenderProject/backend/internal/lib/models"
)

type rollbackI interface {
	GetBid(bid *models.Bid, bidID string) (successfulRequest bool, err error)
	GetCompanyIDbyUser(username string) (companyId string, err error)
	CheckBidExists(bidID string) (bool, error)
	SaveOldBid(bid models.Bid) error
	UpdateBid(bid models.Bid) error
	get.ServerGet
}

func RollbackH(server rollbackI) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "backend.internal.handlers.put.RollbackH"
		bidId := chi.URLParam(request, "bidId")
		username := request.URL.Query().Get("username")
		versionRoll := chi.URLParam(request, "version")

		versionRollInt, ok := versionvalidation.ValidateVersion(writer, request, op, versionRoll)
		if !ok {
			return
		}

		successfulRequest, err := server.CheckBidExists(bidId)
		if err != nil {
			msgErr := fmt.Errorf("failed to retrieve information from the database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !successfulRequest {
			msgErr := fmt.Errorf("couldn't find the bid")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		bid := models.Bid{}
		successfulRequest, err = server.GetBid(&bid, bidId)
		if err != nil {
			msgErr := fmt.Errorf("failed to retrieve information from the database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !successfulRequest {
			msgErr := fmt.Errorf("couldn't find the bid")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		var organizationId string
		switch bid.AuthorType {
		case models.AuthorTypeEnum[0]: // user
			cheking, err := server.CheckUserExists(username)
			if err != nil {
				msgErr := fmt.Errorf("cannot check user exists %w", err)
				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
				return
			}
			if !cheking {
				msgErr := fmt.Errorf("the user does not exist or is incorrect.")
				response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
				return
			}
			organizationId, err = server.GetCompanyIDbyUser(username)
			if err != nil {
				msgErr := fmt.Errorf("failed to retrieve organization data :%s", err)
				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
				return
			}
			if bid.AuthorId != organizationId && bid.AuthorId != username {
				msgErr := fmt.Errorf("this bid is not available to you")
				response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
				return
			}
		case models.AuthorTypeEnum[1]: // organization
			сhecking, err := server.CheckOrganizationExists(username)
			if err != nil {
				msgErr := fmt.Errorf("cannot check organization exists %w", err)
				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
				return
			}
			if !сhecking {
				msgErr := fmt.Errorf("this user cannot get information for this tender")
				response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
				return
			}
		}

		if bid.Version < versionRollInt {
			msgErr := fmt.Errorf("requested version is too large")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		err = server.SaveOldBid(bid)
		if err != nil {
			msgErr := fmt.Errorf("cannot save bid in bids_version database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		bid.Version += 1

		err = server.UpdateBid(bid)
		if err != nil {
			msgErr := fmt.Errorf("failed to update the bid in bid database %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		bid.TenderId = ""
		bid.Description = ""

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, bid)
	}
}
