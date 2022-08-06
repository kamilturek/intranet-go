//go:build integration

package intranet_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kamilturek/intranet-go"
)

var now = time.Now()

func TestAuthenticationError(t *testing.T) {
	client, err := intranet.NewClient(
		intranet.WithSessionID("foo"),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ListHourEntries(&intranet.ListHourEntriesInput{
		Date: intranet.Date(now),
	})
	if err == nil {
		t.Fatalf("want error, got nil")
	}

	want := "unexpected response status: 302"
	got := err.Error()
	if want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestHourEntry(t *testing.T) {
	if _, set := os.LookupEnv(intranet.SessionIDEnvVar); !set {
		t.Fatalf("%s environment variable must be set", intranet.SessionIDEnvVar)
	}

	ignoreFields := cmpopts.IgnoreFields(intranet.Entry{}, "ID", "Project.Name", "Project.Client.Name")

	client, err := intranet.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// Create
	want := &intranet.Entry{
		ID:          "",
		Date:        intranet.Date(now),
		Description: "Test",
		Project: struct {
			ID     int
			Name   string
			Client struct{ Name string }
		}{
			ID: 422,
		},
		Time: 1.5,
		Ticket: struct{ ID string }{
			ID: "TEST",
		},
	}
	got, err := client.CreateHourEntry(
		&intranet.CreateHourEntryInput{
			Date:        intranet.Date(now),
			Description: "Test",
			ProjectID:   422,
			TicketID:    "TEST",
			Time:        1.5,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, ignoreFields); diff != "" {
		t.Error(diff)
	}

	// Get
	got, err = client.GetHourEntry(
		&intranet.GetHourEntryInput{
			ID:   got.ID,
			Date: intranet.Date(now),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, ignoreFields); diff != "" {
		t.Error(diff)
	}

	// Update
	want = &intranet.Entry{
		ID:          "",
		Date:        intranet.Date(now),
		Description: "Test Updated",
		Project: struct {
			ID     int
			Name   string
			Client struct{ Name string }
		}{
			ID: 422,
		},
		Time: 2.5,
		Ticket: struct{ ID string }{
			ID: "TEST-UPDATED",
		},
	}
	got, err = client.UpdateHourEntry(
		&intranet.UpdateHourEntryInput{
			Date:        intranet.Date(now),
			Description: "Test Updated",
			ID:          got.ID,
			ProjectID:   422,
			TicketID:    "TEST-UPDATED",
			Time:        2.5,
		})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, ignoreFields); diff != "" {
		t.Error(diff)
	}

	// Delete
	err = client.DeleteHourEntry(
		&intranet.DeleteHourEntryInput{
			ID: got.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
