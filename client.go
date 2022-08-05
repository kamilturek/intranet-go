package intranet

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BaseURL         string = "https://intranet.stxnext.pl/api"
	DateFormat      string = "2006-01-02"
	SessionIDEnvVar string = "INTRANET_SESSION_ID"
)

type Client struct {
	baseURL    string
	sessionID  string
	httpClient *http.Client
}

func NewClient(opts ...option) *Client {
	return &Client{
		baseURL:   BaseURL,
		sessionID: os.Getenv(SessionIDEnvVar),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

type option func(c *Client) error

func WithSessionID(sessionID string) option {
	return func(c *Client) error {
		if sessionID == "" {
			return errors.New("blank session ID")
		}

		c.sessionID = sessionID

		return nil
	}
}

func (c *Client) sendRequest(req *http.Request) (status int, data []byte, err error) {
	req.Header.Set("Cookie", fmt.Sprintf("beaker.session.id=%s", c.sessionID))

	res, err := c.httpClient.Do(req)
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
