package repository

import (
	"errors"
	"testing"

	"github.com/VadimBorzenkov/WalletAPI/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestWalletRepository_Deposit тестирует метод Deposit в WalletRepository
func TestWalletRepository_Deposit(t *testing.T) {
	// Создаем контроллер mock для управления моками
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для WalletRepository
	mockRepo := mock.NewMockWalletRepository(ctrl)

	// Определяем тестовые случаи
	tests := []struct {
		name     string  // Название теста
		walletID string  // ID кошелька
		amount   float64 // Сумма депозита
		wantErr  bool    // Ожидаемая ошибка
	}{
		{"Valid deposit", "wallet123", 100.0, false},
		{"Zero deposit", "wallet123", 0.0, true},
		{"Negative deposit", "wallet123", -50.0, true},
	}

	// Выполняем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем ожидания на основании условия wantErr
			if !tt.wantErr {
				mockRepo.EXPECT().Deposit(tt.walletID, tt.amount).Return(nil)
			} else {
				mockRepo.EXPECT().Deposit(tt.walletID, tt.amount).Return(errors.New("deposit error"))
			}

			// Вызываем метод Deposit и проверяем результат
			err := mockRepo.Deposit(tt.walletID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWalletRepository_Withdraw тестирует метод Withdraw в WalletRepository
func TestWalletRepository_Withdraw(t *testing.T) {
	// Создаем контроллер mock для управления моками
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для WalletRepository
	mockRepo := mock.NewMockWalletRepository(ctrl)

	// Определяем тестовые случаи
	tests := []struct {
		name     string  // Название теста
		walletID string  // ID кошелька
		amount   float64 // Сумма вывода
		wantErr  bool    // Ожидаемая ошибка
	}{
		{"Valid withdrawal", "wallet123", 50.0, false},
		{"Withdrawal more than balance", "wallet123", 100.0, true},
		{"Negative withdrawal", "wallet123", -30.0, true},
	}

	// Выполняем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем ожидания на основании условия wantErr
			if !tt.wantErr {
				mockRepo.EXPECT().Withdraw(tt.walletID, tt.amount).Return(nil)
			} else {
				mockRepo.EXPECT().Withdraw(tt.walletID, tt.amount).Return(errors.New("insufficient funds"))
			}

			// Вызываем метод Withdraw и проверяем результат
			err := mockRepo.Withdraw(tt.walletID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWalletRepository_GetWalletBalance тестирует метод GetWalletBalance в WalletRepository
func TestWalletRepository_GetWalletBalance(t *testing.T) {
	// Создаем контроллер mock для управления моками
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для WalletRepository
	mockRepo := mock.NewMockWalletRepository(ctrl)

	// Определяем тестовые случаи
	tests := []struct {
		name            string  // Название теста
		walletID        string  // ID кошелька
		expectedBalance float64 // Ожидаемый баланс
		wantErr         bool    // Ожидаемая ошибка
	}{
		{"Valid balance retrieval", "wallet123", 150.0, false},
		{"Wallet not found", "wallet999", 0.0, true},
	}

	// Выполняем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем ожидания на основании условия wantErr
			if !tt.wantErr {
				mockRepo.EXPECT().GetWalletBalance(tt.walletID).Return(tt.expectedBalance, nil)
			} else {
				mockRepo.EXPECT().GetWalletBalance(tt.walletID).Return(0.0, errors.New("wallet not found"))
			}

			// Вызываем метод GetWalletBalance и проверяем результат
			balance, err := mockRepo.GetWalletBalance(tt.walletID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0.0, balance)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, balance)
			}
		})
	}
}
