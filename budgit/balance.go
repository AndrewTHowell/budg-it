package budgit

// Balance is a Balance, holding BalanceAmounts of the cleared and effective balances.
type Balance struct {
	ClearedBalance   BalanceAmount
	EffectiveBalance BalanceAmount
}

// BalanceAmount stores a balance amount with no decimals, assuming the maximum number of decimal points is 2.
// e.g. £10 is stored as 100.
type BalanceAmount int64
