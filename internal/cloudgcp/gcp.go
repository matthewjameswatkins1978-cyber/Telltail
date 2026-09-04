package cloudgcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/trace"
)

type Config struct {
	Project        string `json:"project"`
	Region         string `json:"region"`
	Image          string `json:"image"`
	JobID          string `json:"job_id"`
	Command        string `json:"command"`
	CPU            int    `json:"cpu_milli"`
	MemoryMiB      int    `json:"memory_mib"`
	MaxRetries     int    `json:"max_retries"`
	MaxRun         string `json:"max_run_duration"`
	ServiceAccount string `json:"service_account,omitempty"`
}

type RunOptions struct {
	TracePath string
	Poll      time.Duration
	Timeout   time.Duration
}

type batchJob struct {
	TaskGroups       []taskGroup       `json:"taskGroups"`
	AllocationPolicy map[string]any    `json:"allocationPolicy,omitempty"`
	LogsPolicy       map[string]string `json:"logsPolicy"`
	Labels           map[string]string `json:"labels,omitempty"`
}
type taskGroup struct {
	TaskSpec    taskSpec `json:"taskSpec"`
	TaskCount   int      `json:"taskCount"`
	Parallelism int      `json:"parallelism"`
}
type taskSpec struct {
	Runnables       []runnable `json:"runnables"`
	ComputeResource resource   `json:"computeResource"`
	MaxRetryCount   int        `json:"maxRetryCount"`
	MaxRunDuration  string     `json:"maxRunDuration"`
}
type runnable struct {
	Container container `json:"container"`
}
type container struct {
	ImageURI   string   `json:"imageUri"`
	EntryPoint string   `json:"entrypoint"`
	Commands   []string `json:"commands"`
}
type resource struct {
	CPUMilli  int `json:"cpuMilli"`
	MemoryMiB int `json:"memoryMib"`
}

type describeResult struct {
	Status struct {
		State string `json:"state"`
	} `json:"status"`
}

func Defaults(c Config) Config {
	if c.CPU == 0 {
		c.CPU = 2000
	}
	if c.MemoryMiB == 0 {
		c.MemoryMiB = 2048
	}
	if c.MaxRun == "" {
		c.MaxRun = "1800s"
	}
	if c.JobID == "" {
		c.JobID = fmt.Sprintf("telltail-shift-%d", time.Now().Unix())
	}
	return c
}

func Spec(c Config) ([]byte, error) {
	c = Defaults(c)
	if c.Image == "" {
		return nil, fmt.Errorf("image is required")
	}
	if c.Command == "" {
		return nil, fmt.Errorf("command is required")
	}
	j := batchJob{
		TaskGroups: []taskGroup{{TaskSpec: taskSpec{
			Runnables:       []runnable{{Container: container{ImageURI: c.Image, EntryPoint: "/bin/sh", Commands: []string{"-lc", c.Command}}}},
			ComputeResource: resource{CPUMilli: c.CPU, MemoryMiB: c.MemoryMiB}, MaxRetryCount: c.MaxRetries, MaxRunDuration: c.MaxRun,
		}, TaskCount: 1, Parallelism: 1}},
		LogsPolicy: map[string]string{"destination": "CLOUD_LOGGING"},
		Labels:     map[string]string{"app": "telltail", "kind": "cloud-shift"},
	}
	if c.ServiceAccount != "" {
		j.AllocationPolicy = map[string]any{"serviceAccount": map[string]string{"email": c.ServiceAccount}}
	}
	return json.MarshalIndent(j, "", "  ")
}

func WriteSpec(path string, c Config) error {
	b, e := Spec(c)
	if e != nil {
		return e
	}
	return os.WriteFile(path, b, 0644)
}

func commandOutput(ctx context.Context, name string, args ...string) ([]byte, time.Duration, error) {
	start := time.Now()
	b, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return b, time.Since(start), err
}

func Submit(c Config, specPath string) error {
	c = Defaults(c)
	cmd := exec.Command("gcloud", "batch", "jobs", "submit", c.JobID, "--project", c.Project, "--location", c.Region, "--config", specPath)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

func Describe(c Config) error {
	c = Defaults(c)
	cmd := exec.Command("gcloud", "batch", "jobs", "describe", c.JobID, "--project", c.Project, "--location", c.Region, "--format=json")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func classifyCloudError(s string) string {
	x := strings.ToLower(s)
	switch {
	case strings.Contains(x, "attribute condition") || strings.Contains(x, "unauthorized_client"):
		return "identity_rejected"
	case strings.Contains(x, "does not have permission to act as service account"):
		return "service_account_authority"
	case strings.Contains(x, "invalid json payload") || strings.Contains(x, "invalid_argument"):
		return "job_spec_invalid"
	case strings.Contains(x, "permission_denied") || strings.Contains(x, "permission denied"):
		return "permission_denied"
	case strings.Contains(x, "resource_exhausted") || strings.Contains(x, "quota"):
		return "resource_exhausted"
	default:
		return "cloud_submit_failed"
	}
}

// Run submits a real Batch job and records the observable remote lifecycle in
// the same hash-chained trace format used by local shifts. Worker-internal
// tool/model/file events can be merged later from the worker's durable trace.
func Run(c Config, specPath string, opt RunOptions) error {
	c = Defaults(c)
	if c.Project == "" {
		return fmt.Errorf("project is required")
	}
	if c.Region == "" {
		return fmt.Errorf("region is required")
	}
	if opt.TracePath == "" {
		opt.TracePath = "telltail-cloud-trace.jsonl"
	}
	if opt.Poll <= 0 {
		opt.Poll = 5 * time.Second
	}
	if opt.Timeout <= 0 {
		opt.Timeout = 30 * time.Minute
	}
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		if err := WriteSpec(specPath, c); err != nil {
			return err
		}
	}
	rec, err := trace.NewRecorder(opt.TracePath)
	if err != nil {
		return err
	}
	defer rec.Close()
	ok := true
	_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventShiftStart, Actor: "telltail", Name: "cloud_shift", Success: &ok, Metadata: map[string]string{"backend": "gcp-batch", "job_id": c.JobID, "region": c.Region}})

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()
	out, dur, err := commandOutput(ctx, "gcloud", "batch", "jobs", "submit", c.JobID, "--project", c.Project, "--location", c.Region, "--config", specPath)
	submitOK := err == nil
	md := map[string]string{"stage": "submit"}
	if err != nil {
		md["error_kind"] = classifyCloudError(string(out))
	}
	_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventCommandResult, Actor: "cloud", Name: "batch_submit", Output: string(out), Success: &submitOK, DurationMS: dur.Milliseconds(), Metadata: md})
	if err != nil {
		endOK := false
		_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventShiftEnd, Actor: "telltail", Name: "cloud_shift", Success: &endOK, Metadata: map[string]string{"stage": "submit"}})
		return fmt.Errorf("batch submit: %w", err)
	}

	submitted := time.Now()
	lastState := ""
	stateSince := submitted
	for {
		if err := ctx.Err(); err != nil {
			endOK := false
			_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventShiftEnd, Actor: "telltail", Name: "cloud_shift", Success: &endOK, Metadata: map[string]string{"stage": "poll", "error_kind": "timeout"}})
			return err
		}
		b, _, derr := commandOutput(ctx, "gcloud", "batch", "jobs", "describe", c.JobID, "--project", c.Project, "--location", c.Region, "--format=json")
		if derr != nil {
			time.Sleep(opt.Poll)
			continue
		}
		var d describeResult
		if json.Unmarshal(b, &d) != nil || d.Status.State == "" {
			time.Sleep(opt.Poll)
			continue
		}
		state := strings.ToUpper(d.Status.State)
		if state != lastState {
			now := time.Now()
			if lastState != "" {
				elapsed := now.Sub(stateSince).Milliseconds()
				evType := model.EventInfo
				switch lastState {
				case "QUEUED":
					evType = model.EventQueue
				case "SCHEDULED":
					evType = model.EventProvision
				}
				_, _ = rec.Append(model.Event{Time: now.UTC(), Type: evType, Actor: "cloud", Name: strings.ToLower(lastState), DurationMS: elapsed, Metadata: map[string]string{"stage": "remote_lifecycle"}})
			}
			lastState, stateSince = state, now
		}
		switch state {
		case "SUCCEEDED":
			endOK := true
			_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventShiftEnd, Actor: "telltail", Name: "cloud_shift", Success: &endOK, Metadata: map[string]string{"stage": "complete"}})
			return nil
		case "FAILED", "DELETION_IN_PROGRESS":
			endOK := false
			_, _ = rec.Append(model.Event{Time: time.Now().UTC(), Type: model.EventShiftEnd, Actor: "telltail", Name: "cloud_shift", Success: &endOK, Metadata: map[string]string{"stage": "remote", "state": state}})
			return fmt.Errorf("batch job ended in state %s", state)
		}
		time.Sleep(opt.Poll)
	}
}
