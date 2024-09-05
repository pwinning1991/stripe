package stripe_test

import (
	"flag"
	"fmt"
	"github.com/pwinning1991/stripe"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	apiKey string
	update bool
)

const (
	tokenAmex        = "tok_amex"
	tokenInvalid     = "tok_alsdkjfa"
	tokenCardExpired = "tok_chargeDeclinedExpiredCard"
)

//sk_test_4eC39HqLyjWDarjtT1zdp7dc

func init() {
	flag.StringVar(&apiKey, "key", "", "Your TEST secret key for the Stripe API."+
		"If present, integartion tests will be run using this API key.")
	flag.BoolVar(&update, "update", false, "Update the responses used in local tests.")
}

func TestClient_Local(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"id": "cus_Dy08A71mTUtap1",
			"object": "customer",
			"account_balance": 0,
			"created": 1542124116,
			"currency": "usd",
			"default_source": null,
			"delinquent": false,
			"description": null,
			"discount": null,
			"email": null,
			"invoice_prefix": "1C04993",
			"livemode": false,
			"metadata": {
			},
			"shipping": null,
			"sources": {
				"object": "list",
				"data": [

				],
				"has_more": false,
				"total_count": 0,
				"url": "/v1/customers/cus_Dy08A71mTUtap1/sources"
			},
			"subscriptions": {
				"object": "list",
				"data": [

				],
				"has_more": false,
				"total_count": 0,
				"url": "/v1/customers/cus_Dy08A71mTUtap1/subscriptions"
			},
			"tax_info": null,
			"tax_info_verification": null
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	c := stripe.Client{
		Key:     "gibberish-key",
		BaseURL: server.URL,
	}
	_, err := c.Customer("random token", "random email")
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No TEST API key provided.")
	}

	type checkFn func(*testing.T, *stripe.Customer, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}

	}
	hasErrType := func(typee string) checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Log("Unexpected error type", se)
				t.Fatalf("err isn't a stripe.Error")
			}
			if se.Type != typee {
				t.Errorf("stripe.Error.Type = %s; want %s", se.Type, typee)
			}
		}
	}
	hasIDPrefix := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if !strings.HasPrefix(cus.ID, "cus_") {
				t.Errorf("customer.ID = %s; want prefix %s", cus.ID, "cus_")
			}
		}
	}

	hasCardDefaultSource := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if !strings.HasPrefix(cus.DefaultSource, "card_") {
				t.Errorf("customer.DefaultSource = %s; want prefix %s", cus.DefaultSource, "card_")
			}
		}
	}

	hasEmail := func(email string) checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if cus.Email != email {
				t.Errorf("cus.email = %s, want %s", cus.Email, email)
			}

		}

	}

	c := stripe.Client{
		Key: apiKey,
	}

	test := map[string]struct {
		token  string
		email  string
		checks []checkFn
	}{
		"valid customer with an amex": {
			token:  tokenAmex,
			email:  "test@example.com",
			checks: check(hasNoErr(), hasIDPrefix(), hasCardDefaultSource(), hasEmail("test@example.com")),
		},
		"invalid token": {
			token:  tokenInvalid,
			email:  "test@example.com",
			checks: check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
		"expired card": {
			token:  tokenCardExpired,
			email:  "test@example.com",
			checks: check(hasErrType(stripe.ErrTypeCardError)),
		},
	}
	for name, tc := range test {
		t.Run(name, func(t *testing.T) {
			cus, err := c.Customer(tc.token, tc.email)
			for _, check := range tc.checks {
				check(t, cus, err)
			}
		})
	}

}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No TEST API key provided.")
	}

	type checkFn func(*testing.T, *stripe.Charge, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, _ *stripe.Charge, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}

	}
	hasAmount := func(amount int) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			if charge.Amount != amount {
				t.Errorf("charge.Amount = %d; want %d", charge.Amount, amount)
			}
		}
	}
	hasErrType := func(typee string) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Log("Unexpected error type", se)
				t.Fatalf("err ins't a stripe.Error")
			}
			if se.Type != typee {
				t.Errorf("stripe.Error.Type = %s; want %s", se.Type, typee)
			}
		}
	}

	c := stripe.Client{
		Key: apiKey,
	}
	//create customer for the test
	tok := tokenAmex
	email := "test@test.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Errorf("Customer() err = %v; want %v", err, nil)
	}

	tests := map[string]struct {
		customerID string
		amount     int
		checks     []checkFn
	}{
		"valid charge": {
			customerID: cus.ID,
			amount:     1234,
			checks:     check(hasNoErr(), hasAmount(1234)),
		},
		"invalid customer id": {
			customerID: "cus_missing",
			amount:     1234,
			checks:     check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			amount := 1234
			charge, err := c.Charge(tc.customerID, amount)
			for _, check := range tc.checks {
				check(t, charge, err)
			}
		})
	}

}
