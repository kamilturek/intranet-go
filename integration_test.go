//go:build integration
// +build integration

package intranet_test

import (
	"os"
	"strconv"
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

func TestGetHourEntry(t *testing.T) {
	c := getClient(t)

	resCreate, err := c.CreateHourEntry(&intranet.CreateHourEntryInput{
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
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: resCreate.ID}); err != nil {
			t.Fatalf("failed to clean up after the test: %v", err)
		}
	}()

	id, err := strconv.Atoi(resCreate.ID)
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.GetHourEntry(&intranet.GetHourEntryInput{
		ID:   id,
		Date: time.Now().Format(intranet.DateFormat),
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedDescription := "Test"
	gotDescription := res.Description
	if expectedDescription != gotDescription {
		t.Fatalf("expected: %s, got: %s", expectedDescription, gotDescription)
	}

	expectedProjectID := TestProjectID
	gotProjectID := res.Project.ID
	if expectedProjectID != gotProjectID {
		t.Fatalf("expected: %d, got %d", expectedProjectID, gotProjectID)
	}

	expectedTime := 0.25
	gotTime := res.Time
	if expectedTime != gotTime {
		t.Fatalf("expected: %f, got: %f", expectedTime, gotTime)
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

	resCreate, err := c.CreateHourEntry(&intranet.CreateHourEntryInput{
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
		if err := c.DeleteHourEntry(&intranet.DeleteHourEntryInput{ID: resCreate.ID}); err != nil {
			t.Fatalf("failed to clean up after the test: %v", err)
		}
	}()

	res, err := c.UpdateHourEntry(&intranet.UpdateHourEntryInput{
		Date:        time.Now().Format(intranet.DateFormat),
		Description: "Test Updated",
		ID:          resCreate.ID,
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

func TestDeleteHourEntry(t *testing.T) {
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

	err = c.DeleteHourEntry(&intranet.DeleteHourEntryInput{
		ID: res.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(res.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.GetHourEntry(&intranet.GetHourEntryInput{
		ID:   id,
		Date: time.Now().Format(intranet.DateFormat),
	})
	if err == nil {
		t.Fatalf("expected err, got nil")
	}

	expecterErr := "hour entry not found"
	gotErr := err.Error()
	if expecterErr != gotErr {
		t.Fatalf("expected: %s, got: %s", expecterErr, gotErr)
	}
}
