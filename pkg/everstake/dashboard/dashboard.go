package dashboard

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type (
	Dashboard struct {
		cli http.Client
		url string
	}
)

func NewDashboard(dashboardURL string, overrideHTTPClient ...http.Client) *Dashboard {
	cli := http.Client{}
	if len(overrideHTTPClient) > 0 {
		cli = overrideHTTPClient[0]
	}
	return &Dashboard{
		cli: cli,
		url: dashboardURL,
	}
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
