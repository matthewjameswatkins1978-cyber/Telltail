package cloudgcp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSpecUsesExplicitServiceAccountAndValidJSON(t *testing.T) {
	b, err := Spec(Config{Image: "example.invalid/worker:1", Command: "echo ok", ServiceAccount: "worker@example.iam.gserviceaccount.com"})
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	ap := v["allocationPolicy"].(map[string]any)
	sa := ap["serviceAccount"].(map[string]any)
	if sa["email"] != "worker@example.iam.gserviceaccount.com" {
		t.Fatalf("service account missing: %#v", sa)
	}
}

func TestClassifyCloudErrors(t *testing.T) {
	cases := map[string]string{
		`unauthorized_client rejected by the attribute condition`:                     "identity_rejected",
		`PERMISSION_DENIED caller does not have permission to act as service account`: "service_account_authority",
		`INVALID_ARGUMENT Invalid JSON payload received`:                              "job_spec_invalid",
		`RESOURCE_EXHAUSTED quota exceeded`:                                           "resource_exhausted",
	}
	for in, want := range cases {
		if got := classifyCloudError(in); got != want {
			t.Fatalf("%q => %q want %q", in, got, want)
		}
	}
}

func TestSpecRequiresImageAndCommand(t *testing.T) {
	if _, err := Spec(Config{Command: "echo ok"}); err == nil || !strings.Contains(err.Error(), "image") {
		t.Fatalf("expected image error, got %v", err)
	}
	if _, err := Spec(Config{Image: "x"}); err == nil || !strings.Contains(err.Error(), "command") {
		t.Fatalf("expected command error, got %v", err)
	}
}
