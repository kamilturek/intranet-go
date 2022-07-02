package intranet_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/kamilturek/intranet-go"
)

func GetClient(t *testing.T, cassetteName string) (*intranet.Client, func()) {
	r, err := recorder.New(fmt.Sprintf("fixtures/%s", cassetteName))
	if err != nil {
		t.Fatal(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Cookie")
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	sessionID := os.Getenv(intranet.SessionIDEnvVar)
	client := intranet.NewClient(sessionID)
	client.HTTPClient.Transport = r

	deferFunc := func() {
		err = r.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}

	return client, deferFunc
}
