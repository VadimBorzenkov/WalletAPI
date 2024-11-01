package service

import (
	"fmt"

	"github.com/VadimBorzenkov/WalletAPI/internal/repository"
	"github.com/sirupsen/logrus"
)

// Интерфейс сервиса для кошелька
type WalletService interface {
	GetBalance(walletID string) (float64, error)
	Deposit(walletID string, amount float64) error
	Withdraw(walletID string, amount float64) error
}

// Структура сервиса для API-кошелька
type ApiWalletService struct {
	repo   repository.WalletRepository
	logger *logrus.Logger
}

// Конструктор для ApiWalletService
func NewApiWalletService(repo repository.WalletRepository, logger *logrus.Logger) *ApiWalletService {
	return &ApiWalletService{
		repo:   repo,
		logger: logger,
	}
}

// Получение баланса кошелька
func (s *ApiWalletService) GetBalance(walletID string) (float64, error) {
	balance, err := s.repo.GetWalletBalance(walletID)
	if err != nil {
		s.logger.Errorf("Failed to get balance for wallet %s: %v", walletID, err)
		return 0, fmt.Errorf("could not retrieve balance: %w", err)
	}
	s.logger.Infof("Retrieved balance for wallet %s: %f", walletID, balance)
	return balance, nil
}

// Депозит средств на кошелек
func (s *ApiWalletService) Deposit(walletID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	err := s.repo.Deposit(walletID, amount)
	if err != nil {
		s.logger.Errorf("Failed to deposit %f to wallet %s: %v", amount, walletID, err)
		return fmt.Errorf("could not deposit amount: %w", err)
	}
	s.logger.Infof("Deposited %f to wallet %s", amount, walletID)
	return nil
}

// Вывод средств с кошелька
func (s *ApiWalletService) Withdraw(walletID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	err := s.repo.Withdraw(walletID, amount)
	if err != nil {
		s.logger.Errorf("Failed to withdraw %f from wallet %s: %v", amount, walletID, err)
		return fmt.Errorf("could not withdraw amount: %w", err)
	}
	s.logger.Infof("Withdrew %f from wallet %s", amount, walletID)
	return nil
}
