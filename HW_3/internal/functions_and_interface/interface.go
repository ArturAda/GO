package functions_and_interface

type BalanceRepository interface {
	Get(userID string) (int64, error)

	Deposit(userID string, amount int64) (int64, error)

	Transfer(from, to string, amount int64) (fromBalance, toBalance int64, err error)
}
