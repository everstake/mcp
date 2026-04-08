package dashboard

import (
	"context"

	"github.com/pkg/errors"
)

type (
	DatafeedResponse[T any] struct {
		Value T      `json:"value"`
		ID    string `json:"id"` // method name
	}
	GlobalStatsValue struct {
		NetworksSupported    string  `json:"networksSupported" jsonschema:"Number of networks supported"`
		NumberOfDelegators   int64   `json:"numberOfDelegators" jsonschema:"Number of delegators across all networks"`
		TotalStakeUSD        int64   `json:"totalStakeUsd" jsonschema:"Total amount staked in USD"`
		RewardsPaidUSD       int64   `json:"rewardsPaidUsd" jsonschema:"Total rewards paid out in USD"`
		InstitutionalClients int64   `json:"institutionalClients" jsonschema:"Number of institutional clients"`
		ClientRetention      int64   `json:"clientRetention" jsonschema:"Client retention percentage"`
		ServiceUptime        float64 `json:"serviceUptime" jsonschema:"Service uptime percentage"`
	}
)

func (d *Dashboard) GetGlobalStats(ctx context.Context) (*GlobalStatsValue, error) {
	res, err := get[DatafeedResponse[GlobalStatsValue]](d, ctx, "/datafeed/global-stats")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch")
	}
	return &res.Value, nil
}
