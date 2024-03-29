package monday

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httputilmore"
)

const (
	MondayAPIURL = "https://api.monday.com/v2"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(token string) Client {
	return Client{httpClient: authutil.NewClientAuthzTokenSimple("", token)}
}

func (c *Client) DoJSON(data []byte) (*http.Response, error) {
	if c.httpClient == nil {
		return nil, errors.New("no auth token")
	}
	req, err := http.NewRequest(
		http.MethodPost,
		MondayAPIURL,
		bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)
	return c.httpClient.Do(req)
}

func (c *Client) DoGraphQLString(gql string) (*http.Response, error) {
	req := QueryRequest{Query: gql}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.DoJSON(data)
}

func (c *Client) DoGraphQL(gql Query) (*http.Response, error) {
	return c.DoGraphQLString(gql.String())
}

type QueryRequest struct {
	Query string `json:"query"`
}
