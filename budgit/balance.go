package budgit

import "fmt"

// Balance is a Balance, holding BalanceAmounts of the cleared and effective balances.
type Balance struct {
	ClearedBalance   BalanceAmount
	EffectiveBalance BalanceAmount
}

func (b Balance) Add(balance Balance) Balance {
	b.ClearedBalance += balance.ClearedBalance
	b.EffectiveBalance += balance.EffectiveBalance
	return b
}

func (b Balance) AddAmount(amount BalanceAmount, cleared bool) Balance {
	b.EffectiveBalance += amount
	if cleared {
		b.ClearedBalance += amount
	}
	return b
}

// BalanceAmount stores a balance amount with no decimals, assuming the maximum number of decimal points is 2.
// e.g. £10 is stored as 100.
type BalanceAmount int64

func (b BalanceAmount) String() string {
	return fmt.Sprintf("%d.%d", b/100, b%100)
}
