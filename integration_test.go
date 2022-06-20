//go:build integration
// +build integration

package intranet_test

import (
	"os"
	"testing"
	"time"

	"github.com/kamilturek/intranet-go"
)

const (
	TestProjectID int = 1185 // Cross-company initiatives & improvements
)

func getClient(t *testing.T) *intranet.Client {
	sessionID := os.Getenv(intranet.SessionIDEnvVar)
	if len(sessionID) == 0 {
		t.Fatalf("%s environment variable is not set.", intranet.SessionIDEnvVar)
	}

	return intranet.NewClient(sessionID)
}

func TestGetHourEntries(t *testing.T) {
	c := getClient(t)

	res, err := c.GetHourEntries(&intranet.GetHourEntriesInput{Date: "2022-05-20"})
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if len(res.Entries) != 8 {
		t.Fatalf("expected: 8, got: %d", len(res.Entries))
	}

	if res.Entries[0].Project.ClientName != "Scurri Web Services Limited" {
		t.Fatalf("expected: Scurri Web Services Limited, got: %s", res.Entries[0].Project.ClientName)
	}

	if res.Entries[0].Project.ID != 422 {
		t.Fatalf("expected: 422, got: %d", res.Entries[0].Project.ID)
	}

	if res.Entries[0].Project.Name != "Shadow Unicorn (Scurri) / WRO / AyeAye / Billable" {
		t.Fatalf("expected: Shadow Unicorn (Scurri) / WRO / AyeAye / Billable, got: %s", res.Entries[0].Project.Name)
	}
}

func TestCreateHourEntry(t *testing.T) {
	c := getClient(t)

	res, err := c.CreateHourEntry(&intranet.CreateHourEntryInput{
		Date:        time.Now().Format(intranet.DateFormat),
		Description: "Test",
		ProjectID:   TestProjectID,
		TicketID:    "",
		Time:        0.25,
	})
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	defer func() {
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: res.ID}); err != nil {
			t.Fatalf("failed to clean up the test: %v", err)
		}
	}()

	if res.Added != time.Now().Format(intranet.DateFormat) {
		t.Fatalf("expected: %s, got: %s", time.Now().Format(intranet.DateFormat), res.Added)
	}

	if res.Time != 0.25 {
		t.Fatalf("expected: %f, got: %f", 0.25, res.Time)
	}

	if res.Description != "Test" {
		t.Fatalf("expected: %s, got: %s", "Test", res.Description)
	}
}

func TestUpdateHourEntry(t *testing.T) {
	c := getClient(t)

	res, err := c.CreateHourEntry(&intranet.CreateHourEntryInput{
		Date:        time.Now().Format(intranet.DateFormat),
		Description: "Test",
		ProjectID:   TestProjectID,
		TicketID:    "",
		Time:        0.25,
	})
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	defer func() {
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: res.ID}); err != nil {
			t.Fatalf("failed to clean up the test: %v", err)
		}
	}()

	res, err = c.UpdateHourEntry(&intranet.UpdateHourEntryInput{
		Date:        time.Now().Format(intranet.DateFormat),
		Description: "Test Updated",
		ID:          res.ID,
		ProjectID:   TestProjectID,
		TicketID:    "",
		Time:        0.5,
	})
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if res.Time != 0.5 {
		t.Fatalf("expected: %f, got: %f", 0.5, res.Time)
	}

	if res.Description != "Test Updated" {
		t.Fatalf("expected: %s, got: %s", "Test Updated", res.Description)
	}
}
