package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

func (s *server) GetUserBids(ctx echo.Context) error {
	var err error

	defaultLimit := int32(10)
	defaultOffset := int32(0)
	defaultUsername := ""

	params := repos.GetUserBidsParams{
		Limit:    &defaultLimit,
		Offset:   &defaultOffset,
		Username: &defaultUsername,
	}

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

	bids, err := s.bidHandler.GetUserBids(context.TODO(), params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bids)
}

func (s *server) CreateBid(ctx echo.Context) error {
	var err error
	var requestBody CreateBidJSONRequestBody

	if err = ctx.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request format: %s", err))
	}

	params := repos.CreateBidParams{
		Name:            &requestBody.Name,
		Description:     &requestBody.Description,
		Status:          &requestBody.Status,
		TenderID:        &requestBody.TenderId,
		OrganizationID:  &requestBody.OrganizationId,
		CreatorUsername: &requestBody.CreatorUsername,
	}

	// здесь надо создавать еще контекст с таймаутом, например
	// TODO: добавить контекст

	// Валидация здесь + потом создание записи в бд, если все гуд
	// Return структура бида + err
	bid, err := s.bidHandler.CreateBid(context.TODO(), params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
}
func (s *server) EditBid(ctx echo.Context) error {
	var err error

	var bidId repos.BidId
	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	var username repos.Username
	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	var requestBody EditBidJSONRequestBody

	if err := ctx.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
	}

	params := repos.EditBidParams{
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}

	bid, err := s.bidHandler.EditBid(context.TODO(), bidId, username, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
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

	// здесь булет идти валидация запроса.
	// проверяем, является ли username автором бида с bidId,
	// либо он состоит в орагнизации, которая является автором бида
	// возвращаем бид (? зачем?)
	bid, err := s.bidHandler.SubmitBidFeedback(context.TODO(), bidId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
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

	bid, err := s.bidHandler.RollbackBid(context.TODO(), bidId, version, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
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

	status, err := s.bidHandler.GetBidStatus(context.TODO(), bidId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, status)
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

	bid, err := s.bidHandler.UpdateBidStatus(context.TODO(), bidId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
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

	bid, err := s.bidHandler.SubmitBidDecision(context.TODO(), bidId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bid)
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

	bids, err := s.bidHandler.GetBidsForTender(context.TODO(), tenderId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, bids)
}

func (s *server) GetBidReviews(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var params repos.GetBidReviewsParams

	err = runtime.BindQueryParameter("form", true, true, "authorUsername", ctx.QueryParams(), &params.AuthorUsername)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorUsername: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, true, "requesterUsername", ctx.QueryParams(), &params.RequesterUsername)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter requesterUsername: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	revs, err := s.bidHandler.GetBidReviews(context.TODO(), tenderId, params)
	if err != nil {
		httpStatus, errResp := getStatusByError(err)
		return ctx.JSON(httpStatus, errResp)
	}
	return ctx.JSON(http.StatusOK, revs)
}
