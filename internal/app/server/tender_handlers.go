package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

func (s *server) GetTenders(ctx echo.Context) error {
	var err error

	var params repos.GetTendersParams

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Construction, Delivery, Manufacture
	err = runtime.BindQueryParameter("form", true, false, "service_type", ctx.QueryParams(), &params.ServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter service_type: %s", err))
	}

	// валидируем запрос, делаем запросик в бд, получаем список
	// TODO: изменить возвращаемое значение
	err = s.tenderHandler.GetTenders(context.TODO(), params)
	return err
}

func (s *server) GetUserTenders(ctx echo.Context) error {
	var err error
	var params repos.GetUserTendersParams

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

	// Получаем списко тендеров, но перед этим
	// валидируем данные, проверяем доступ юзера к тендерам
	err = s.tenderHandler.GetUserTenders(context.TODO(), params)
	return err
}

func (s *server) CreateTender(ctx echo.Context) error {
	var err error

	var requestBody CreateTenderJSONRequestBody

	if err = ctx.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request format: %s", err))
	}

	params := repos.CreateTenderParams{
		Name:            &requestBody.Name,
		Description:     &requestBody.Description,
		ServiceType:     &requestBody.ServiceType,
		Status:          &requestBody.Status,
		OrganizationID:  &requestBody.OrganizationId,
		CreatorUsername: &requestBody.CreatorUsername,
	}

	// Валидация здесь + потом создание записи в бд, если все гуд
	// Return структура бида + err
	err = s.tenderHandler.CreateTender(context.TODO(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create bid: %s", err))
	}

	err = s.tenderHandler.CreateTender(context.TODO(), repos.CreateTenderParams{})
	return err
}

func (s *server) EditTender(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var username repos.Username
	err = runtime.BindQueryParameter("form", true, false, "username", ctx.QueryParams(), &username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	var requestBody EditTenderJSONRequestBody

	if err = ctx.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request format: %s", err))
	}

	params := repos.EditTenderParams{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		ServiceType: requestBody.ServiceType,
	}

	err = s.tenderHandler.EditTender(context.TODO(), tenderId, username, params)
	return err
}

func (s *server) RollbackTender(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var version int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "version", runtime.ParamLocationPath, ctx.Param("version"), &version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter version: %s", err))
	}

	var params repos.RollbackTenderParams

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.tenderHandler.RollbackTender(context.TODO(), tenderId, version, params)
	return err
}

func (s *server) GetTenderStatus(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var params repos.GetTenderStatusParams

	err = runtime.BindQueryParameter("form", true, false, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.tenderHandler.GetTenderStatus(context.TODO(), tenderId, params)
	return err
}

func (s *server) UpdateTenderStatus(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var params repos.UpdateTenderStatusParams

	err = runtime.BindQueryParameter("form", true, true, "status", ctx.QueryParams(), &params.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter status: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.tenderHandler.UpdateTenderStatus(context.TODO(), tenderId, params)
	return err
}
