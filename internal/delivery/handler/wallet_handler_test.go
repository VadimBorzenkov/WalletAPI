package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VadimBorzenkov/WalletAPI/internal/service/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestHandleTransaction проверяет обработчик HandleTransaction для различных сценариев транзакций.
func TestHandleTransaction(t *testing.T) {
	// Определяем тестовые случаи для метода HandleTransaction
	tests := []struct {
		name         string                                                // Название теста
		requestBody  TransactionRequest                                    // Тело запроса транзакции
		mockService  func(ctrl *gomock.Controller) *mock.MockWalletService // Мок сервис для тестирования
		expectedCode int                                                   // Ожидаемый HTTP-код ответа
	}{
		{
			// Успешный сценарий депозита
			name: "Deposit Success",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "DEPOSIT",
				Amount:        100.0,
			},
			// Настраиваем mock для успешного вызова Deposit
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Deposit("wallet-123", 100.0).Return(nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			// Успешный сценарий вывода средств
			name: "Withdraw Success",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "WITHDRAW",
				Amount:        50.0,
			},
			// Настраиваем mock для успешного вызова Withdraw
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Withdraw("wallet-123", 50.0).Return(nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			// Некорректный тип операции
			name: "Invalid Operation Type",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "TRANSFER",
				Amount:        50.0,
			},
			// Ожидаем, что сервис не вызовет методы, так как операция некорректна
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				return mock.NewMockWalletService(ctrl)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			// Ошибка из-за недостатка средств на счете
			name: "Withdraw Insufficient Funds",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "WITHDRAW",
				Amount:        200.0,
			},
			// Настраиваем mock для вызова Withdraw, возвращающего ошибку
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Withdraw("wallet-123", 200.0).Return(fmt.Errorf("insufficient funds"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	// Выполняем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем новое приложение Fiber и настраиваем обработчик и mock сервис
			app := fiber.New()
			mockService := tt.mockService(ctrl)
			logger := logrus.New()
			apiHandler := NewApiWalletHandler(mockService, logger)
			app.Post("/api/v1/transaction", apiHandler.HandleTransaction)

			// Создаем запрос с телом запроса из тестового случая
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/transaction", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Выполняем запрос и проверяем код ответа
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

// TestHandleBalance проверяет обработчик HandleBalance для разных случаев получения баланса.
func TestHandleBalance(t *testing.T) {
	// Определяем тестовые случаи для метода HandleBalance
	tests := []struct {
		name         string                                                // Название теста
		walletID     string                                                // ID кошелька для запроса баланса
		mockService  func(ctrl *gomock.Controller) *mock.MockWalletService // Mock сервис для тестирования
		expectedCode int                                                   // Ожидаемый HTTP-код ответа
		expectedBody string                                                // Ожидаемое тело ответа
	}{
		{
			// Успешный сценарий получения баланса
			name:     "Balance Success",
			walletID: "wallet-123",
			// Настраиваем mock для успешного вызова GetBalance
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().GetBalance("wallet-123").Return(100.0, nil)
				return s
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"balance":100}`,
		},
		{
			// Ошибка при получении баланса
			name:     "Balance Error",
			walletID: "wallet-123",
			// Настраиваем mock для вызова GetBalance, возвращающего ошибку
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().GetBalance("wallet-123").Return(0.0, fmt.Errorf("could not retrieve balance"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"could not retrieve balance"}`,
		},
	}

	// Выполняем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем приложение Fiber и настраиваем mock сервис
			app := fiber.New()
			mockService := tt.mockService(ctrl)
			logger := logrus.New()
			apiHandler := NewApiWalletHandler(mockService, logger)
			app.Get("/api/v1/wallet/:walletID", apiHandler.HandleBalance)

			// Создаем GET запрос для проверки баланса
			req := httptest.NewRequest(http.MethodGet, "/api/v1/wallet/"+tt.walletID, nil)
			resp, _ := app.Test(req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Декодируем тело ответа для проверки содержимого
			var respBody map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&respBody)

			if tt.expectedCode == http.StatusOK {
				// Проверяем баланс, если запрос успешен
				assert.Equal(t, respBody["balance"], float64(100.0))
			} else {
				// Проверяем наличие ошибки в ответе, если запрос завершился ошибкой
				assert.Contains(t, respBody["error"], "could not retrieve balance")
			}
		})
	}
}
