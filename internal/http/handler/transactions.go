package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
    transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
    return TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	input := new(binder.TransactionCreateRequest)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid input"))
	}

	// Get the fee from the pricing table based on event_id
	pricing, err := h.transactionService.GetPricingByEventID(input.EventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to get pricing"))
	}

	totalAmount := input.TicketQuantity * pricing.Fee

	newTransaction := entity.NewTransaction(input.EventID, input.UserID, input.TicketQuantity, totalAmount, input.IsPaid)

	transaction, err := h.transactionService.CreateTransaction(newTransaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Transaction created successfully", transaction))
}


func (h *TransactionHandler) FindTransactionByID(c echo.Context) error {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
    }

    transaction, err := h.transactionService.FindTransactionByID(id)
    if err != nil {
        return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Transaction not found"))
    }

    return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data transaksi", transaction))
}

func (h *TransactionHandler) FindAllTransactions(c echo.Context) error {
    transactions, err := h.transactionService.FindAllTransactions()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
    }

    return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data transaksi", transactions))
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
    }

    var input binder.TransactionUpdateRequest
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
    }

    transaction := &entity.Transaction{
        ID:             id,
        EventID:        input.EventID,
        UserID:         input.UserID,
        TicketQuantity: input.TicketQuantity,
        TotalAmount:    input.TotalAmount,
        IsPaid:         input.IsPaid,
    }

    updatedTransaction, err := h.transactionService.UpdateTransaction(transaction)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
    }

    return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mengupdate transaksi", updatedTransaction))
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
    }

    err = h.transactionService.DeleteTransaction(id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
    }

    return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Transaction deleted successfully", nil))
}