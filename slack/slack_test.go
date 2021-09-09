package slack

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"
	"github.com/teleivo/diskmon/usage"
)

func TestFormatMessage(t *testing.T) {
	r := usage.Report{
		Limit: 50,
		Limits: []usage.Stats{
			{
				Path:  "/mnt/volume1",
				Free:  2500000000,
				Total: 5000000000,
			},
			{
				Path:  "/mnt/volume2",
				Free:  500000000,
				Total: 5000000000,
			},
		},
		Errors: []error{
			errors.New("error reading stats: no such file 'disk123'"),
			errors.New("error reading stats: no permission to read 'disk981'"),
		},
	}

	// This test ensures that usage reports are turned into pretty Slack block
	// messages.
	// We only create the blocks field of a Slack message (JSON)
	// and compare it to a stored snapshot. Refactoring is safe as long as the
	// formatMessage() output is equal to the snapshot.
	// The snapshot can be copy & pasted into
	//
	// https://app.slack.com/block-kit-builder
	//
	// and iterated on. Once the new design is done, update the snapshot and
	// make the test pass again :)
	msg := struct {
		Blocks slack.Blocks `json:"blocks,omitempty"`
	}{
		slack.Blocks{BlockSet: formatMessage(r, "staging-01")},
	}
	got, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		t.Fatalf("failed to marshal the slack block message: %v", err)
	}
	// got := string(b)

	want, err := os.ReadFile("testdata/slack-message.snapshot")
	if err != nil {
		t.Fatalf("failed to read the slack block message snapshot: %v", err)
	}
	// want := string(c)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("formatMessage() mismatch (-want +got): \n%s", diff)
	}
}
