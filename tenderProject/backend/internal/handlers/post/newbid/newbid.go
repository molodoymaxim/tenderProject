package newbid

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"tenderProject/backend/internal/handlers"
	"tenderProject/backend/internal/lib/api/response"
	"tenderProject/backend/internal/lib/api/typecheck"
	"tenderProject/backend/internal/lib/models"

	"github.com/go-chi/render"
)

type createBid interface {
	CreateNewBid(bid models.Bid) (string, error)
	handlers.DataServerChecks
}

func NewBidH(server createBid) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.post.NewBidH"
		var bid models.Bid

		err := render.DecodeJSON(request.Body, &bid)

		if errors.Is(err, io.EOF) {
			msgErr := fmt.Errorf("handle an error if we receive a request with an empty body.")
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}
		if err != nil {
			msgErr := fmt.Errorf("failed to deserialize the request. : %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		err = bid.Validate()
		if err != nil {
			msgErr := fmt.Errorf("data validation failed. : %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}

		if typecheck.IsAuthorTypeEnumIncorrect(bid.AuthorType) {
			msgErr := fmt.Errorf("Incorrect author type format received")
			response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
			return
		}

		switch bid.AuthorType {

		case models.AuthorTypeEnum[0]: // models.AuthorTypeEnum[0] is  User
			сhecking, err := server.CheckUserExists(bid.AuthorId)
			if err != nil {
				msgErr := fmt.Errorf("cannot Check User Exists :%s", err)
				response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
				return
			}
			if !сhecking {
				msgErr := fmt.Errorf("The user does not exist or is incorrect.")
				response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
				return
			}
		case models.AuthorTypeEnum[1]: // models.AuthorTypeEnum[1] is Organization
			сhecking, err := server.CheckOrganizationExists(bid.AuthorId)
			if err != nil {
				msgErr := fmt.Errorf("Cannot check organization exists :%s", err)
				response.AnswerError(writer, request, op, http.StatusBadRequest, msgErr)
				return
			}

			if !сhecking {
				msgErr := fmt.Errorf("The organization does not exist or is incorrect.")
				response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
				return
			}
		}

		сhecking, err := server.CheckTenderExists(bid.TenderId)
		if err != nil {
			msgErr := fmt.Errorf("cannot check tender Exists %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		if !сhecking {
			msgErr := fmt.Errorf("tender number not found")
			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
			return
		}

		bid.CreatedAt = time.Now()
		bid.Version = 1
		bid.Status = "Created"
		id, err := server.CreateNewBid(bid)
		if err != nil {
			msgErr := fmt.Errorf("failed to create a new tender %w", err)
			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
			return
		}
		bid.Id = id
		render.Status(request, http.StatusOK)
		render.JSON(writer, request, bid)

	}
}
