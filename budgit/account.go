package budgit

import "time"

// Account is an account internal to budgit.
type Account struct {
	ID              string
	Name            string
	Balance         Balance
	ExternalAccount *ExternalAccount
}

// ExternalAccount is an Account representing some real, external Account that is attached to a budgit Account.
type ExternalAccount struct {
	ID                string
	Name              string
	IntegrationID     string
	LastSyncTimestamp time.Time
	Balance           Balance
}
