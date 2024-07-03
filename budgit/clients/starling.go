package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/integrations/starling"
)

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

func (c Client) GetAccounts(ctx context.Context) ([]*starling.AccountV2, error) {

	resp, err := c.client.GetAccountsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Accounts: %w", err)
	}
	if resp.JSON4XX != nil {
		return nil, fmt.Errorf("getting Accounts: %w", format4XXError(resp.JSON4XX))
	}

	for _, account := range *resp.JSON200.Accounts {
		fmt.Println(fmt.Sprintf("%+v", account))
	}
	return nil, nil
}

func toAccount(starlingAccount *starling.AccountV2) (*budgit.Account, error) {
	acct := budgit.NewAccount("", "")
	return acct, nil
}

func format4XXError(errResp *starling.ErrorResponse) error {
	errs := make([]error, 0, len(*errResp.Errors))
	for _, errDetail := range *errResp.Errors {
		msg := *errDetail.Message
		errs = append(errs, fmt.Errorf(msg))
	}
	return errors.Join(errs...)
}
