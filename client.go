package intranet

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL         string = "https://intranet.stxnext.pl/api"
	DateFormat      string = "2006-01-02"
	SessionIDEnvVar string = "INTRANET_SESSION_ID"
)

type Client struct {
	BaseURL    string
	SessionID  string
	HTTPClient *http.Client
}

func NewClient(sessionID string) *Client {
	return &Client{
		BaseURL:   BaseURL,
		SessionID: sessionID,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) sendRequest(req *http.Request) (status int, data []byte, err error) {
	req.Header.Set("Cookie", fmt.Sprintf("beaker.session.id=%s", c.SessionID))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return status, data, fmt.Errorf("error making request: %w", err)
	}

	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return status, data, fmt.Errorf("error reading response body: %w", err)
	}

	return res.StatusCode, data, nil
}
