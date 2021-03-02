package client

import (
	"accountapi-client/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Client struct {
	httpClient *http.Client
	domain     string
}

type createRequest struct {
	Data struct {
		DataType       string `json:"type"`
		ID             string `json:"id"`
		OrganisationID string `json:"organisation_id"`
		Attributes     struct {
			Country       string `json:"country"`
			BaseCurrency  string `json:"base_currency,omitempty"`
			BankID        string `json:"bank_id"`
			Bic           string `json:"bic,omitempty"`
			BankIDCode    string `json:"bank_id_code,omitempty"`
			AccountNumber string `json:"account_number,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

type Response struct {
	Account      Account
	ErrorMessage string `json:"error_message"`
}

type Account struct {
	AccountID               string
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency"`
	BankID                  string   `json:"bank_id"`
	BankIDCode              string   `json:"bank_id_code"`
	AccountNumber           string   `json:"account_number"`
	Bic                     string   `json:"bic"`
	Iban                    string   `json:"iban"`
	CustomerID              string   `json:"customer_id"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_name"`
	AccountCassification    string   `json:"account_classification"`
	JointAccount            bool     `json:"joint_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
	Switched                bool     `json:"swtiched"`
	Status                  string   `json:"status"`
}

//ClientConfig is the default config file, can be modified to use a different file
var ClientConfig string = "config/server.yaml"

//Create is the method to send a create request to the server
func (cli *Client) Create(request map[string]string) (statusCode int, response Response, err error) {
	createRequest := buildRequest(request)
	requestString, err := json.Marshal(createRequest)
	resp, err := cli.httpClient.Post(cli.domain+"/v1/organisation/accounts", "application/json", bytes.NewBuffer(requestString))
	if err != nil {
		return statusCode, response, err
	}
	response = parseResponse(resp)

	return resp.StatusCode, response, err
}

func (cli *Client) Fetch(id string) (statusCode int, response Response, err error) {
	resp, err := cli.httpClient.Get(cli.domain + "/v1/organisation/accounts/" + id)
	if err != nil {
		return statusCode, response, err
	}

	response = parseResponse(resp)

	return resp.StatusCode, response, err
}

func (cli *Client) Delete(id string, version int) (statusCode int, response Response, err error) {
	url := cli.domain + "/v1/organisation/accounts/" + id
	if version != 0 {
		url = url + "?version=" + strconv.Itoa(version)
	}
	request, err := http.NewRequest("DELETE", url, nil)
	resp, err := cli.httpClient.Do(request)
	if err != nil {
		return statusCode, response, err
	}
	response = parseResponse(resp)
	return resp.StatusCode, response, err
}

func NewClient() *Client {
	client := Client{}
	c := config.Schema{}
	config.ReadConfig(ClientConfig, &c)
	addr := c.Server.Host + ":" + c.Server.Port
	client.domain = addr
	client.httpClient = http.DefaultClient
	return &client
}

func buildRequest(request map[string]string) (createRequest createRequest) {
	if _, ok := request["AccountID"]; ok == false {
		request["AccountID"] = uuid.New().String()
	}

	createRequest.Data.DataType = "accounts"
	createRequest.Data.ID = request["AccountID"]
	createRequest.Data.OrganisationID = request["OrganisationID"]
	createRequest.Data.Attributes.AccountNumber = request["AccountNumber"]
	createRequest.Data.Attributes.Country = request["Country"]
	createRequest.Data.Attributes.BaseCurrency = request["BaseCurrency"]
	createRequest.Data.Attributes.BankID = request["BankID"]
	createRequest.Data.Attributes.BankIDCode = request["BankIDCode"]
	createRequest.Data.Attributes.Bic = request["Bic"]
	return
}

func parseResponse(resp *http.Response) (response Response) {
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return
	}
	if resp.StatusCode == 200 {
		result := map[string]map[string]Account{}
		json.Unmarshal(body, &result)
		getid := map[string]map[string]string{}
		json.Unmarshal(body, &getid)
		response.Account = result["data"]["attributes"]
		response.Account.AccountID = getid["data"]["id"]
	} else {
		json.Unmarshal(body, &response)
	}
	return
}
