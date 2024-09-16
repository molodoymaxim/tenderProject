package tenderbids

import (
	"errors"
	"fmt"
	"net/http"

	"tenderProject/backend/internal/handlers/get"
	"tenderProject/backend/internal/lib/api/limitandoffsetcheck"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type tenderBids interface {
	GetBids(bids *[]models.Bid, limit, offset int, username, companyID, tenderID string, bidSearchType int) error
	get.ServerGet
}

func TenderBidsH(server tenderBids) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "internal.get.TenderBidsH"
		tenderId := chi.URLParam(request, "tenderId")

		reqQuery := request.URL.Query()
		limit := reqQuery.Get("limit")
		offset := reqQuery.Get("offset")
		username := reqQuery.Get("username")

		limitInt, offsetInt, err := limitandoffsetcheck.LimitAndOffsetCheck(limit, offset)
		if err != nil {
			err = errors.New("connot convert limit or offset to int")
			response.AnswerError(writer, request, op, http.StatusBadRequest, err)
			return
		}

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

		bids := make([]models.Bid, 0, limitInt)

		err = server.GetBids(&bids, limitInt, offsetInt, "", "", tenderId, get.BidSearchTenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot recive info from DB %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, bids)
	}
}
