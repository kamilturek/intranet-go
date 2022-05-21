//go:build integration
// +build integration

package intranet_test

import (
	"os"
	"testing"
	"time"

	"github.com/kamilturek/intranet"
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

	date, err := time.Parse(intranet.DateFormat, "2022-05-20")
	if err != nil {
		t.Fatalf("failed to parse the date: %v", err)
	}

	res, err := c.GetHourEntries(&intranet.GetHourEntriesOptions{Date: date})
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if len(res.Entries) != 8 {
		t.Fatalf("expected: 8, got: %d", len(res.Entries))
	}
}
