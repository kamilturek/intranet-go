package intranet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FlatProject struct {
	ClientName string `json:"clientName"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type HourEntry struct {
	ID          int         `json:"id"`
	Description string      `json:"description"`
	Time        float64     `json:"time"`
	Project     FlatProject `json:"project"`
}

type ListHourEntriesInput struct {
	Date string
}

type ListHourEntriesOutput struct {
	Entries []HourEntry `json:"entries"`
}

func (c *Client) ListHourEntries(input *ListHourEntriesInput) (*ListHourEntriesOutput, error) {
	url := fmt.Sprintf("%s/intranet4/hours?date=%s", c.BaseURL, input.Date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := ListHourEntriesOutput{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type GetHourEntryInput struct {
	ID   int
	Date string
}

type GetHourEntryOutput = HourEntry

func (c *Client) GetHourEntry(input *GetHourEntryInput) (*GetHourEntryOutput, error) {
	output, err := c.ListHourEntries(&ListHourEntriesInput{
		Date: input.Date,
	})
	if err != nil {
		return nil, err
	}

	for _, entry := range output.Entries {
		if entry.ID == input.ID {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("hour entry not found")
}

type Project struct {
	Client ProjectClient `json:"client"`
	Name   string        `json:"name"`
}

type ProjectClient struct {
	Name string `json:"name"`
}

type CreateHourEntryInput struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	ProjectID   int     `json:"projectId"`
	TicketID    string  `json:"ticketId"`
	Time        float64 `json:"time"`
}

type CreateHourEntryOutput struct {
	Added       string  `json:"added"`
	Date        string  `json:"date"`
	Description string  `json:"desc"`
	ID          string  `json:"id"`
	Modified    string  `json:"modified"`
	Project     Project `json:"project"`
	TicketID    string  `json:"ticketId"`
	Time        float64 `json:"time"`
	UserID      string  `json:"userID"`
}

func (c *Client) CreateHourEntry(input *CreateHourEntryInput) (*CreateHourEntryOutput, error) {
	url := fmt.Sprintf("%s/intranet4/user_times", c.BaseURL)

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
	ID string `json:"id"`
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

type UpdateHourEntryInput struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	ID          string  `json:"timeEntryId"`
	ProjectID   int     `json:"projectId"`
	TicketID    string  `json:"ticketId"`
	Time        float64 `json:"time"`
}

type UpdateHourEntryOutput = CreateHourEntryOutput

func (c *Client) UpdateHourEntry(input *UpdateHourEntryInput) (*UpdateHourEntryOutput, error) {
	url := fmt.Sprintf("%s/intranet4/user_times", c.BaseURL)

	reqBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequest("PUT", url, reqBody)
	if err != nil {
		return nil, err
	}

	res := UpdateHourEntryOutput{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
