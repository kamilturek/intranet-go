//go:build integration

package intranet_test

import (
	"os"
	"strconv"
	"testing"
	"time"

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
	createInput := &intranet.CreateHourEntryInput{
		Date:        intranet.Date(now),
		Description: "Test",
		ProjectID:   422,
		TicketID:    "TEST",
		Time:        1.5,
	}
	createOutput, err := client.CreateHourEntry(createInput)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	updateInput := &intranet.UpdateHourEntryInput{
		Date:        intranet.Date(now),
		Description: "Test Updated",
		ID:          createOutput.ID,
		ProjectID:   422,
		TicketID:    "TEST-UPDATED",
		Time:        2.5,
	}
	updateOutput, err := client.UpdateHourEntry(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	// Get
	id, err := strconv.Atoi(updateOutput.ID)
	if err != nil {
		t.Fatal(err)
	}

	getInput := &intranet.GetHourEntryInput{
		ID:   id,
		Date: intranet.Date(now),
	}
	getOutput, err := client.GetHourEntry(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if strconv.Itoa(getOutput.ID) != createOutput.ID {
		t.Fatalf("want %s, got %d", createOutput.ID, getOutput.ID)
	}

	if getOutput.Description != "Test Updated" {
		t.Fatalf("want %s, got %s", "Test Updated", getOutput.Description)
	}

	if getOutput.Ticket.ID != "TEST-UPDATED" {
		t.Fatalf("want %s, got %s", "TEST-UPDATED", getOutput.Ticket.ID)
	}

	if getOutput.Time != 2.5 {
		t.Fatalf("want %f, got %f", 2.5, getOutput.Time)
	}

	// Delete
	deleteInput := &intranet.DeleteHourEntryInput{
		ID: createOutput.ID,
	}
	err = client.DeleteHourEntry(deleteInput)
	if err != nil {
		t.Fatal(err)
	}
}
