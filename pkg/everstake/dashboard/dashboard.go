package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const apiKeyHeader = "X-API-Key" // #nosec G101 -- key, not value

type (
	Dashboard struct {
		cli    http.Client
		url    string
		apiKey string
	}
)

func NewDashboard(dashboardURL, apiKey string, overrideHTTPClient ...http.Client) *Dashboard {
	cli := http.Client{}
	if len(overrideHTTPClient) > 0 {
		cli = overrideHTTPClient[0]
	}
	return &Dashboard{
		cli:    cli,
		url:    dashboardURL,
		apiKey: apiKey,
	}
}

func post(d *Dashboard, ctx context.Context, path string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.url+path, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(apiKeyHeader, d.apiKey)

	resp, err := d.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request to dashboard")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errBody BodyError
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			return errors.Wrapf(err, "unexpected status code from dashboard: %d", resp.StatusCode)
		}
		return errors.New(errBody.Description)
	}
	return nil
}

func get[T any](d *Dashboard, ctx context.Context, path string) (res T, err error) {
	endpoint := d.url + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return res, errors.Wrap(err, "failed to create request")
	}
	resp, err := d.cli.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "failed to send request to dashboard")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return res, errors.Errorf("unexpected status code from dashboard: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, errors.Wrap(err, "failed to decode dashboard response")
	}
	return res, nil
}
