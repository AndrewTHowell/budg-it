package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/integrations/starling"
	"go.uber.org/zap"
)

const starlingIntegrationID = "starling"

type Client struct {
	log    *zap.SugaredLogger
	client *starling.ClientWithResponses
}

func NewStarlingClient(log *zap.SugaredLogger, url, apiToken string) (*Client, error) {
	log.Debugw("Starting Starling client", zap.String("url", url))

	client, err := starling.NewClientWithResponses(
		url,
		starling.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("initialising Starling client: %w", err)
	}
	return &Client{
		log:    log,
		client: client,
	}, nil
}

func (c Client) ID() string { return starlingIntegrationID }

func (c Client) GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error) {
	c.log.Debug("Getting external Starling accounts")

	resp, err := c.client.GetAccountsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Accounts: %w", err)
	}
	if resp.JSON4XX != nil {
		return nil, fmt.Errorf("getting Accounts: %w", format4XXError(resp.JSON4XX))
	}
	c.log.Debugw("Retrieved external Starling accounts", zap.Int("number_of_accounts", len(*resp.JSON200.Accounts)))

	accounts := make([]*budgit.ExternalAccount, 0, len(*resp.JSON200.Accounts))
	for _, account := range *resp.JSON200.Accounts {
		c.log.Debugw("Getting account balance of Starling account",
			zap.String("account_id", account.AccountUid.String()),
			zap.String("name", *account.Name),
		)

		resp, err := c.client.GetAccountBalanceWithResponse(ctx, *account.AccountUid)
		if err != nil {
			return nil, fmt.Errorf("getting Accounts: %w", err)
		}
		if resp.JSON4XX != nil {
			return nil, fmt.Errorf("getting Accounts: %w", format4XXError(resp.JSON4XX))
		}
		accounts = append(accounts, &budgit.ExternalAccount{
			ID:            account.AccountUid.String(),
			IntegrationID: starlingIntegrationID,
			Balance: budgit.Balance{
				ClearedBalance:   budgit.BalanceAmount(resp.JSON200.TotalClearedBalance.MinorUnits),
				EffectiveBalance: budgit.BalanceAmount(resp.JSON200.TotalEffectiveBalance.MinorUnits),
			},
		})
	}
	return accounts, nil
}

func (c Client) GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error) {
	c.log.Debugw("Getting external Starling account", zap.String("account_id", externalID))

	accounts, err := c.GetExternalAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Account %q: %w", externalID, err)
	}
	idx := slices.IndexFunc(accounts, func(a *budgit.ExternalAccount) bool {
		return a.ID == externalID
	})
	if idx == -1 {
		return nil, ErrAccountNotFound
	}
	return accounts[idx], nil
}

func format4XXError(errResp *starling.ErrorResponse) error {
	errs := make([]error, 0, len(*errResp.Errors))
	for _, errDetail := range *errResp.Errors {
		msg := *errDetail.Message
		errs = append(errs, fmt.Errorf(msg))
	}
	return errors.Join(errs...)
}
