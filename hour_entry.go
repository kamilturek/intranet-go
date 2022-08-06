package intranet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Entry struct {
	ID          string
	Description string
	Time        float64
	Ticket      struct {
		ID string
	}
}

type ListHourEntriesInput struct {
	Date string
}

type listEntriesResponse struct {
	Entries []struct {
		ID          int
		Description string
		Time        float64
		Project     struct {
			ClientName string
			ID         int
			Name       string
		}
		Ticket struct {
			ID string
		}
	}
}

func (c *Client) ListHourEntries(input *ListHourEntriesInput) ([]*Entry, error) {
	url := fmt.Sprintf("%s/intranet4/hours?date=%s", c.baseURL, input.Date)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var output listEntriesResponse

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	entries := []*Entry{}

	for _, rawEntry := range output.Entries {
		entry := &Entry{
			ID:          strconv.Itoa(rawEntry.ID),
			Description: rawEntry.Description,
			Time:        rawEntry.Time,
			Ticket: struct{ ID string }{
				ID: rawEntry.Ticket.ID,
			},
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

type GetHourEntryInput struct {
	ID   string
	Date string
}

func (c *Client) GetHourEntry(input *GetHourEntryInput) (*Entry, error) {
	entries, err := c.ListHourEntries(&ListHourEntriesInput{
		Date: input.Date,
	})
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.ID == input.ID {
			return e, nil
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

type createHourEntryResponse struct {
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

func (c *Client) CreateHourEntry(input *CreateHourEntryInput) (*Entry, error) {
	url := fmt.Sprintf("%s/intranet4/user_times", c.baseURL)

	postData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postData))
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

	var output createHourEntryResponse

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return c.GetHourEntry(&GetHourEntryInput{
		ID:   output.ID,
		Date: output.Date,
	})
}

type DeleteHourEntryInput struct {
	ID string `json:"id"`
}

func (c *Client) DeleteHourEntry(input *DeleteHourEntryInput) error {
	url := fmt.Sprintf("%s/intranet4/user_times", c.baseURL)

	postData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(postData))
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

type updateHourEntryResponse createHourEntryResponse

func (c *Client) UpdateHourEntry(input *UpdateHourEntryInput) (*Entry, error) {
	url := fmt.Sprintf("%s/intranet4/user_times", c.baseURL)

	postData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(postData))
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

	var output updateHourEntryResponse

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}

	return c.GetHourEntry(&GetHourEntryInput{
		ID:   output.ID,
		Date: output.Date,
	})
}
