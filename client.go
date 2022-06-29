package intranet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL                string = "https://intranet.stxnext.pl/api"
	DateFormat             string = "2006-01-02"
	DateFormatAlternative1 string = "02-01-2006"
	DateFormatAlternative2 string = "02.01.2006"
	SessionIDEnvVar        string = "INTRANET_SESSION_ID"
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
		},
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Cookie", fmt.Sprintf("beaker.session.id=%s", c.SessionID))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		resBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		resBody := string(resBytes)
		return fmt.Errorf("HTTP error. Status: %d. Body: %s", res.StatusCode, resBody)
	}

	if v == nil {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func fail_over_here() {

}
