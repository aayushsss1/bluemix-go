package k8sclusterv1

import (
	"fmt"
	gohttp "net/http"
	"path"
	"strings"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/trace"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

//AuthorizationHeader ...
const AuthorizationHeader = "Authorization"

//Client is the Aramda K8s client ...
type Client interface {
	Clusters() Clusters
	Workers() Workers
	WebHooks() Webhooks
	Subnets() Subnets
}

//URLGetter ...
type URLGetter func() string

//ErrHandler ...
type ErrHandler func(statusCode int, rawResponse []byte) error

//BeforeHandler ...
type BeforeHandler func(*rest.Request) error

//TokenRefresher ...
type TokenRefresher interface {
	RefreshToken() (string, error)
}

type clusterClient struct {
	IAMTokenRefresher TokenRefresher
	BaseURL           URLGetter
	OnError           ErrHandler
	Before            BeforeHandler

	config     *bluemix.Config
	HTTPClient *gohttp.Client
}

//NewClient ...
func NewClient(s *session.Session) (Client, error) {
	config := s.Config.Copy()

	_, err := config.EndpointLocator.ContainerEndpoint()
	if err != nil {
		return nil, err
	}
	baseURL := func() string {
		ep, _ := config.EndpointLocator.ContainerEndpoint()
		return ep
	}

	httpClient := http.NewHTTPClient(config)

	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		HTTPClient: httpClient,
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})

	if err != nil {
		return nil, err
	}
	client := &clusterClient{
		BaseURL:           baseURL,
		IAMTokenRefresher: tokenRefreher,
		config:            config,
		HTTPClient:        httpClient,
	}
	return client, nil
}

//Clusters implements Clusters API
func (c *clusterClient) Clusters() Clusters {
	return newClusterAPI(c)
}

//Workers implements Cluster Workers API
func (c *clusterClient) Workers() Workers {
	return newWorkerAPI(c)
}

//Subnets implements Cluster Subnets API
func (c *clusterClient) Subnets() Subnets {
	return newSubnetAPI(c)
}

//Webhooks implements Cluster WebHooks API
func (c *clusterClient) WebHooks() Webhooks {
	return newWebhookAPI(c)
}

func (c *clusterClient) sendRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = gohttp.DefaultClient
	}

	restClient := &rest.Client{
		DefaultHeader: http.DefaultClusterAuthHeader(c.config),
		HTTPClient:    httpClient,
	}

	if c.Before != nil {
		err := c.Before(r)
		if err != nil {
			return new(gohttp.Response), err
		}
	}

	resp, err := restClient.Do(r, respV, nil)

	// The response returned by go HTTP client.Do() could be nil if request timeout.
	// For convenience, we ensure that response returned by this method is always not nil.
	if resp == nil {
		return new(gohttp.Response), err
	}

	if err != nil {
		err = bmxerror.WrapNetworkErrors(resp.Request.URL.Host, err)
	}

	// if token is invalid, refresh and try again
	if resp.StatusCode == 401 && c.IAMTokenRefresher != nil {
		trace.Logger.Println("Authentication token probably expired, attempting refresh ...")
		_, uaaErr := c.IAMTokenRefresher.RefreshToken()
		switch uaaErr.(type) {
		case nil:
			restClient.DefaultHeader = http.DefaultClusterAuthHeader(c.config)
			resp, err = restClient.Do(r, respV, nil)
		case *bmxerror.InvalidTokenError:
			return resp, bmxerror.NewRequestFailure("InvalidToken", fmt.Sprintf("%v", uaaErr), 401)
		default:
			return resp, fmt.Errorf("Authentication failed, Unable to refresh auth token: %v. Try again later", uaaErr)
		}
	}

	if errResp, ok := err.(bmxerror.RequestFailure); ok && c.OnError != nil {
		err = c.OnError(errResp.StatusCode(), []byte(errResp.Description()))
	}

	return resp, err
}

func (c *clusterClient) get(path string, respV interface{}, targetHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.GetRequest(c.url(path))
	for _, t := range targetHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

func (c *clusterClient) put(path string, data interface{}, respV interface{}, targetHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PutRequest(c.url(path)).Body(data)
	for _, t := range targetHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

func (c *clusterClient) patch(path string, data interface{}, respV interface{}, targetHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PatchRequest(c.url(path)).Body(data)
	for _, t := range targetHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

func (c *clusterClient) post(path string, data interface{}, respV interface{}, targetHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PostRequest(c.url(path)).Body(data)
	for _, t := range targetHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

func (c *clusterClient) delete(path string, targetHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.url(path))
	for _, t := range targetHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, nil)
}

func (c *clusterClient) url(path string) string {
	if c.BaseURL == nil {
		return path
	}

	return c.BaseURL() + cleanPath(path)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

const (
	orgIDHeader     = "X-Auth-Resource-Org"
	spaceIDHeader   = "X-Auth-Resource-Space"
	accountIDHeader = "X-Auth-Resource-Account"

	slUserNameHeader = "X-Auth-Softlayer-Username"
	slAPIKeyHeader   = "X-Auth-Softlayer-APIKey"
)

func addToRequestHeader(h interface{}, r *rest.Request) {
	switch v := h.(type) {
	case *ClusterTargetHeader:
		r.Set(orgIDHeader, v.OrgID)
		r.Set(spaceIDHeader, v.SpaceID)
		r.Set(accountIDHeader, v.AccountID)

	case *ClusterSoftlayerHeader:
		r.Set(slUserNameHeader, v.SoftLayerUsername)
		r.Set(slAPIKeyHeader, v.SoftLayerAPIKey)
	}

}
