package eslib

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultProtocol = "http"
	DefaultDomain   = "localhost"
	DefaultPort     = "9200"
)

type EsClient interface {
	String() string
	SetFromFlag(rawUrl string) error
	NewRequest(method string, path string, query string) (*http.Request, error)
}

type EsClientImpl struct {
	Protocol string
	Domain   string
	Port     string
	Username string
	Password string
}

func NewEsClient() EsClientImpl {
	return EsClientImpl{
		Protocol: DefaultProtocol,
		Domain:   DefaultDomain,
		Port:     DefaultPort,
	}
}

func (esClient *EsClientImpl) String() string {
	return fmt.Sprintf("Protocol:%s - Domain:%s - Port: %s - Username: %s - Password: %s", esClient.Protocol, esClient.Domain, esClient.Port, esClient.Username, esClient.Password)
}

func (esClient *EsClientImpl) SetFromFlag(rawUrl string) error {
	if rawUrl == "" {
		return errors.New("Url is empty")
	}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme != "" {
		esClient.Protocol = parsedUrl.Scheme
	}
	esClient.Domain, esClient.Port = splitHostnamePartsFromHost(parsedUrl.Host)

	if parsedUrl.User != nil {
		esClient.Username = parsedUrl.User.Username()
		password, passwordIsSet := parsedUrl.User.Password()
		if passwordIsSet {
			esClient.Password = password
		}
	}

	return nil
}

func (esClient *EsClientImpl) NewRequest(method string, path string, query string) (*http.Request, error) {
	var uri string

	path = strings.Trim(path, "/")

	// If query parameters are provided, the add them to the URL
	if len(query) > 0 {
		uri = fmt.Sprintf("%s://%s:%s/%s?%s", esClient.Protocol, esClient.Domain, esClient.Port, path, query)
	} else {
		uri = fmt.Sprintf("%s://%s:%s/%s", esClient.Protocol, esClient.Domain, esClient.Port, path)
	}

	request, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")

	if esClient.Username != "" || esClient.Password != "" {
		request.SetBasicAuth(esClient.Username, esClient.Password)
	}

	return request, err
}

// Split apart the hostname on colon
// Return the host and a default port if there is no separator
func splitHostnamePartsFromHost(fullHost string) (domain string, port string) {
	splittedHost := strings.Split(fullHost, ":")

	if len(splittedHost) == 2 {
		return splittedHost[0], splittedHost[1]
	}
	return splittedHost[0], DefaultPort
}
