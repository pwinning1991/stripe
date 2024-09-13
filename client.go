package stripe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	Version         = "2018-09-24"
	DefaultCurrency = "usd"
	DefaultBaseURL  = "https://api.stripe.com/v1"
)

type Charge struct {
	ID             string `json:"id"`
	Amount         int    `json:"amount"`
	FailureCode    string `json:"failure_code"`
	FailureMessage string `json:"failure_message"`
	Paid           bool   `json:"paid"`
	Status         string `json:"status"`
}

type Customer struct {
	ID            string `json:"id"`
	DefaultSource string `json:"default_source"`
	Email         string `json:"email"`
}

type Client struct {
	Key        string
	BaseURL    string
	HttpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	httpClient := c.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if req.Method != http.MethodGet {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Stripe-Version", Version)
	req.SetBasicAuth(c.Key, "")
	return httpClient.Do(req)
}

func (c *Client) url(path string) string {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	return fmt.Sprintf("%s/%s", c.BaseURL, path)
}

func (c *Client) Customer(token string, email string) (*Customer, error) {
	endpoint := c.url("customers")
	v := url.Values{}
	v.Set("source", token)
	v.Set("email", email)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Stripe-Version", Version)
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, parseError(body)
	}
	//fmt.Println(string(body))
	var cus Customer
	err = json.Unmarshal(body, &cus)
	if err != nil {
		return nil, err
	}
	return &cus, nil
}

func (c *Client) Charge(customerId string, amount int) (*Charge, error) {
	endpoint := c.url("charges")
	v := url.Values{}
	v.Set("customer", customerId)
	v.Set("amount", strconv.Itoa(amount))
	v.Set("currency", DefaultCurrency)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, parseError(body)
	}
	//fmt.Println(string(body))
	var chg Charge
	err = json.Unmarshal(body, &chg)
	if err != nil {
		return nil, err
	}
	return &chg, nil
}

func (c *Client) GetCharge(chargeID string) (*Charge, error) {
	endpoint := c.url(fmt.Sprintf("/charges/%s", chargeID))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, parseError(body)
	}

	var chg Charge
	err = json.Unmarshal(body, &chg)
	if err != nil {
		return nil, err
	}
	return &chg, nil
}

func parseError(data []byte) error {
	var se Error
	err := json.Unmarshal(data, &se)
	if err != nil {
		return err
	}
	return se
}
