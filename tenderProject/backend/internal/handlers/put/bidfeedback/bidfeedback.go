package bidfeedback

//
//import (
//	"tenderProject/backend/internal/lib/api/response"
//	"tenderProject/backend/internal/lib/models"
//	"fmt"
//	"github.com/go-chi/chi/v5"
//	"github.com/go-chi/render"
//	"net/http"
//)
//
//type createFeedBackI interface {
//	CheckBidExists(bidID string) (bool, error)
//	GetBid(bid *models.Bid, bidID string) (successfulRequest bool, err error)
//	GetUserIDUsername(username string) (string, error)
//}
//
//func CreateFeedBack(server createFeedBackI) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		const op = "backend.internal.handlers.put.CreateFeedBack"
//		bidId := chi.URLParam(request, "bidId")
//		username := request.URL.Query().Get("username")
//		bidFeedback := request.URL.Query().Get("bidFeedback")
//
//		if len([]rune(username)) > 100 {
//			msgErr := fmt.Errorf("username to big")
//			response.AnswerError(writer, request, op, http.StatusRequestEntityTooLarge, msgErr)
//			return
//
//		}
//		if len([]rune(bidFeedback)) > 1000 {
//			msgErr := fmt.Errorf("bidFeedback to big")
//			response.AnswerError(writer, request, op, http.StatusRequestEntityTooLarge, msgErr)
//			return
//
//		}
//
//		successfulRequest, err := server.CheckBidExists(bidId)
//		if err != nil {
//			msgErr := fmt.Errorf("failed to retrieve information from the database %w", err)
//			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//			return
//		}
//		if !successfulRequest {
//			msgErr := fmt.Errorf("couldn't find the bid")
//			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
//			return
//		}
//
//		bid := models.Bid{}
//		successfulRequest, err = server.GetBid(&bid, bidId)
//		if err != nil {
//			msgErr := fmt.Errorf("failed to retrieve information from the database %w", err)
//			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//			return
//		}
//		if !successfulRequest {
//			msgErr := fmt.Errorf("couldn't find the bid")
//			response.AnswerError(writer, request, op, http.StatusNotFound, msgErr)
//			return
//		}
//
//		var organizationId string
//		server.GetUserIDUsername(username)
//		switch bid.AuthorType {
//		case models.AuthorTypeEnum[0]: // user
//			cheking, err := server.CheckUserExists(username)
//			if err != nil {
//				msgErr := fmt.Errorf("cannot check user exists %w", err)
//				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//				return
//			}
//			if !cheking {
//				msgErr := fmt.Errorf("the user does not exist or is incorrect.")
//				response.AnswerError(writer, request, op, http.StatusUnauthorized, msgErr)
//				return
//			}
//			organizationId, err = server.GetCompanyIDbyUser(username)
//			if err != nil {
//				msgErr := fmt.Errorf("failed to retrieve organization data :%s", err)
//				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//				return
//			}
//			if bid.AuthorId != organizationId && bid.AuthorId != username {
//				msgErr := fmt.Errorf("this bid is not available to you")
//				response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
//				return
//			}
//		case models.AuthorTypeEnum[1]: // organization
//			сhecking, err := server.CheckOrganizationExists(username)
//			if err != nil {
//				msgErr := fmt.Errorf("cannot check organization exists %w", err)
//				response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//				return
//			}
//			if !сhecking {
//				msgErr := fmt.Errorf("this user cannot get information for this tender")
//				response.AnswerError(writer, request, op, http.StatusForbidden, msgErr)
//				return
//			}
//		}
//
//		err = server.UpdateBidStatus(bidId, status)
//		if err != nil {
//			msgErr := fmt.Errorf("cannot update bid status %w", err)
//			response.AnswerError(writer, request, op, http.StatusInternalServerError, msgErr)
//			return
//		}
//
//		bid.Status = status
//
//		render.Status(request, http.StatusOK)
//		render.JSON(writer, request, bid)
//	}
//}
