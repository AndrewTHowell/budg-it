package convert

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
)

func ToTransaction(transaction *db.Transaction) *budgit.Transaction {
	return &budgit.Transaction{
		ID:              transaction.ID.String,
		EffectiveDate:   transaction.EffectiveDate.Time,
		AccountID:       transaction.AccountID.String,
		PayeeID:         transaction.PayeeID.String,
		IsPayeeInternal: transaction.IsPayeeInternal.Bool,
		Amount:          budgit.BalanceAmount(transaction.Amount.Int64),
		Cleared:         transaction.Cleared.Bool,
	}
}

func FromTransaction(transaction *budgit.Transaction) *db.Transaction {
	return &db.Transaction{
		ID:              toText(transaction.ID),
		EffectiveDate:   toDate(transaction.EffectiveDate),
		AccountID:       toText(transaction.AccountID),
		PayeeID:         toText(transaction.PayeeID),
		IsPayeeInternal: toBool(transaction.IsPayeeInternal),
		Amount:          toInt8(int64(transaction.Amount)),
		Cleared:         toBool(transaction.Cleared),
	}
}
