package analysis

import (
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/scenario"
	"testing"
)

func bp(v bool) *bool { return &v }
func TestDetectors(t *testing.T) {
	s := scenario.Scenario{ID: "x", Tools: []scenario.Tool{{Name: "git", Available: true}, {Name: "powershell", Available: false}}}
	ev := []model.Event{{Seq: 1, Type: model.EventToolCall, Name: "powershell"}, {Seq: 2, Type: model.EventCommandResult, Name: "build", Input: "make", Success: bp(false), DurationMS: 10}, {Seq: 3, Type: model.EventCommandResult, Name: "build", Input: "make", Success: bp(false), DurationMS: 20}, {Seq: 4, Type: model.EventClaim, Name: "success"}}
	r := Analyze(ev, s, "w", "local")
	counts := map[string]int{}
	for _, f := range r.Findings {
		counts[f.Code]++
	}
	if counts["phantom_tool"] != 1 {
		t.Fatalf("phantom=%d", counts["phantom_tool"])
	}
	if counts["repeated_failure"] != 1 {
		t.Fatalf("repeat=%d", counts["repeated_failure"])
	}
	if counts["false_success"] != 1 {
		t.Fatalf("false success=%d", counts["false_success"])
	}
	if r.Profile.RepeatedWorkMS != 20 {
		t.Fatalf("repeat ms=%d", r.Profile.RepeatedWorkMS)
	}
}

func TestRecidivismThrashAndFalseBlocker(t *testing.T) {
	s := scenario.Scenario{ID: "x", Tools: []scenario.Tool{{Name: "ghost", Available: false}}}
	ev := []model.Event{
		{Seq: 1, Type: model.EventToolCall, Name: "ghost"},
		{Seq: 2, Type: model.EventToolCall, Name: "ghost"},
		{Seq: 3, Type: model.EventCommandResult, Name: "a", Input: "a", Success: bp(false)},
		{Seq: 4, Type: model.EventCommandResult, Name: "b", Input: "b", Success: bp(false)},
		{Seq: 5, Type: model.EventCommandResult, Name: "c", Input: "c", Success: bp(false)},
		{Seq: 6, Type: model.EventCommandResult, Name: "d", Input: "d", Success: bp(false)},
		{Seq: 7, Type: model.EventClaim, Name: "blocked", Metadata: map[string]string{"route_available": "true"}},
	}
	r := Analyze(ev, s, "w", "local")
	counts := map[string]int{}
	for _, f := range r.Findings {
		counts[f.Code]++
	}
	if counts["repeated_phantom_tool"] != 1 {
		t.Fatalf("repeated phantom=%d", counts["repeated_phantom_tool"])
	}
	if counts["thrash_signal"] != 1 {
		t.Fatalf("thrash=%d", counts["thrash_signal"])
	}
	if counts["false_blocker"] != 1 {
		t.Fatalf("false blocker=%d", counts["false_blocker"])
	}
}

func TestOutcomeClassification(t *testing.T) {
	s := scenario.Scenario{ID: "x"}
	clean := Analyze([]model.Event{{Seq: 1, Type: model.EventAcceptance, Success: bp(true)}}, s, "w", "local")
	if clean.Outcome != "clean_success" {
		t.Fatalf("clean outcome=%q", clean.Outcome)
	}
	useful := Analyze([]model.Event{{Seq: 1, Type: model.EventProgress, Progress: true}}, s, "w", "local")
	if useful.Outcome != "useful_failure" {
		t.Fatalf("useful outcome=%q", useful.Outcome)
	}
}
