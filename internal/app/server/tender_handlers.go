package server

import (
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

	err = runtime.BindQueryParameter("form", true, false, "service_type", ctx.QueryParams(), &params.ServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter service_type: %s", err))
	}

	err = s.tenderHandler.GetTenders(ctx, params)
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

	err = s.tenderHandler.GetUserTenders(ctx, params)
	return err
}

func (s *server) CreateTender(ctx echo.Context) error {
	var err error

	err = s.tenderHandler.CreateTender(ctx, repos.CreateTenderParams{})
	return err
}

func (s *server) EditTender(ctx echo.Context) error {
	var err error
	var tenderId repos.TenderId

	err = runtime.BindStyledParameterWithLocation("simple", false, "tenderId", runtime.ParamLocationPath, ctx.Param("tenderId"), &tenderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tenderId: %s", err))
	}

	var params repos.EditTenderParams

	err = runtime.BindQueryParameter("form", true, true, "username", ctx.QueryParams(), &params.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter username: %s", err))
	}

	err = s.tenderHandler.EditTender(ctx, tenderId, params)
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

	err = s.tenderHandler.RollbackTender(ctx, tenderId, version, params)
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

	err = s.tenderHandler.GetTenderStatus(ctx, tenderId, params)
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

	err = s.tenderHandler.UpdateTenderStatus(ctx, tenderId, params)
	return err
}
