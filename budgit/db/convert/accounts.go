package convert

import (
	"time"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToAccount(account *db.Account) *budgit.Account {
	var externalAccount *budgit.ExternalAccount
	if account.ExternalID.Valid {
		externalAccount = &budgit.ExternalAccount{
			ID:                account.ExternalID.String,
			IntegrationID:     account.ExternalIntegrationID.String,
			LastSyncTimestamp: account.ExternalLastSyncTimestamp.Time,
			Balance: budgit.Balance{
				ClearedBalance:   budgit.BalanceAmount(account.ExternalClearedBalance.Int64),
				EffectiveBalance: budgit.BalanceAmount(account.ExternalEffectiveBalance.Int64),
			},
		}
	}
	return &budgit.Account{
		ID:   account.ID.String,
		Name: account.Name.String,
		Balance: budgit.Balance{
			ClearedBalance:   budgit.BalanceAmount(account.ClearedBalance.Int64),
			EffectiveBalance: budgit.BalanceAmount(account.EffectiveBalance.Int64),
		},
		ExternalAccount: externalAccount,
	}
}

func FromAccount(account *budgit.Account) *db.Account {
	dbAccount := &db.Account{
		ID:               toText(account.ID),
		Name:             toText(account.Name),
		ClearedBalance:   toInt8(int64(account.Balance.ClearedBalance)),
		EffectiveBalance: toInt8(int64(account.Balance.EffectiveBalance)),
	}
	if account.ExternalAccount != nil {
		dbAccount.ExternalID = toText(account.ExternalAccount.ID)
		dbAccount.ExternalIntegrationID = toText(account.ExternalAccount.IntegrationID)
		dbAccount.ExternalLastSyncTimestamp = toTimestamptz(account.ExternalAccount.LastSyncTimestamp)
		dbAccount.ExternalClearedBalance = toInt8(int64(account.ExternalAccount.Balance.ClearedBalance))
		dbAccount.ExternalEffectiveBalance = toInt8(int64(account.ExternalAccount.Balance.EffectiveBalance))

	}
	return dbAccount
}

func toText(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: str, Valid: true}
}

func toInt8(i int64) pgtype.Int8 {
	if i == 0 {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: i, Valid: true}
}

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}
