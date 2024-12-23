package transaction

import (
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/common"
	"github.com/safayildirim/asset-management-service/internal/transaction/request"
	walletpkg "github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"net/http"
	"reflect"
	"strings"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.RegisterConverter([]string{}, func(value string) reflect.Value {
		return reflect.ValueOf(strings.Split(value, ","))
	})
}

type Handler struct {
	transactionService Service
}

func NewHandler(transactionService Service) *Handler {
	return &Handler{transactionService: transactionService}
}

func (h Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/transactions/schedule", h.ScheduleTransaction)
	e.GET("/transactions", h.GetTransactions)
	e.DELETE("/transactions", h.DeleteTransaction)
}

func (h Handler) ScheduleTransaction(ctx echo.Context) error {
	var req request.ScheduleTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transaction, err := h.transactionService.ScheduleTransaction(ctx.Request().Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, walletpkg.ErrWalletNotFound):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, common.Response{Data: transaction})
}

func (h Handler) GetTransactions(ctx echo.Context) error {
	var req request.GetTransactionsParams
	params := ctx.QueryParams()

	err := decoder.Decode(&req, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = req.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	assets, err := h.transactionService.GetTransactions(ctx.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, common.Response{Data: assets})
}

func (h Handler) DeleteTransaction(ctx echo.Context) error {
	id, err := common.ParseIntFromString[uint](ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.transactionService.CancelTransaction(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}
