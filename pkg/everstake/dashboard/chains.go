package dashboard

import (
	"context"

	"github.com/pkg/errors"
)

type (
	Chain struct {
		Chain             string  `json:"chain" jsonschema:"Chain name"`
		Title             string  `json:"title,omitempty" jsonschema:"Chain title for display"`
		Description       string  `json:"description,omitempty" jsonschema:"Chain description"`
		Keywords          string  `json:"keywords,omitempty" jsonschema:"SEO keywords"`
		URL               string  `json:"url,omitempty" jsonschema:"URL slug for chain"`
		Status            string  `json:"status,omitempty" jsonschema:"Chain status"`
		LogoWhite         string  `json:"logo_white,omitempty"`
		LogoBlack         string  `json:"logo_black,omitempty"`
		Website           string  `json:"website,omitempty" jsonschema:"Chain official website"`
		RewardFrequency   string  `json:"reward_frequency,omitempty" jsonschema:"Frequency of reward distribution"`
		UnboundPeriod     string  `json:"unbound_period,omitempty" jsonschema:"Unbounding period duration"`
		Picture           string  `json:"picture,omitempty" jsonschema:"Chain logo image URL"`
		B2Text            string  `json:"b2_text,omitempty"`
		B3Text            string  `json:"b3_text,omitempty"`
		CreatedAt         string  `json:"created_at,omitempty" jsonschema:"Timestamp when chain was created"`
		UpdatedAt         string  `json:"updated_at,omitempty" jsonschema:"Timestamp of last update"`
		CurrencyCode      string  `json:"currency_code,omitempty" jsonschema:"Currency code (e.g., eth, btc)"`
		CurrencyCreatedAt string  `json:"currency_created_at,omitempty" jsonschema:"Timestamp when currency was created"`
		CurrencyUpdatedAt string  `json:"currency_updated_at,omitempty" jsonschema:"Timestamp of currency last update"`
		Fee               float64 `json:"fee,omitempty" jsonschema:"Staking fee percentage"`
		Apr               float64 `json:"apr,omitempty" jsonschema:"Annual percentage rate for staking"`
		TotalDelegated    float64 `json:"total_delegated,omitempty" jsonschema:"Total amount delegated"`
		ChainID           int     `json:"chain_id" jsonschema:"Unique chain identifier"`
		Index             int     `json:"index,omitempty" jsonschema:"Display order index"`
		CurrencyID        int     `json:"currency_id,omitempty" jsonschema:"Associated currency identifier"`
		B2EnableEmail     bool    `json:"b2_enable_email,omitempty"`
		B3EnableCall      bool    `json:"b3_enable_call,omitempty"`
	}
)

func (d *Dashboard) GetChains(ctx context.Context) ([]Chain, error) {
	res, err := get[[]Chain](d, ctx, "/chains")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch")
	}
	return res, nil
}
