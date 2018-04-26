package v3

import (
	"configcenter/src/common"
	"configcenter/src/common/http/httpclient"
)

// Client the http client
type Client struct {
	httpCli         *httpclient.HttpClient
	address         string
	supplierAccount string
	user            string
}

var client = &Client{}

func init() {

	client.httpCli = httpclient.NewHttpClient()
	client.httpCli.SetHeader("Content-Type", "application/json")
	client.httpCli.SetHeader("Accept", "application/json")
}

// GetV3Client get the v3 client
func GetV3Client() *Client {

	return client
}

// SetAddress set a new address
func (cli *Client) SetAddress(address string) {
	cli.address = address
}

// SetSupplierAccount set a new supplieraccount
func (cli *Client) SetSupplierAccount(supplierAccount string) {
	cli.supplierAccount = supplierAccount
	cli.httpCli.SetHeader(common.BKHTTPOwnerID, supplierAccount)
}

// SetUser set a new user
func (cli *Client) SetUser(user string) {
	cli.user = user
	cli.httpCli.SetHeader(common.BKHTTPHeaderUser, user)
}

// GetUser get the user
func (cli *Client) GetUser() string {
	return cli.user
}

// GetSupplierAccount get the supplier account
func (cli *Client) GetSupplierAccount() string {
	return cli.supplierAccount
}

// GetAddress get the address
func (cli *Client) GetAddress() string {
	return cli.address
}
