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
		t.Fatal(err)
	}

	defer func() {
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: res.ID}); err != nil {
			t.Fatalf("failed to clean up after the test: %v", err)
		}
	}()

	expectedAdded := time.Now().Format(intranet.DateFormat)
	gotAdded := res.Added
	if expectedAdded != gotAdded {
		t.Fatalf("expected: %s, got: %s", expectedAdded, gotAdded)
	}

	expectedTime := 0.25
	gotTime := res.Time
	if expectedTime != gotTime {
		t.Fatalf("expected: %f, got: %f", expectedTime, gotTime)
	}

	expectedDescription := "Test"
	gotDescription := res.Description
	if expectedDescription != gotDescription {
		t.Fatalf("expected: %s, got: %s", expectedDescription, gotDescription)
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
		t.Fatal(err)
	}

	defer func() {
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: res.ID}); err != nil {
			t.Fatalf("failed to clean up after the test: %v", err)
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

	expectedTime := 0.5
	gotTime := res.Time
	if expectedTime != gotTime {
		t.Fatalf("expected: %f, got: %f", expectedTime, gotTime)
	}

	expectedDescription := "Test Updated"
	gotDescription := res.Description
	if expectedDescription != gotDescription {
		t.Fatalf("expected: %s, got: %s", expectedDescription, gotDescription)
	}
}
