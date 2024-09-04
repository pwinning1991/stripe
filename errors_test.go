package stripe_test

import (
	"encoding/json"
	"github.com/pwinning1991/stripe"
	"testing"
)

var errorJSON = []byte(`{
  "error": {
    "code": "resource_missing",
    "doc_url": "https://stripe.com/docs/error-codes/resource-missing",
    "message": "No such customer: cus_123",
    "param": "customer",
    "type": "invalid_request_error"
  }
}`)

func TestError_Unmarshal(t *testing.T) {
	var se stripe.Error
	err := json.Unmarshal(errorJSON, &se)
	if err != nil {
		t.Fatalf("Unmarshal expected error to be nil, got %s", err)
	}
	wantDocURL := "https://stripe.com/docs/error-codes/resource-missing"
	if se.DocURL != wantDocURL {
		t.Errorf("DOCURL = %s wamt %s", se.DocURL, wantDocURL)
	}
	wantType := "invalid_request_error"
	if se.Type != wantType {
		t.Errorf("Type = %s wamt %s", se.Type, wantType)
	}
	wantMessage := "No such customer: cus_123"
	if se.Message != wantMessage {
		t.Errorf("Message = %s wamt %s", se.Message, wantMessage)
	}
}

func TestError_Marshal(t *testing.T) {
	se := stripe.Error{
		Code:   "test-code",
		DocURL: "test-docURL",
		Type:   "test-type",
		Param:  "test-param",
	}
	data, err := json.Marshal(se)
	if err != nil {
		t.Fatalf("Marshal expected error to be nil, got %s", err)
	}
	var got stripe.Error
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Unmarshal expected error to be nil, got %s", err)
	}
	if got != se {
		t.Errorf("got = %v; want %v", got, se)
		t.Log("Is Unmarshal working? It is required for this test to pass.")

	}
}
