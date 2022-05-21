package intranet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HourEntry struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Time        float32 `json:"time"`
}

type Project struct {
	Client ProjectClient `json:"client"`
	Name   string        `json:"name"`
}

type ProjectClient struct {
	Name string `json:"name"`
}

type GetHourEntriesInput struct {
	Date time.Time
}

type GetHourEntriesOutput struct {
	Entries []HourEntry `json:"entries"`
}

func (c *Client) GetHourEntries(input *GetHourEntriesInput) (*GetHourEntriesOutput, error) {
	date := input.Date
	url := fmt.Sprintf("%s/intranet4/hours?date=%s", c.BaseURL, date.Format(DateFormat))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := GetHourEntriesOutput{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type CreateHourEntryInput struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	ProjectID   int     `json:"project_id"`
	TicketID    string  `json:"ticket_id"`
	Time        float32 `json:"time"`
}

type CreateHourEntryOutput struct {
	Added       string  `json:"added"`
	Date        string  `json:"date"`
	Description string  `json:"desc"`
	ID          int     `json:"id"`
	Modified    string  `json:"modified"`
	Project     Project `json:"project"`
	TicketID    string  `json:"ticketId"`
	Time        float32 `json:"time"`
	UserID      string  `json:"userID"`
}

func (c *Client) CreateHourEntry(input *CreateHourEntryInput) (*CreateHourEntryOutput, error) {
	url := fmt.Sprintf("%s/user_times", c.BaseURL)

	reqBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	res := CreateHourEntryOutput{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type DeleteHourEntryInput struct {
	ID int `json:"id"`
}

func (c *Client) DeleteHourEntry(input *DeleteHourEntryInput) error {
	url := fmt.Sprintf("%s/intranet4/user_times", c.BaseURL)

	reqBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequest("DELETE", url, reqBody)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
