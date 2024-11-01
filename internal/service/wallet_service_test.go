package service_test

import (
	"fmt"
	"testing"

	"github.com/VadimBorzenkov/WalletAPI/internal/repository/mock"
	"github.com/VadimBorzenkov/WalletAPI/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApiWalletService_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	expectedBalance := 100.0

	// Настройка ожиданий мока для репозитория
	mockRepo.EXPECT().GetWalletBalance(walletID).Return(expectedBalance, nil)

	// Вызов метода
	balance, err := service.GetBalance(walletID)

	// Проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestApiWalletService_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 50.0

	// Настройка ожиданий мока для репозитория
	mockRepo.EXPECT().Deposit(walletID, amount).Return(nil)

	// Вызов метода
	err := service.Deposit(walletID, amount)

	// Проверка результатов
	assert.NoError(t, err)
}

func TestApiWalletService_Deposit_NegativeAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := -50.0

	// Вызов метода Deposit
	err := service.Deposit(walletID, amount)

	// Проверка, что ошибка возвращается правильно
	assert.Error(t, err)
	assert.Equal(t, "deposit amount must be positive", err.Error())

	// Убедитесь, что репозиторий не вызывался
	mockRepo.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(0)
}

func TestApiWalletService_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 30.0

	// Настройка ожиданий мока для репозитория
	mockRepo.EXPECT().Withdraw(walletID, amount).Return(nil)

	// Вызов метода
	err := service.Withdraw(walletID, amount)

	// Проверка результатов
	assert.NoError(t, err)
}

func TestApiWalletService_Withdraw_NegativeAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := -30.0

	// Вызов метода
	err := service.Withdraw(walletID, amount)

	// Проверка, что ошибка возвращается правильно
	assert.Error(t, err)
	assert.Equal(t, "withdrawal amount must be positive", err.Error())

	// Убедитесь, что репозиторий не вызывался
	mockRepo.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(0)
}

func TestApiWalletService_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := service.NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 30.0

	// Настройка ожиданий мока с ошибкой
	mockRepo.EXPECT().Withdraw(walletID, amount).Return(fmt.Errorf("insufficient funds"))

	// Вызов метода
	err := service.Withdraw(walletID, amount)

	// Проверка результатов
	assert.Error(t, err)
	assert.Equal(t, "could not withdraw amount: insufficient funds", err.Error())
}
