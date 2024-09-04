package stripe_test

import (
	"flag"
	"github.com/pwinning1991/stripe"
	"strings"
	"testing"
)

var (
	apiKey string
)

//sk_test_4eC39HqLyjWDarjtT1zdp7dc

func init() {
	flag.StringVar(&apiKey, "key", "", "Your TEST secret key for the Stripe API."+
		"If present, integartion tests will be run using this API key.")
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No TEST API key provided.")
	}
	c := stripe.Client{
		Key: apiKey,
	}
	tok := "tok_amex"
	email := "test@test.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Errorf("Customer() err = %v; want %v", err, nil)
	}
	if cus == nil {
		t.Fatalf("Customer() cus = nil; want non nil value")
	}

	if !strings.HasPrefix(cus.ID, "cus_") {
		t.Errorf("Customer() CustomerID = %s; want prefix %q", cus.ID, "cus_")
	}
	if !strings.HasPrefix(cus.DefaultSource, "card_") {
		t.Errorf("Customer() DefaultSource = %s; want prefix %q", cus.DefaultSource, "cus_")
	}
	if cus.Email != email {
		t.Errorf("Customer() CustomerEmail = %s; want %s", cus.Email, email)
	}
}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No TEST API key provided.")
	}
	c := stripe.Client{
		Key: apiKey,
	}
	//create customer for the test
	tok := "tok_amex"
	email := "test@test.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Errorf("Customer() err = %v; want %v", err, nil)
	}
	amount := 1234
	charge, err := c.Charge(cus.ID, amount)
	if err != nil {
		t.Errorf("Charge() err = %v; want %v", err, nil)
	}
	if charge == nil {
		t.Fatalf("Charge() = nil; want non nil value")
	}
	if charge.Amount != amount {
		t.Errorf("Charge() Amount = %d; want %d", charge.Amount, amount)
	}
}
