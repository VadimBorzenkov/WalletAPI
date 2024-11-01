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

func TestHandleTransaction(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  TransactionRequest
		mockService  func(ctrl *gomock.Controller) *mock.MockWalletService
		expectedCode int
	}{
		{
			name: "Deposit Success",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "DEPOSIT",
				Amount:        100.0,
			},
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Deposit("wallet-123", 100.0).Return(nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Withdraw Success",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "WITHDRAW",
				Amount:        50.0,
			},
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Withdraw("wallet-123", 50.0).Return(nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Operation Type",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "TRANSFER",
				Amount:        50.0,
			},
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				return mock.NewMockWalletService(ctrl)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw Insufficient Funds",
			requestBody: TransactionRequest{
				WalletID:      "wallet-123",
				OperationType: "WITHDRAW",
				Amount:        200.0,
			},
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().Withdraw("wallet-123", 200.0).Return(fmt.Errorf("insufficient funds"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()
			mockService := tt.mockService(ctrl)
			logger := logrus.New()
			apiHandler := NewApiWalletHandler(mockService, logger)
			app.Post("/api/v1/transaction", apiHandler.HandleTransaction)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/transaction", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestHandleBalance(t *testing.T) {
	tests := []struct {
		name         string
		walletID     string
		mockService  func(ctrl *gomock.Controller) *mock.MockWalletService
		expectedCode int
		expectedBody string
	}{
		{
			name:     "Balance Success",
			walletID: "wallet-123",
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().GetBalance("wallet-123").Return(100.0, nil)
				return s
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"balance":100}`,
		},
		{
			name:     "Balance Error",
			walletID: "wallet-123",
			mockService: func(ctrl *gomock.Controller) *mock.MockWalletService {
				s := mock.NewMockWalletService(ctrl)
				s.EXPECT().GetBalance("wallet-123").Return(0.0, fmt.Errorf("could not retrieve balance"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"could not retrieve balance"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()
			mockService := tt.mockService(ctrl)
			logger := logrus.New()
			apiHandler := NewApiWalletHandler(mockService, logger)
			app.Get("/api/v1/wallet/:walletID", apiHandler.HandleBalance)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/wallet/"+tt.walletID, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var respBody map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&respBody)

			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, respBody["balance"], float64(100.0))
			} else {
				assert.Contains(t, respBody["error"], "could not retrieve balance")
			}
		})
	}
}
