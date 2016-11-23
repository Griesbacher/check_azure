package azureHttp

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest/azure"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	//ApiURL the base url of azure
	ApiURL = "https://management.azure.com/subscriptions/"
	//ContentTypeJSON the result is json
	ContentTypeJSON = "application/json; charset=utf-8"
	//ContentTypeXML the result is xml
	ContentTypeXML = "application/xml; charset=utf-8"
)

func ContentError(expexted, got string) error {
	return fmt.Errorf("The expected content-type(%s) has not been returned: %s", expexted, got)
}

type AzureConnector struct {
	clientId       string
	clientSecret   string
	subscriptionId string
	tenantId       string
	azureClient    storage.AccountsClient
}

func NewAzureConnector(clientId, clientSecret, subscriptionId, tenantId string) (*AzureConnector, error) {
	c := &AzureConnector{
		clientId:       clientId,
		clientSecret:   clientSecret,
		subscriptionId: subscriptionId,
		tenantId:       tenantId,
	}
	err := c.authorize()
	//TODO: test given secrets
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a *AzureConnector) authorize() error {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(a.tenantId)
	if err != nil {
		return err
	}
	spt, err := azure.NewServicePrincipalToken(*oauthConfig, a.clientId, a.clientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return err
	}
	ac := storage.NewAccountsClient(a.subscriptionId)
	ac.Authorizer = spt
	a.azureClient = ac
	return nil
}

func (a AzureConnector) HttpRequest(u url.URL) ([]byte, string, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := a.azureClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", errors.New("Returncode is not 2XX.")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if typ, ok := resp.Header["Content-Type"]; ok {
		return body, typ[0], nil
	}
	return body, "", nil
}
func (a AzureConnector) RequestWithSub(apiVersion, path, filter string) ([]byte, string, error) {
	return a.Request(apiVersion, fmt.Sprintf("%s/%s/", a.subscriptionId, path), filter)
}
func (a AzureConnector) Request(apiVersion, path, filter string) ([]byte, string, error) {
	u, err := url.Parse(ApiURL)
	if err != nil {
		return nil, "", err
	}
	u.Path += path
	params := url.Values{}
	params.Add("api-version", apiVersion)
	if filter != "" {
		params.Add("$filter", filter)
	}
	u.RawQuery = params.Encode()
	//fmt.Println(u.String())
	body, typ, err := a.HttpRequest(*u)
	if err != nil {
		return nil, "", err
	}
	return body, typ, nil
}
