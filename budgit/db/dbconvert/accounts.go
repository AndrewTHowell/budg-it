package dbconvert

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
)

func ToAccounts(dbAccounts ...*db.Account) []*budgit.Account {
	accounts := make([]*budgit.Account, 0, len(dbAccounts))
	for _, dbAccount := range dbAccounts {
		accounts = append(accounts, toAccount(dbAccount))
	}
	return accounts
}

func toAccount(account *db.Account) *budgit.Account {
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

func FromAccounts(accounts ...*budgit.Account) []*db.Account {
	dbAccounts := make([]*db.Account, 0, len(accounts))
	for _, account := range accounts {
		dbAccounts = append(dbAccounts, fromAccount(account))
	}
	return dbAccounts
}

func fromAccount(account *budgit.Account) *db.Account {
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
