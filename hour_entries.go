package intranet

import (
	"fmt"
	"net/http"
	"time"
)

type HourEntriesList struct {
	Entries []HourEntry `json:"entries"`
}

type HourEntry struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Time        float64 `json:"time"`
}

type GetHourEntriesOptions struct {
	Date time.Time
}

func (c *Client) GetHourEntries(options *GetHourEntriesOptions) (*HourEntriesList, error) {
	date := time.Now()

	if options != nil {
		date = options.Date
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/intranet4/hours?date=%s", c.BaseURL, date.Format(DateFormat)), nil)
	if err != nil {
		return nil, err
	}

	res := HourEntriesList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
