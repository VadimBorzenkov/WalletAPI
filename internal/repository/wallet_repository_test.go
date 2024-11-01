package repository

import (
	"errors"
	"testing"

	"github.com/VadimBorzenkov/WalletAPI/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	tests := []struct {
		name     string
		walletID string
		amount   float64
		wantErr  bool
	}{
		{"Valid deposit", "wallet123", 100.0, false},
		{"Zero deposit", "wallet123", 0.0, true},
		{"Negative deposit", "wallet123", -50.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.EXPECT().Deposit(tt.walletID, tt.amount).Return(nil)
			} else {
				mockRepo.EXPECT().Deposit(tt.walletID, tt.amount).Return(errors.New("deposit error"))
			}

			err := mockRepo.Deposit(tt.walletID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWalletRepository_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	tests := []struct {
		name     string
		walletID string
		amount   float64
		wantErr  bool
	}{
		{"Valid withdrawal", "wallet123", 50.0, false},
		{"Withdrawal more than balance", "wallet123", 100.0, true},
		{"Negative withdrawal", "wallet123", -30.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.EXPECT().Withdraw(tt.walletID, tt.amount).Return(nil)
			} else {
				mockRepo.EXPECT().Withdraw(tt.walletID, tt.amount).Return(errors.New("insufficient funds"))
			}

			err := mockRepo.Withdraw(tt.walletID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWalletRepository_GetWalletBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockWalletRepository(ctrl)
	tests := []struct {
		name            string
		walletID        string
		expectedBalance float64
		wantErr         bool
	}{
		{"Valid balance retrieval", "wallet123", 150.0, false},
		{"Wallet not found", "wallet999", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.EXPECT().GetWalletBalance(tt.walletID).Return(tt.expectedBalance, nil)
			} else {
				mockRepo.EXPECT().GetWalletBalance(tt.walletID).Return(0.0, errors.New("wallet not found"))
			}

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
