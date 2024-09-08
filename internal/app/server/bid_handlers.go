package server

import (
	"fmt"
	"net/http"

	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

func (s *server) GetUserBids(ctx echo.Context) error {
	var err error
	var params repos.GetUserBidsParams

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.GetUserBids(ctx, params)
	return err
}

func (s *server) CreateBid(ctx echo.Context) error {
	var err error

	err = s.bidHandler.CreateBid(ctx, repos.CreateBidParams{})
	return err
}

func (s *server) EditBid(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var params repos.EditBidParams

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.EditBid(ctx, bidId, params)
	return err
}

func (s *server) SubmitBidFeedback(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var params repos.SubmitBidFeedbackParams

	err = runtime.BindQueryParameter("form", true, true, "bidFeedback", ctx.QueryParams(), &params.BidFeedback)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidFeedback: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.SubmitBidFeedback(ctx, bidId, params)
	return err
}

func (s *server) RollbackBid(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var version int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "version", runtime.ParamLocationPath, ctx.Param("version"), &version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter version: %s", err))
	}

	var params repos.RollbackBidParams
	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.RollbackBid(ctx, bidId, version, params)
	return err
}

func (s *server) GetBidStatus(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var params repos.GetBidStatusParams

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.GetBidStatus(ctx, bidId, params)
	return err
}

func (s *server) UpdateBidStatus(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var params repos.UpdateBidStatusParams

	err = runtime.BindQueryParameter("form", true, true, "status", ctx.QueryParams(), &params.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter status: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.UpdateBidStatus(ctx, bidId, params)
	return err
}

func (s *server) SubmitBidDecision(ctx echo.Context) error {
	var err error
	var bidId repos.BidId

	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var params repos.SubmitBidDecisionParams
	err = runtime.BindQueryParameter("form", true, true, "decision", ctx.QueryParams(), &params.Decision)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter decision: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.bidHandler.SubmitBidDecision(ctx, bidId, params)
	return err
}

func (s *server) GetBidsForTender(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var params repos.GetBidsForTenderParams

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	err = s.bidHandler.GetBidsForTender(ctx, tenderId, params)
	return err
}

// GetBidReviews converts echo context to params.
func (s *server) GetBidReviews(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tenderId" -------------
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params repos.GetBidReviewsParams
	// ------------- Required query parameter "authorUsername" -------------

	err = runtime.BindQueryParameter("form", true, true, "authorUsername", ctx.QueryParams(), &params.AuthorUsername)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorUsername: %s", err))
	}

	// ------------- Required query parameter "requesterUsername" -------------

	err = runtime.BindQueryParameter("form", true, true, "requesterUsername", ctx.QueryParams(), &params.RequesterUsername)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter requesterUsername: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	err = s.bidHandler.GetBidReviews(ctx, tenderId, params)
	return err
}
