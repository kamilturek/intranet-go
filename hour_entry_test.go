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
		Date:        "2022-07-01",
		Description: "Working on feature A",
		ProjectID:   123,
		TicketID:    "ABC123",
		Time:        0.5,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &intranet.CreateHourEntryOutput{
		Added:       "2022-07-01",
		Date:        "2022-07-01",
		Description: "Working on feature A",
		ID:          "2177996",
		Modified:    "2022-07-01",
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
