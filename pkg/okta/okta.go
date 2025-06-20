package okta

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/iamNilotpal/iam/internal/config"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type Client struct {
	sdk *okta.APIClient
}

func NewClient(cfg *config.OktaConfig) (*Client, error) {
	oktaConfig, err := okta.NewConfiguration(
		okta.WithToken(cfg.APIToken),
		okta.WithOrgUrl(fmt.Sprintf("https://%s", cfg.Domain)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create okta config : %w", err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
		},
	}

	oktaConfig.HTTPClient = httpClient
	return &Client{sdk: okta.NewAPIClient(oktaConfig)}, nil
}

func (c *Client) SDK() *okta.APIClient {
	return c.sdk
}

func (c *Client) TestConnection(ctx context.Context) error {
	_, resp, err := c.sdk.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return fmt.Errorf("failed to validate Okta connection: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("okta api returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
