package handler

import (
	"github.com/VadimBorzenkov/WalletAPI/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WalletHandler interface {
	HandleBalance(c *fiber.Ctx) error
	HandleTransaction(c *fiber.Ctx) error
}

type ApiWalletHandler struct {
	walletService service.WalletService
	logger        *logrus.Logger
}

func NewApiWalletHandler(walletService service.WalletService, logger *logrus.Logger) *ApiWalletHandler {
	return &ApiWalletHandler{
		walletService: walletService,
		logger:        logger,
	}
}

// HandleBalance обрабатывает запрос на получение баланса кошелька
func (h *ApiWalletHandler) HandleBalance(c *fiber.Ctx) error {
	walletID := c.Params("walletID")
	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "walletID is required"})
	}

	balance, err := h.walletService.GetBalance(walletID)
	if err != nil {
		h.logger.Errorf("Failed to get balance for wallet %s: %v", walletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve balance"})
	}

	return c.JSON(fiber.Map{"balance": balance})
}

type TransactionRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"` // "DEPOSIT" or "WITHDRAW"
	Amount        float64 `json:"amount"`
}

// HandleTransaction обрабатывает запрос на выполнение операции с кошельком
func (h *ApiWalletHandler) HandleTransaction(c *fiber.Ctx) error {
	var req TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "amount must be positive"})
	}

	var err error
	switch req.OperationType {
	case "DEPOSIT":
		err = h.walletService.Deposit(req.WalletID, req.Amount)
	case "WITHDRAW":
		err = h.walletService.Withdraw(req.WalletID, req.Amount)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid operation type"})
	}

	if err != nil {
		h.logger.Errorf("Failed to process transaction for wallet %s: %v", req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "transaction successful"})
}
