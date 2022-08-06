//go:build integration

package intranet_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	client, err := intranet.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// Create
	createdEntry, err := client.CreateHourEntry(
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

	// Get
	gotEntry, err := client.GetHourEntry(
		&intranet.GetHourEntryInput{
			ID:   createdEntry.ID,
			Date: intranet.Date(now),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(createdEntry, gotEntry) {
		t.Fatal(cmp.Diff(createdEntry, gotEntry))
	}

	// Update
	updatedEntry, err := client.UpdateHourEntry(
		&intranet.UpdateHourEntryInput{
			Date:        intranet.Date(now),
			Description: "Test Updated",
			ID:          gotEntry.ID,
			ProjectID:   422,
			TicketID:    "TEST-UPDATED",
			Time:        2.5,
		})
	if err != nil {
		t.Fatal(err)
	}

	// Get again
	gotUpdatedEntry, err := client.GetHourEntry(
		&intranet.GetHourEntryInput{
			ID:   updatedEntry.ID,
			Date: intranet.Date(now),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(updatedEntry, gotUpdatedEntry) {
		t.Fatal(cmp.Diff(updatedEntry, gotUpdatedEntry))
	}

	// Delete
	err = client.DeleteHourEntry(
		&intranet.DeleteHourEntryInput{
			ID: gotEntry.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
