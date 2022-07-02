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

type Entry struct {
	ID          int
	Description string
	Time        float64
	Project     struct {
		ClientName string
		ID         int
		Name       string
	}
}

type ListHourEntriesOutput struct {
	Entries []Entry `json:"entries"`
}

func (c *Client) ListHourEntries(input *ListHourEntriesInput) (*ListHourEntriesOutput, error) {
	url := fmt.Sprintf("%s/intranet4/hours?date=%s", c.BaseURL, input.Date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	status, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", status)
	}

	var output ListHourEntriesOutput

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type GetHourEntryInput struct {
	ID   int
	Date string
}

type GetHourEntryOutput Entry

func (c *Client) GetHourEntry(input *GetHourEntryInput) (*GetHourEntryOutput, error) {
	output, err := c.ListHourEntries(&ListHourEntriesInput{
		Date: input.Date,
	})
	if err != nil {
		return nil, err
	}

	for _, e := range output.Entries {
		if e.ID == input.ID {
			entry := GetHourEntryOutput(e)

			return &entry, nil
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

	postData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	status, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response status: %d", status)
	}

	var output CreateHourEntryOutput

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type DeleteHourEntryInput struct {
	ID string `json:"id"`
}

func (c *Client) DeleteHourEntry(input *DeleteHourEntryInput) error {
	url := fmt.Sprintf("%s/intranet4/user_times", c.BaseURL)

	postData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(postData))
	if err != nil {
		return err
	}

	status, _, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if status != http.StatusNoContent {
		return fmt.Errorf("unexpected response status: %d", status)
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

	postData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	status, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", status)
	}

	var output UpdateHourEntryOutput

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
