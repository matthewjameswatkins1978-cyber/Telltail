package trace

import (
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
	"path/filepath"
	"testing"
)

func TestHashChain(t *testing.T) {
	p := filepath.Join(t.TempDir(), "t.jsonl")
	r, e := NewRecorder(p)
	if e != nil {
		t.Fatal(e)
	}
	ok := true
	if _, e = r.Append(model.Event{Type: model.EventToolCall, Name: "git"}); e != nil {
		t.Fatal(e)
	}
	if _, e = r.Append(model.Event{Type: model.EventToolResult, Name: "git", Success: &ok}); e != nil {
		t.Fatal(e)
	}
	r.Close()
	ev, e := Read(p)
	if e != nil {
		t.Fatal(e)
	}
	if e = Verify(ev); e != nil {
		t.Fatal(e)
	}
	ev[1].Name = "evil"
	if Verify(ev) == nil {
		t.Fatal("expected tamper detection")
	}
}
