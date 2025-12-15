package memory

import (
	"sync"

	"github.com/ArturAda/GO/HW_3/internal/exceptions"
	"github.com/ArturAda/GO/HW_3/internal/service"
)

type storage struct {
	mutex sync.RWMutex
	data  map[string]int64
}

func New() service.BalanceRepository {
	return &storage{
		data: make(map[string]int64),
	}
}

func (str *storage) Get(userID string) (int64, error) {
	str.mutex.RLock()
	defer str.mutex.RUnlock()
	balance, ok := str.data[userID]
	if !ok {
		return 0, exceptions.ErrNotFound
	}
	return balance, nil
}

func (str *storage) Deposit(userID string, amount int64) (int64, error) {
	str.mutex.Lock()
	defer str.mutex.Unlock()
	str.data[userID] += amount
	return str.data[userID], nil
}

func (str *storage) Transfer(from, to string, amount int64) (int64, int64, error) {
	str.mutex.Lock()
	defer str.mutex.Unlock()
	fromBalance, ok := str.data[from]
	if !ok {
		return 0, 0, exceptions.ErrNotFound
	}
	if fromBalance < amount {
		return fromBalance, 0, exceptions.ErrInsufficientFunds
	}
	fromBalance -= amount
	str.data[to] += amount
	str.data[from] = fromBalance
	return fromBalance, str.data[to], nil
}
