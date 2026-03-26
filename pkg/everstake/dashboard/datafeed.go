package dashboard

import "github.com/pkg/errors"

type (
	DatafeedResponse[T any] struct {
		ID    string `json:"id"` // method name
		Value T      `json:"value"`
	}
	GlobalStatsValue struct {
		NumberOfDelegators   int64   `json:"numberOfDelegators" jsonschema:"Number of delegators across all networks"`
		TotalStakeUSD        int64   `json:"totalStakeUsd" jsonschema:"Total amount staked in USD"`
		NetworksSupported    string  `json:"networksSupported" jsonschema:"Number of networks supported"`
		RewardsPaidUSD       int64   `json:"rewardsPaidUsd" jsonschema:"Total rewards paid out in USD"`
		InstitutionalClients int64   `json:"institutionalClients" jsonschema:"Number of institutional clients"`
		ClientRetention      int64   `json:"clientRetention" jsonschema:"Client retention percentage"`
		ServiceUptime        float64 `json:"serviceUptime" jsonschema:"Service uptime percentage"`
	}
)

func (d *Dashboard) GetGlobalStats() (*GlobalStatsValue, error) {
	res, err := get[DatafeedResponse[GlobalStatsValue]](d, "/datafeed/global-stats")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch")
	}
	return &res.Value, nil
}
