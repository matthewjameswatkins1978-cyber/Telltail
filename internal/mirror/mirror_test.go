package mirror

import "testing"

func TestScore(t *testing.T) {
	truth := []TruthCase{{ID: "a", Required: map[string]int{"phantom_tool": 2}}}
	s := Submission{Evaluator: "e", Cases: []EvaluatorCase{{ID: "a", Findings: map[string]int{"phantom_tool": 2}}}}
	got := ScoreSubmission(truth, s)
	if got.Recall != 1 || got.Precision != 1 || got.IntegrityScore != 1 {
		t.Fatalf("%+v", got)
	}
}
