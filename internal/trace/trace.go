package trace

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
)

type Recorder struct {
	f        *os.File
	enc      *json.Encoder
	seq      int64
	prevHash string
}

func NewRecorder(path string) (*Recorder, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	return &Recorder{f: f, enc: json.NewEncoder(f)}, nil
}

func hashEvent(e model.Event) (string, error) {
	clone := e
	clone.Hash = ""
	b, err := json.Marshal(clone)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

func (r *Recorder) Append(e model.Event) (model.Event, error) {
	r.seq++
	e.Version = 1
	e.Seq = r.seq
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	e.PrevHash = r.prevHash
	h, err := hashEvent(e)
	if err != nil {
		return e, err
	}
	e.Hash = h
	if err := r.enc.Encode(e); err != nil {
		return e, err
	}
	r.prevHash = h
	return e, nil
}

func (r *Recorder) Head() string { return r.prevHash }
func (r *Recorder) Close() error { return r.f.Close() }

func Read(path string) ([]model.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []model.Event
	s := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, 4*1024*1024)
	for s.Scan() {
		var e model.Event
		if err := json.Unmarshal(s.Bytes(), &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, s.Err()
}

func Verify(events []model.Event) error {
	prev := ""
	for i, e := range events {
		if e.Seq != int64(i+1) {
			return fmt.Errorf("sequence break at event %d", i+1)
		}
		if e.PrevHash != prev {
			return fmt.Errorf("hash-chain predecessor mismatch at seq %d", e.Seq)
		}
		expected, err := hashEvent(e)
		if err != nil {
			return err
		}
		if e.Hash != expected {
			return fmt.Errorf("hash mismatch at seq %d", e.Seq)
		}
		prev = e.Hash
	}
	return nil
}
