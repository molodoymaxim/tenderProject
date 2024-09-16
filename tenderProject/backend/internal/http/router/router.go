package router

import (
	"tenderProject/backend/internal/handlers/get/bidstatus"
	"tenderProject/backend/internal/handlers/get/mybids"
	"tenderProject/backend/internal/handlers/get/mytenders"
	"tenderProject/backend/internal/handlers/get/ping"
	"tenderProject/backend/internal/handlers/get/tenderbids"
	"tenderProject/backend/internal/handlers/get/tenders"
	"tenderProject/backend/internal/handlers/get/tenderstatus"
	"tenderProject/backend/internal/handlers/patch/editbid"
	"tenderProject/backend/internal/handlers/patch/edittender"
	"tenderProject/backend/internal/handlers/post/newbid"
	"tenderProject/backend/internal/handlers/post/newtender"
	"tenderProject/backend/internal/handlers/put/changebidstatus"
	"tenderProject/backend/internal/handlers/put/changetenderstatus"
	"tenderProject/backend/internal/handlers/put/rollbackbid"
	"tenderProject/backend/internal/handlers/put/rollbacktender"
	logmiddleware "tenderProject/backend/internal/middleware/log"
	"tenderProject/backend/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageData *storage.Storage) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(logmiddleware.LoggingMiddleware)

	router.Get("/api/ping", ping.Ping)
	router.Get("/api/tenders", tenders.GetTenderH(storageData))
	router.Get("/api/tenders/my", mytenders.GetMyTender(storageData))
	router.Get("/api/tenders/{tenderId}/status", tenderstatus.TenderStatus(storageData))
	router.Get("/api/bids/my", mybids.GetMyBids(storageData))
	router.Get("/api/bids/{tenderId}/list", tenderbids.TenderBidsH(storageData))
	router.Get("/api/bids/{bidId}/status", bidstatus.BidStatus(storageData))

	router.Put("/api/tenders/{tenderId}/status", changetenderstatus.ChangeTenderStatus(storageData))
	router.Put("/api/tenders/{tenderId}/rollback/{version}", rollbacktender.RollbackH(storageData))
	router.Put("/api/bids/{bidId}/status", changebidstatus.BidStatusChange(storageData))
	router.Put("/api/bids/{bidId}/rollback/{version}", rollbackbid.RollbackH(storageData))

	router.Post("/api/tenders/new", newtender.NewTenderH(storageData))
	router.Post("/api/bids/new", newbid.NewBidH(storageData))

	router.Patch("/api/tenders/{tenderId}/edit", edittender.EditTenderH(storageData))
	router.Patch("/api/bids/{bidId}/edit", editbid.EditBidH(storageData))

	router.NotFoundHandler()
	return router
}
