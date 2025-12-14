package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ArturAda/GO/HW_3/internal/exceptions"
	"github.com/ArturAda/GO/HW_3/internal/service"
)

type storage struct {
	mutex    sync.RWMutex
	filepath string
	data     map[string]int64
}

func New(path string) (service.BalanceRepository, error) {
	str := &storage{
		filepath: path,
		data:     make(map[string]int64),
	}
	if err := str.load(); err != nil {
		return nil, err
	}
	return str, nil
}

func (str *storage) load() error {
	str.mutex.Lock()
	defer str.mutex.Unlock()
	dir := filepath.Dir(str.filepath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.Open(str.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return str.flushLocked()
		}
		return err
	}
	defer func() { _ = f.Close() }()
	dec := json.NewDecoder(f)
	return dec.Decode(&str.data)
}

func (str *storage) flushLocked() error {
	tmp := str.filepath + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(str.data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, str.filepath)
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
	if err := str.flushLocked(); err != nil {
		return 0, fmt.Errorf("flush: %w", err)
	}
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
	if err := str.flushLocked(); err != nil {
		return 0, 0, fmt.Errorf("flush: %w", err)
	}
	return fromBalance, str.data[to], nil
}
