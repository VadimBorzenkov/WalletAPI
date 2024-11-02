package service

import (
	"fmt"
	"testing"

	"github.com/VadimBorzenkov/WalletAPI/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestApiWalletService_GetBalance тестирует метод GetBalance в ApiWalletService
func TestApiWalletService_GetBalance(t *testing.T) {
	// Создаем контроллер для управления mock-объектами
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для WalletRepository и экземпляр логгера
	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()

	// Создаем сервис ApiWalletService, используя mock репозиторий и логгер
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	expectedBalance := 100.0

	// Определяем ожидание: вызов метода GetWalletBalance должен вернуть баланс 100.0
	mockRepo.EXPECT().GetWalletBalance(walletID).Return(expectedBalance, nil)

	// Вызываем метод GetBalance и проверяем результат
	balance, err := service.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

// TestApiWalletService_Deposit тестирует успешный случай метода Deposit в ApiWalletService
func TestApiWalletService_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для WalletRepository и логгер
	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()

	// Создаем сервис
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 50.0

	// Настраиваем mock: метод Deposit должен завершиться без ошибок
	mockRepo.EXPECT().Deposit(walletID, amount).Return(nil)

	// Вызываем метод Deposit и проверяем, что ошибок нет
	err := service.Deposit(walletID, amount)
	assert.NoError(t, err)
}

// TestApiWalletService_Deposit_NegativeAmount проверяет ошибку при попытке депозита отрицательной суммы
func TestApiWalletService_Deposit_NegativeAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := -50.0

	// Вызываем метод Deposit с отрицательной суммой и проверяем, что возникает ошибка
	err := service.Deposit(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, "deposit amount must be positive", err.Error())

	// Проверяем, что метод Deposit не должен был вызываться с какими-либо аргументами
	mockRepo.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(0)
}

// TestApiWalletService_Withdraw тестирует успешный случай метода Withdraw в ApiWalletService
func TestApiWalletService_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 30.0

	// Ожидаем, что вызов Withdraw выполнится успешно
	mockRepo.EXPECT().Withdraw(walletID, amount).Return(nil)

	// Вызываем метод Withdraw и проверяем, что ошибок нет
	err := service.Withdraw(walletID, amount)
	assert.NoError(t, err)
}

// TestApiWalletService_Withdraw_NegativeAmount проверяет ошибку при попытке вывода отрицательной суммы
func TestApiWalletService_Withdraw_NegativeAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := -30.0

	// Вызываем метод Withdraw с отрицательной суммой и проверяем, что возникает ошибка
	err := service.Withdraw(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, "withdrawal amount must be positive", err.Error())

	// Проверяем, что метод Withdraw не должен был вызываться с какими-либо аргументами
	mockRepo.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(0)
}

// TestApiWalletService_Withdraw_Failure тестирует случай неудачного вывода средств (например, недостаточно средств)
func TestApiWalletService_Withdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	logger := logrus.New()
	service := NewApiWalletService(mockRepo, logger)

	walletID := "test_wallet"
	amount := 30.0

	// Ожидаем, что метод Withdraw вернет ошибку "insufficient funds"
	mockRepo.EXPECT().Withdraw(walletID, amount).Return(fmt.Errorf("insufficient funds"))

	// Вызываем метод Withdraw и проверяем, что ошибка соответствует ожиданию
	err := service.Withdraw(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, "could not withdraw amount: insufficient funds", err.Error())
}
