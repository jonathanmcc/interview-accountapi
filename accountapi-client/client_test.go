package client

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingHttpClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

const (
	okResponse = `{
		"data": {
			"attributes": {
				"alternative_bank_account_names": null,
				"bank_id": "400300",
				"bank_id_code": "GBDSC",
				"base_currency": "GBP",
				"bic": "NWBKGB22",
				"country": "GB"
			},
			"created_on": "2021-02-16T17:34:43.404Z",
			"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4de",
			"modified_on": "2021-02-16T17:34:43.404Z",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"type": "accounts",
			"version": 0
		},
		"links": {
			"self": "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4de"
		}
	}`
	errResponse = `{
		"error_message": "Account cannot be created as it violates a duplicate constraint"
	}`
	notFoundResponse = `{
		"error_message": "record ad27e265-9605-4b4b-a0e5-3003ea9cc4de does not exist"
	}`
)

func TestNewClient(t *testing.T) {
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	assert.Equal(t, "http://localhost:8080", cli.domain)
}

// buildRequest should generate account id if it does not exist
func TestBuildRequest(t *testing.T) {
	request := map[string]string{
		"Country":      "GB",
		"BaseCurrency": "GBP",
		"BankID":       "400300",
		"Bic":          "NWBKGB22",
		"BankIDCode":   "GBDSC",
	}
	createRequest := buildRequest(request)
	assert.NotEqual(t, createRequest.Data.ID, "")
	assert.Equal(t, createRequest.Data.OrganisationID, "") //should not generate organisation id
}

func TestCreate(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts", r.URL.Path)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, string(body), `{"data":{"type":"accounts","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4de","organisation_id":"aede1368-00d3-45d9-b3cb-1907f6ae9d4d","attributes":{"country":"GB","base_currency":"GBP","bank_id":"400300","bic":"NWBKGB22","bank_id_code":"GBDSC"}}}`)
		w.Write([]byte(okResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	request := map[string]string{
		"AccountID":      "ad27e265-9605-4b4b-a0e5-3003ea9cc4de",
		"OrganisationID": "aede1368-00d3-45d9-b3cb-1907f6ae9d4d",
		"Country":        "GB",
		"BaseCurrency":   "GBP",
		"BankID":         "400300",
		"Bic":            "NWBKGB22",
		"BankIDCode":     "GBDSC",
	}
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	status, resp, err := cli.Create(request)
	expectedResponse := Account{}
	expectedResponse.AccountID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4de"
	expectedResponse.BankID = "400300"
	expectedResponse.BankIDCode = "GBDSC"
	expectedResponse.BaseCurrency = "GBP"
	expectedResponse.Bic = "NWBKGB22"
	expectedResponse.Country = "GB"
	assert.Equal(t, expectedResponse, resp.Account)
	assert.Equal(t, status, 200)
	assert.NoError(t, err)
}

func TestCreateErr(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts", r.URL.Path)
		w.WriteHeader(409)
		w.Write([]byte(errResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	request := map[string]string{
		"AccountID":    "ad27e265-9605-4b4b-a0e5-3003ea9cc4de",
		"Country":      "GB",
		"BaseCurrency": "GBP",
		"BankID":       "400300",
		"Bic":          "NWBKGB22",
		"BankIDCode":   "GBDSC",
	}
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	statusCode, resp, _ := cli.Create(request)
	assert.Equal(t, statusCode, 409)
	assert.Equal(t, resp.ErrorMessage, "Account cannot be created as it violates a duplicate constraint")
	// fmt.Println(resp)
	// fmt.Println(err)
}

func TestFetch(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4de", r.URL.Path)
		w.WriteHeader(200)
		w.Write([]byte(okResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	statusCode, resp, err := cli.Fetch("ad27e265-9605-4b4b-a0e5-3003ea9cc4de")
	expectedResponse := Account{}
	expectedResponse.AccountID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4de"
	expectedResponse.BankID = "400300"
	expectedResponse.BankIDCode = "GBDSC"
	expectedResponse.BaseCurrency = "GBP"
	expectedResponse.Bic = "NWBKGB22"
	expectedResponse.Country = "GB"
	assert.Equal(t, expectedResponse, resp.Account)
	assert.Equal(t, statusCode, 200)
	assert.NoError(t, err)
}

func TestFetchNotFound(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4de", r.URL.Path)
		w.WriteHeader(404)
		w.Write([]byte(notFoundResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	statusCode, resp, err := cli.Fetch("ad27e265-9605-4b4b-a0e5-3003ea9cc4de")
	assert.Equal(t, resp.ErrorMessage, "record ad27e265-9605-4b4b-a0e5-3003ea9cc4de does not exist")
	assert.Equal(t, statusCode, 404)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4de", r.URL.Path)
		w.WriteHeader(204)
		w.Write([]byte(okResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	statusCode, _, err := cli.Delete("ad27e265-9605-4b4b-a0e5-3003ea9cc4de", 0)
	assert.Equal(t, statusCode, 204)
	assert.NoError(t, err)
}

func TestDeleteNotFound(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "localhost:8080", r.Host)
		assert.Equal(t, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4de", r.URL.Path)
		w.WriteHeader(404)
		w.Write([]byte(notFoundResponse))
	})
	ClientConfig = "config/examples/server.yaml"
	cli := NewClient()
	httpClient, teardown := testingHttpClient(h)
	cli.httpClient = httpClient
	defer teardown()
	statusCode, resp, err := cli.Delete("ad27e265-9605-4b4b-a0e5-3003ea9cc4de", 0)
	assert.Equal(t, resp.ErrorMessage, "record ad27e265-9605-4b4b-a0e5-3003ea9cc4de does not exist")
	assert.Equal(t, statusCode, 404)
	assert.NoError(t, err)
}
