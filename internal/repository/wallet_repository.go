package repository

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type WalletRepository interface {
	GetWalletBalance(walletID string) (float64, error)
	Deposit(walletID string, amount float64) error
	Withdraw(walletID string, amount float64) error
}

type ApiWalletRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewApiWalletRepository(db *sql.DB, logger *logrus.Logger) *ApiWalletRepository {
	return &ApiWalletRepository{
		db:     db,
		logger: logger,
	}
}

// Получение баланса кошелька по ID
func (r *ApiWalletRepository) GetWalletBalance(walletID string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM wallets WHERE wallet_id = $1`
	err := r.db.QueryRow(query, walletID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warnf("Wallet with ID %s not found", walletID)
			return 0, fmt.Errorf("wallet not found")
		}
		r.logger.Errorf("Error retrieving balance for wallet %s: %v", walletID, err)
		return 0, err
	}
	r.logger.Infof("Retrieved balance for wallet %s: %f", walletID, balance)
	return balance, nil
}

// Депозит средств на кошелек
func (r *ApiWalletRepository) Deposit(walletID string, amount float64) error {
	_, err := r.db.Exec(`UPDATE wallets SET balance = balance + $1 WHERE wallet_id = $2`, amount, walletID)
	if err != nil {
		r.logger.Errorf("Error depositing %f to wallet %s: %v", amount, walletID, err)
		return err
	}
	r.logger.Infof("Deposited %f to wallet %s", amount, walletID)
	return nil
}

// Вывод средств с кошелька
func (r *ApiWalletRepository) Withdraw(walletID string, amount float64) error {
	// Проверяем текущий баланс
	currentBalance, err := r.GetWalletBalance(walletID)
	if err != nil {
		return err
	}

	// Проверяем, достаточно ли средств для вывода
	if currentBalance < amount {
		r.logger.Warnf("Insufficient funds for wallet %s: current balance is %f, requested withdrawal is %f", walletID, currentBalance, amount)
		return fmt.Errorf("insufficient funds")
	}

	// Выполняем вывод
	_, err = r.db.Exec(`UPDATE wallets SET balance = balance - $1 WHERE wallet_id = $2`, amount, walletID)
	if err != nil {
		r.logger.Errorf("Error withdrawing %f from wallet %s: %v", amount, walletID, err)
		return err
	}
	r.logger.Infof("Withdrew %f from wallet %s", amount, walletID)
	return nil
}
