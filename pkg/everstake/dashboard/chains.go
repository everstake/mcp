package dashboard

import "github.com/pkg/errors"

type (
	Chain struct {
		ChainID           int     `json:"chain_id" jsonschema:"Unique chain identifier"`
		Chain             string  `json:"chain" jsonschema:"Chain name"`
		Title             string  `json:"title,omitempty" jsonschema:"Chain title for display"`
		Description       string  `json:"description,omitempty" jsonschema:"Chain description"`
		Keywords          string  `json:"keywords,omitempty" jsonschema:"SEO keywords"`
		Url               string  `json:"url,omitempty" jsonschema:"URL slug for chain"`
		Status            string  `json:"status,omitempty" jsonschema:"Chain status"`
		LogoWhite         string  `json:"logo_white,omitempty"`
		LogoBlack         string  `json:"logo_black,omitempty"`
		Fee               float64 `json:"fee,omitempty" jsonschema:"Staking fee percentage"`
		Apr               float64 `json:"apr,omitempty" jsonschema:"Annual percentage rate for staking"`
		Website           string  `json:"website,omitempty" jsonschema:"Chain official website"`
		TotalDelegated    float64 `json:"total_delegated,omitempty" jsonschema:"Total amount delegated"`
		RewardFrequency   string  `json:"reward_frequency,omitempty" jsonschema:"Frequency of reward distribution"`
		UnboundPeriod     string  `json:"unbound_period,omitempty" jsonschema:"Unbounding period duration"`
		Picture           string  `json:"picture,omitempty" jsonschema:"Chain logo image URL"`
		Index             int     `json:"index,omitempty" jsonschema:"Display order index"`
		B2EnableEmail     bool    `json:"b2_enable_email,omitempty"`
		B2Text            string  `json:"b2_text,omitempty"`
		B3EnableCall      bool    `json:"b3_enable_call,omitempty"`
		B3Text            string  `json:"b3_text,omitempty"`
		CreatedAt         string  `json:"created_at,omitempty" jsonschema:"Timestamp when chain was created"`
		UpdatedAt         string  `json:"updated_at,omitempty" jsonschema:"Timestamp of last update"`
		CurrencyID        int     `json:"currency_id,omitempty" jsonschema:"Associated currency identifier"`
		CurrencyCode      string  `json:"currency_code,omitempty" jsonschema:"Currency code (e.g., eth, btc)"`
		CurrencyCreatedAt string  `json:"currency_created_at,omitempty" jsonschema:"Timestamp when currency was created"`
		CurrencyUpdatedAt string  `json:"currency_updated_at,omitempty" jsonschema:"Timestamp of currency last update"`
	}
)

func (d *Dashboard) GetChains() ([]Chain, error) {
	res, err := get[[]Chain](d, "/chains")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch")
	}
	return res, nil
}
