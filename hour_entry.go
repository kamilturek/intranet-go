package intranet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ListHourEntriesInput struct {
	Date string
}

type entry struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Time        float64 `json:"time"`
	Project     struct {
		ClientName string `json:"clientName"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
	} `json:"project"`
}

type ListHourEntriesOutput struct {
	Entries []entry `json:"entries"`
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

type GetHourEntryOutput entry

func (c *Client) GetHourEntry(input *GetHourEntryInput) (*GetHourEntryOutput, error) {
	output, err := c.ListHourEntries(&ListHourEntriesInput{
		Date: input.Date,
	})
	if err != nil {
		return nil, err
	}

	for _, entry := range output.Entries {
		if entry.ID == input.ID {
			entryOutput := GetHourEntryOutput(entry)
			return &entryOutput, nil
		}
	}

	return nil, fmt.Errorf("hour entry not found")
}

type CreateHourEntryInput struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	ProjectID   int     `json:"projectId"`
	TicketID    string  `json:"ticketId"`
	Time        float64 `json:"time"`
}

type CreateHourEntryOutput struct {
	Added       string
	Date        string
	Description string `json:"desc"`
	ID          string
	Modified    string
	Project     struct {
		Client struct {
			Name string
		}
		Name string
	}
	TicketID string
	Time     float64
	UserID   string
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

type UpdateHourEntryOutput CreateHourEntryOutput

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
