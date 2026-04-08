package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type PDLead struct {
	FirstName              string `json:"first_name" jsonschema:"First name of the contact" required:"true"`
	LastName               string `json:"last_name" jsonschema:"Last name of the contact" required:"false"`
	Email                  string `json:"email" jsonschema:"Work email address of the contact" required:"true"`
	CompanyName            string `json:"company_name,omitempty" jsonschema:"Name of the company"`
	CompanyType            string `json:"company_type,omitempty" jsonschema:"Type of company" enum:"[Custodian,Wallet,Exchange,Asset Manager,Treasury,Family Office,Individual Investor,Other"`
	CompanySite            string `json:"company_site,omitempty" jsonschema:"Company website URL (optional, can be inferred from email domain)"`
	JobTitle               string `json:"job_title,omitempty" jsonschema:"Job title of the contact"`
	PrimaryRegion          string `json:"primary_region,omitempty" jsonschema:"Primary region of the company: North America, Europe, LATAM, APAC, MENA"`
	ProductOfInterest      string `json:"product_of_interest,omitempty" jsonschema:"Product of interest: Staking, Vaults, Data & Analytics"`
	CustodySolution        string `json:"custody_solution,omitempty" jsonschema:"Custody solution: Anchorage, BitGo, Circle, Coinbase Custody, Copper, Fireblocks, Gemini, Ledger Enterprise, MetaMask, Sygnum, Taurus Group, Self Custody, Other"`
	ApproximateStakeSize   string `json:"approximate_stake_size,omitempty" jsonschema:"Approximate stake size: <$1M, $2M–$5M, $5M–$10M, >$10M"`
	ImplementationTimeline string `json:"implementation_timeline,omitempty" jsonschema:"Implementation timeline: ASAP, 1–3 Months, 3+ Months"`
	LeadSource             string `json:"lead_source,omitempty" jsonschema:"Source of the lead — set automatically by the server, do not ask the user"`
	SendNewsletter         bool   `json:"send_newsletter,omitempty" jsonschema:"Whether the contact consents to receive Everstake marketing communications"`
}

// CreatePDLead submits a lead to the dashboard backend. The backend returns 200 with an empty body on success.
func (d *Dashboard) CreatePDLead(ctx context.Context, lead *PDLead) error {
	payload, err := json.Marshal(lead)
	if err != nil {
		return errors.Wrap(err, "failed to marshal lead")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.url+"/pd/lead", bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request to dashboard")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code from dashboard: %d", resp.StatusCode)
	}

	return nil
}
