package intranet_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/google/go-cmp/cmp"
	"github.com/kamilturek/intranet-go"
)

func GetClient(t *testing.T, cassetteName string) (*intranet.Client, func()) {
	r, err := recorder.New(fmt.Sprintf("fixtures/%s", cassetteName))
	if err != nil {
		t.Fatal(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Cookie")
		return nil
	})

	sessionID := os.Getenv(intranet.SessionIDEnvVar)
	client := intranet.NewClient(sessionID)
	client.HTTPClient.Transport = r

	deferFunc := func() {
		r.Stop()
	}

	return client, deferFunc
}

func TestCreateHourEntry(t *testing.T) {
	client, deferFunc := GetClient(t, "create")
	defer deferFunc()

	got, err := client.CreateHourEntry(&intranet.CreateHourEntryInput{
		Date:        "2022-07-02",
		Description: "Working on feature A",
		ProjectID:   123,
		TicketID:    "ABC123",
		Time:        0.5,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteHourEntry(&intranet.DeleteHourEntryInput{
		ID: got.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &intranet.CreateHourEntryOutput{
		Added:       "2022-07-02",
		Date:        "2022-07-02",
		Description: "Working on feature A",
		ID:          "2178010",
		Modified:    "2022-07-02",
		Project: struct {
			Client struct {
				Name string
			}
			Name string
		}{
			Client: struct {
				Name string
			}{
				Name: "Test Client",
			},
			Name: "Test Project",
		},
		TicketID: "ABC123",
		Time:     0.5,
		UserID:   "7777",
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestListHourEntries(t *testing.T) {
	client, deferFunc := GetClient(t, "list")
	defer deferFunc()

	got, err := client.ListHourEntries(&intranet.ListHourEntriesInput{
		Date: "2022-07-02",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &intranet.ListHourEntriesOutput{
		Entries: []intranet.Entry{
			{
				ID:          2177998,
				Description: "Test",
				Time:        0.25,
				Project: struct {
					ClientName string
					ID         int
					Name       string
				}{
					ClientName: "Test Client",
					ID:         123,
					Name:       "Test Project",
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestGetHourEntry(t *testing.T) {
	client, deferFunc := GetClient(t, "get")
	defer deferFunc()

	got, err := client.GetHourEntry(&intranet.GetHourEntryInput{
		ID:   2177998,
		Date: "2022-07-02",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &intranet.GetHourEntryOutput{
		ID:          2177998,
		Description: "Test",
		Time:        0.25,
		Project: struct {
			ClientName string
			ID         int
			Name       string
		}{
			ClientName: "Test Client",
			ID:         123,
			Name:       "Test Project",
		},
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestGetHourEntryNotFound(t *testing.T) {
	client, deferFunc := GetClient(t, "get_not_found")
	defer deferFunc()

	_, err := client.GetHourEntry(&intranet.GetHourEntryInput{
		ID:   1,
		Date: "2022-07-03",
	})
	if err == nil {
		t.Fatal("want error, got nil")
	}

	want := "hour entry not found"
	got := err.Error()

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestUpdateHourEntry(t *testing.T) {
	client, deferFunc := GetClient(t, "update")
	defer deferFunc()

	created, err := client.CreateHourEntry(&intranet.CreateHourEntryInput{
		Date:        "2022-07-01",
		Description: "Working on feature A",
		ProjectID:   123,
		TicketID:    "ABC123",
		Time:        0.5,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.UpdateHourEntry(&intranet.UpdateHourEntryInput{
		Date:        "2022-07-02",
		Description: "Working on feature B",
		ID:          created.ID,
		ProjectID:   456,
		TicketID:    "CDE456",
		Time:        1,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteHourEntry(&intranet.DeleteHourEntryInput{
		ID: created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &intranet.UpdateHourEntryOutput{
		Added:       "2022-07-02",
		Date:        "2022-07-02",
		Description: "Working on feature B",
		ID:          "2178009",
		Modified:    "2022-07-02",
		Project: struct {
			Client struct{ Name string }
			Name   string
		}{
			Client: struct{ Name string }{
				Name: "Test Client",
			},
			Name: "Test Project",
		},
		TicketID: "CDE456",
		Time:     1,
		UserID:   "7777",
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
