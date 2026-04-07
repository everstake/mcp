package dashboard

import (
	"bytes"
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

func NewDashboard(dashboardUrl string, overrideHttpClient ...http.Client) *Dashboard {
	cli := http.Client{}
	if len(overrideHttpClient) > 0 {
		cli = overrideHttpClient[0]
	}
	return &Dashboard{
		cli: cli,
		url: dashboardUrl,
	}
}

func get[T any](d *Dashboard, path string) (res T, err error) {
	endpoint := d.url + path
	resp, err := d.cli.Get(endpoint)
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

func post[T any](d *Dashboard, path string, body any) (res T, err error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return res, errors.Wrap(err, "failed to marshal request body")
	}

	endpoint := d.url + path
	resp, err := d.cli.Post(endpoint, "application/json", bytes.NewReader(payload))
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
