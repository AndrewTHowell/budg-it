package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/integrations/starling"
)

const providerStarling = "starling"

type Client struct {
	client *starling.ClientWithResponses
}

func NewStarlingClient(url, apiToken string) (*Client, error) {
	client, err := starling.NewClientWithResponses(
		url,
		starling.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("initialising starling client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c Client) GetAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error) {
	resp, err := c.client.GetAccountsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Accounts: %w", err)
	}
	if resp.JSON4XX != nil {
		return nil, fmt.Errorf("getting Accounts: %w", format4XXError(resp.JSON4XX))
	}

	accounts := make([]*budgit.ExternalAccount, 0, len(*resp.JSON200.Accounts))
	for _, account := range *resp.JSON200.Accounts {
		accounts = append(accounts, toAccount(&account))
	}
	return accounts, nil
}

func (c Client) GetAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Account %q: %w", externalID, err)
	}
	idx := slices.IndexFunc(accounts, func(a *budgit.ExternalAccount) bool {
		return a.ExternalID == externalID
	})
	if idx == -1 {
		return nil, ErrAccountNotFound
	}
	return accounts[idx], nil
}

func toAccount(starlingAccount *starling.AccountV2) *budgit.ExternalAccount {
	return budgit.NewExternalAccount(providerStarling, starlingAccount.AccountUid.String())
}

func format4XXError(errResp *starling.ErrorResponse) error {
	errs := make([]error, 0, len(*errResp.Errors))
	for _, errDetail := range *errResp.Errors {
		msg := *errDetail.Message
		errs = append(errs, fmt.Errorf(msg))
	}
	return errors.Join(errs...)
}
