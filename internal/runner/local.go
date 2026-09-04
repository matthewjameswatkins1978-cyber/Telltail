package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/trace"
)

type LocalConfig struct {
	Dir       string
	Worker    string
	Command   string
	TracePath string
	Timeout   time.Duration
	Shell     string
}

func shellCommand(ctx context.Context, shell, command string) *exec.Cmd {
	if shell != "" {
		base := strings.ToLower(strings.TrimSpace(shell))
		if strings.Contains(base, "powershell") || strings.Contains(base, "pwsh") {
			return exec.CommandContext(ctx, shell, "-NoProfile", "-Command", command)
		}
		if strings.Contains(base, "cmd") {
			return exec.CommandContext(ctx, shell, "/C", command)
		}
		return exec.CommandContext(ctx, shell, "-lc", command)
	}
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd.exe", "/C", command)
	}
	return exec.CommandContext(ctx, "/bin/sh", "-lc", command)
}

func RunLocal(c LocalConfig) error {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Minute
	}
	r, err := trace.NewRecorder(c.TracePath)
	if err != nil {
		return err
	}
	defer r.Close()
	_, _ = r.Append(model.Event{Type: model.EventShiftStart, Actor: c.Worker, Name: "local"})
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	start := time.Now()
	cmd := shellCommand(ctx, c.Shell, c.Command)
	cmd.Dir = c.Dir
	cmd.Env = append(os.Environ(), "TELLTAIL_TRACE="+c.TracePath)
	out, runErr := cmd.CombinedOutput()
	dur := time.Since(start).Milliseconds()
	ok := runErr == nil
	_, _ = r.Append(model.Event{Type: model.EventCommandResult, Actor: c.Worker, Name: "worker_process", Input: c.Command, Output: string(out), Success: &ok, DurationMS: dur, Progress: ok})
	_, _ = r.Append(model.Event{Type: model.EventShiftEnd, Actor: c.Worker, Name: "local", Success: &ok})
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("worker timeout after %s", c.Timeout)
	}
	if runErr != nil {
		return fmt.Errorf("worker failed: %w: %s", runErr, strings.TrimSpace(string(out)))
	}
	return nil
}
