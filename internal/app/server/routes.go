package server

func (s *server) RegisterHandlers() {
	s.r.GET("/api/bids/my", s.GetUserBids)
	s.r.POST("/api/bids/new", s.CreateBid)
	s.r.PATCH("/api/bids/:bidId/edit", s.EditBid)
	s.r.PUT("/api/bids/:bidId/feedback", s.SubmitBidFeedback)
	s.r.PUT("/api/bids/:bidId/rollback/:version", s.RollbackBid)
	s.r.GET("/api/bids/:bidId/status", s.GetBidStatus)
	s.r.PUT("/api/bids/:bidId/status", s.UpdateBidStatus)
	s.r.PUT("/api/bids/:bidId/submit_decision", s.SubmitBidDecision)
	s.r.GET("/api/bids/:tenderId/list", s.GetBidsForTender)
	s.r.GET("/api/bids/:tenderId/reviews", s.GetBidReviews)
	s.r.GET("/api/ping", s.CheckServer)
	s.r.GET("/api/tenders", s.GetTenders)
	s.r.GET("/api/tenders/my", s.GetUserTenders)
	s.r.POST("/api/tenders/new", s.CreateTender)
	s.r.PATCH("/api/tenders/:tenderId/edit", s.EditTender)
	s.r.PUT("/api/tenders/:tenderId/rollback/:version", s.RollbackTender)
	s.r.GET("/api/tenders/:tenderId/status", s.GetTenderStatus)
	s.r.PUT("/api/tenders/:tenderId/status", s.UpdateTenderStatus)
}
