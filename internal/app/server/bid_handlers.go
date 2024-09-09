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

	// default params
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

	err = s.bidHandler.GetUserBids(ctx, params)
	return err
}

func (s *server) CreateBid(ctx echo.Context) error {
	var err error
	var params repos.CreateBidParams

	if err = ctx.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request format: %s", err))
	}

	// Валидация здесь + потом создание записи в бд, если все гуд
	// Return структура бида + err
	err = s.bidHandler.CreateBid(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create bid: %s", err))
	}

	// Возвращаем успешный ответ
	// TODO: возврат бида
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Bid successfully created",
	})
}
func (s *server) EditBid(ctx echo.Context) error {
	var err error

	// Извлечение параметра пути
	var bidId repos.BidId
	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	// Извлечение параметра запроса
	var username repos.Username
	err = runtime.BindQueryParameter("form", true, false, "username", ctx.QueryParams(), &username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	var requestBody EditBidJSONRequestBody

	if err := ctx.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
	}

	// Создание структуры для параметров
	params := repos.EditBidParams{
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}

	// Вызов метода для редактирования предложения
	bid, err := s.bidHandler.EditBid(ctx, bidId, params)
	if err != nil {
		return err
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

	// Параметры для изменения предложения
	var params repos.SubmitBidFeedbackParams

	// получаем фидбек
	err = runtime.BindQueryParameter("form", true, true, "bidFeedback", ctx.QueryParams(), &params.BidFeedback)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidFeedback: %s", err))
	}

	// получаем username
	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	// здесь булет идти валидация запроса.
	// проверяем, является ли username автором бида с bidId,
	// либо он состоит в орагнизации, которая является автором бида
	// возвращаем бид (? зачем?)
	err = s.bidHandler.SubmitBidFeedback(ctx, bidId, params)
	return err
}

// Откат версии бида с возвратом к старым значениям (значениям той версии - name и description)
func (s *server) RollbackBid(ctx echo.Context) error {
	var err error

	// получаем bidId из path
	var bidId repos.BidId
	err = runtime.BindStyledParameterWithLocation("simple", false, "bidId", runtime.ParamLocationPath, ctx.Param("bidId"), &bidId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter bidId: %s", err))
	}

	// получаем version из path
	var version int32
	err = runtime.BindStyledParameterWithLocation("simple", false, "version", runtime.ParamLocationPath, ctx.Param("version"), &version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter version: %s", err))
	}

	// получаем username из query
	var params repos.RollbackBidParams
	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	// валидиурем доступ юзера к версионированию
	// делаем откат и вовзращаем bid
	err = s.bidHandler.RollbackBid(ctx, bidId, version, params)
	return err
}

// Поулчаем статус предложения
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

	// Валидируем пользователя и возвращаем статус
	err = s.bidHandler.GetBidStatus(ctx, bidId, params)
	return err
}

// Апдейтим статус предложения
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

	// здесь будем валидировать пользователя и статус
	// далее изменяем статус
	err = s.bidHandler.UpdateBidStatus(ctx, bidId, params)
	return err
}

// Approved или Rejected
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

	// валидируем данные (доступ юзера + решение, а то вдруг передадим в решении, например, "нельзя(")
	// возвращаем бид и ошибку
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

	// Необязательные параметры
	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Список бидов, отсортированных по АЛФАВИТУ (какому алфавиту? по какому полю сортировать?)
	err = s.bidHandler.GetBidsForTender(ctx, tenderId, params)
	return err
}

// Достаем все прошлые reviews по тендеру
// TODO: отредачить, а то щас лень(
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

	err = s.bidHandler.GetBidReviews(ctx, tenderId, params)
	return err
}
