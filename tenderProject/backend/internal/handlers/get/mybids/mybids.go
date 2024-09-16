package mybids

import (
	"errors"
	"fmt"
	"net/http"

	"tenderProject/backend/internal/handlers/get"
	"tenderProject/backend/internal/lib/api/limitandoffsetcheck"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/render"
)

type getMyBidsI interface {
	GetCompanyIDbyUser(username string) (companyId string, err error)
	GetBids(bids *[]models.Bid, limit, offset int, username, companyID, tenderID string, bidSearchType int) error
	get.ServerGet
}

func GetMyBids(server getMyBidsI) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "backend.internal.handlers.get.GetMyBids"

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

		userChecking, err := server.CheckUserExists(username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check User Exists :%s", err)
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}
		companyChecking, err := server.CheckOrganizationExists(username)
		if err != nil {
			msgErr := fmt.Errorf("cannot Check User Exists :%s", err)
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}

		if !userChecking && !companyChecking {
			msgErr := fmt.Errorf("the user or company does not exist or is incorrect")
			response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
			return
		}

		searchType := get.BidSearchByUser
		var companyID string
		if userChecking {
			companyID, err = server.GetCompanyIDbyUser(username)
			if err != nil {
				msgErr := fmt.Errorf("failed to retrieve organization data :%s", err)
				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
				return
			}
			searchType = get.BidSearchByUserAncCompanyID
		}

		bids := make([]models.Bid, 0, limitInt)
		err = server.GetBids(&bids, limitInt, offsetInt, username, companyID, "", searchType)
		if err != nil {
			err = fmt.Errorf("error when receiving data from the server. %s", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, err)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, bids)
	}
}
