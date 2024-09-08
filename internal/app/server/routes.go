package server

import (
	"github.com/labstack/echo/v4"
)

// RegisterHandlers регистрирует маршруты и связывает их с обработчиками
func (s *server) RegisterHandlers(e *echo.Echo) {
	s.r.GET("/bids/my", s.GetUserBids)
	s.r.POST("/bids/new", s.CreateBid)
	s.r.PATCH("/bids/:bidId/edit", s.EditBid)
	s.r.PUT("/bids/:bidId/feedback", s.SubmitBidFeedback)
	s.r.PUT("/bids/:bidId/rollback/:version", s.RollbackBid)
	s.r.GET("/bids/:bidId/status", s.GetBidStatus)
	s.r.PUT("/bids/:bidId/status", s.UpdateBidStatus)
	s.r.PUT("/bids/:bidId/submit_decision", s.SubmitBidDecision)
	s.r.GET("/bids/:tenderId/list", s.GetBidsForTender)
	s.r.GET("/bids/:tenderId/reviews", s.GetBidReviews)
	s.r.GET("/ping", s.CheckServer)
	s.r.GET("/tenders", s.GetTenders)
	s.r.GET("/tenders/my", s.GetUserTenders)
	s.r.POST("/tenders/new", s.CreateTender)
	s.r.PATCH("/tenders/:tenderId/edit", s.EditTender)
	s.r.PUT("/tenders/:tenderId/rollback/:version", s.RollbackTender)
	s.r.GET("/tenders/:tenderId/status", s.GetTenderStatus)
	s.r.PUT("/tenders/:tenderId/status", s.UpdateTenderStatus)
}
