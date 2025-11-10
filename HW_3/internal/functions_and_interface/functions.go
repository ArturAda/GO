package functions_and_interface

import (
	"fmt"

	"github.com/ArturAda/GO/HW_3/internal/exceptions"
)

type BalanceService struct {
	storage BalanceRepository
}

func NewBalanceService(newStorage BalanceRepository) *BalanceService {
	return &BalanceService{storage: newStorage}
}

func (service *BalanceService) GetBalance(userID string) (int64, error) {
	return service.storage.Get(userID)
}

func (service *BalanceService) Deposit(userID string, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, exceptions.ErrInvalidAmount
	}
	return service.storage.Deposit(userID, amount)
}

func (service *BalanceService) Transfer(from, to string, amount int64) (int64, int64, error) {
	if amount <= 0 {
		return 0, 0, exceptions.ErrInvalidAmount
	}
	if from == to {
		return 0, 0, exceptions.ErrSelfTransfer
	}
	fromBalance, toBalance, err := service.storage.Transfer(from, to, amount)
	return fromBalance, toBalance, err
}

func FormatKopeck(kopecks int64) string {
	sign := ""
	if kopecks < 0 {
		sign = "-"
		kopecks *= -1
	}
	return fmt.Sprintf("%s%d.%02d", sign, kopecks/100, kopecks%100)
}
