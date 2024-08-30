package controller

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/service"
	"car-rental-application/pkg"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransactionController interface {
	UpdateTransactionStatus(c echo.Context) error
	GetTransactionByID(c echo.Context) error
	GetAllTransactions(c echo.Context) error
}

type transactionController struct {
	transactionService service.TransactionService
}

func (t *transactionController) UpdateTransactionStatus(c echo.Context) error {
	var payload struct {
		ExternalID  string `json:"external_id"`
		Status      string `json:"status"`
		PaymentType string `json:"payment_type"`
	}
	if err := c.Bind(&payload); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid payload: "+err.Error())
	}
	if payload.ExternalID == "" || payload.Status == "" {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "ExternalID or Status missing in request payload")
	}
	if payload.Status == "PAID" {
		err := t.transactionService.UpdateTransactionStatus(payload.ExternalID, "paid")
		if err != nil {
			return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to update transaction status: "+err.Error())
		}
		return pkg.RespondJSON(c, http.StatusOK, nil, "Transaction status updated to paid successfully")
	}
	return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Unexpected payment status")
}

func (t *transactionController) GetTransactionByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid transaction ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))
	userRole := claims["role"].(string)

	transaction, err := t.transactionService.GetTransactionById(uint(id))
	if err != nil {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "Transaction not found: "+err.Error())
	}

	if userRole != "admin" && transaction.UserID != userId {
		return pkg.RespondJSON(c, http.StatusForbidden, nil, "Access denied: You do not have permission to view this transaction")
	}

	return pkg.RespondJSON(c, http.StatusOK, transaction, "Transaction fetched successfully")
}

func (t *transactionController) GetAllTransactions(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))
	userRole := claims["role"].(string)

	transactions, err := t.transactionService.GetAllTransaction()
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to fetch transactions: "+err.Error())
	}

	if userRole != "admin" {
		var userTransactions []models.Transaction
		for _, transaction := range transactions {
			if transaction.UserID == userId {
				userTransactions = append(userTransactions, transaction)
			}
		}
		return pkg.RespondJSON(c, http.StatusOK, userTransactions, "Transactions fetched successfully")
	}

	return pkg.RespondJSON(c, http.StatusOK, transactions, "All transactions fetched successfully")
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}
