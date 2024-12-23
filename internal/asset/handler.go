package asset

import (
	"fmt"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/asset/request"
	"github.com/safayildirim/asset-management-service/internal/common"
	walletpkg "github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"net/http"
	"reflect"
	"strings"
)

// Initialize a schema decoder for parsing query parameters
var decoder = schema.NewDecoder()

func init() {
	// Configure the schema decoder to ignore unknown query keys
	decoder.IgnoreUnknownKeys(false)
	// Register a custom converter to handle string slices (comma-separated values)
	decoder.RegisterConverter([]string{}, func(value string) reflect.Value {
		return reflect.ValueOf(strings.Split(value, ","))
	})
}

type Handler struct {
	assetService Service
}

// NewHandler initializes a new Handler instance with the provided asset service
func NewHandler(assetService Service) *Handler {
	return &Handler{assetService: assetService}
}

// RegisterRoutes registers the asset-related API routes with the provided Echo router group
func (h Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/assets", h.CreateAsset)
	e.GET("/assets", h.GetAssets)
	e.POST("/assets/deposit", h.Deposit)
	e.POST("/assets/withdraw", h.Withdraw)
}

// CreateAsset handles requests to create a new asset
func (h Handler) CreateAsset(ctx echo.Context) error {
	var req request.CreateAssetRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	asset, err := h.assetService.CreateAsset(ctx.Request().Context(), nil, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateAsset):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, common.Response{Data: asset})
}

// GetAssets handles requests to fetch a list of assets
func (h Handler) GetAssets(ctx echo.Context) error {
	var req request.GetAssetsParams
	params := ctx.QueryParams()

	err := decoder.Decode(&req, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	assets, err := h.assetService.GetAssets(ctx.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, common.Response{Data: assets})
}

// Deposit handles requests to deposit an asset into a wallet
func (h Handler) Deposit(ctx echo.Context) error {
	var req request.CreateDepositRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	asset, err := h.assetService.Deposit(ctx.Request().Context(), nil, &req)
	if err != nil {
		switch {
		case errors.Is(err, walletpkg.ErrWalletNotFound):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, common.Response{Data: asset})
}

// Withdraw handles requests to withdraw an asset from a wallet
func (h Handler) Withdraw(ctx echo.Context) error {
	var req request.CreateWithdrawRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	asset, err := h.assetService.Withdraw(ctx.Request().Context(), nil, &req)
	if err != nil {
		switch {
		case errors.Is(err, walletpkg.ErrWalletNotFound):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, common.Response{Data: asset})
}
