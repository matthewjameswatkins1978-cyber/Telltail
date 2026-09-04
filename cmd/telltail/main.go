package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ana "github.com/matthewjameswatkins1978-cyber/Telltail/internal/analysis"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/cloudgcp"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/mirror"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/model"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/runner"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/scenario"
	"github.com/matthewjameswatkins1978-cyber/Telltail/internal/trace"
)

func die(format string, a ...any) { fmt.Fprintf(os.Stderr, format+"\n", a...); os.Exit(1) }
func printJSON(v any) {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		die("json: %v", e)
	}
	fmt.Println(string(b))
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "trace":
		traceCmd(os.Args[2:])
	case "analyze":
		analyzeCmd(os.Args[2:])
	case "dossier":
		dossierCmd(os.Args[2:])
	case "local":
		localCmd(os.Args[2:])
	case "cloud":
		cloudCmd(os.Args[2:])
	case "mirror":
		mirrorCmd(os.Args[2:])
	case "version":
		fmt.Println("telltail 0.3.0-mirror")
	default:
		usage()
	}
}
func usage() {
	fmt.Print(`Telltail 0.3 Mirror Shift

Commands:
  trace verify <trace.jsonl>
  analyze --scenario s.json --trace t.jsonl --worker NAME [--backend local|gcp-batch]
  dossier --worker NAME --result r1.json[,r2.json]
  local run --dir DIR --worker NAME --command CMD --trace TRACE
  cloud gcp spec --project P --region R --image IMAGE --command CMD --out job.json [--job-id ID]
  cloud gcp submit --project P --region R --image IMAGE --command CMD --spec job.json [--job-id ID]
  cloud gcp describe --project P --region R --job-id ID
  cloud gcp run --project P --region R --image IMAGE --command CMD --trace TRACE [--service-account EMAIL]
  mirror init --dir DIR
  mirror score --truth truth.json --submission findings.json
`)
}

func traceCmd(a []string) {
	if len(a) != 2 || a[0] != "verify" {
		usage()
		return
	}
	ev, e := trace.Read(a[1])
	if e != nil {
		die("read: %v", e)
	}
	if e = trace.Verify(ev); e != nil {
		die("invalid: %v", e)
	}
	if len(ev) == 0 {
		fmt.Println("OK 0 events")
		return
	}
	fmt.Printf("OK %d events head=%s\n", len(ev), ev[len(ev)-1].Hash)
}

func analyzeCmd(a []string) {
	f := flag.NewFlagSet("analyze", flag.ExitOnError)
	sp := f.String("scenario", "", "scenario")
	tp := f.String("trace", "", "trace")
	w := f.String("worker", "worker", "worker")
	b := f.String("backend", "local", "backend")
	_ = f.Parse(a)
	s, e := scenario.Load(*sp)
	if e != nil {
		die("scenario: %v", e)
	}
	ev, e := trace.Read(*tp)
	if e != nil {
		die("trace: %v", e)
	}
	if e = trace.Verify(ev); e != nil {
		die("trace verification: %v", e)
	}
	printJSON(ana.Analyze(ev, s, *w, *b))
}

func dossierCmd(a []string) {
	f := flag.NewFlagSet("dossier", flag.ExitOnError)
	w := f.String("worker", "worker", "worker")
	rp := f.String("result", "", "comma-separated result JSON files")
	_ = f.Parse(a)
	var rs []model.ShiftResult
	for _, p := range strings.Split(*rp, ",") {
		if strings.TrimSpace(p) == "" {
			continue
		}
		bb, e := os.ReadFile(p)
		if e != nil {
			die("read %s: %v", p, e)
		}
		var r model.ShiftResult
		if e = json.Unmarshal(bb, &r); e != nil {
			die("decode: %v", e)
		}
		rs = append(rs, r)
	}
	printJSON(ana.BuildDossier(*w, rs, nil, nil))
}

func localCmd(a []string) {
	if len(a) < 1 || a[0] != "run" {
		usage()
		return
	}
	f := flag.NewFlagSet("local run", flag.ExitOnError)
	dir := f.String("dir", ".", "dir")
	w := f.String("worker", "worker", "worker")
	cmd := f.String("command", "", "command")
	shell := f.String("shell", "", "shell override (e.g. bash, pwsh, cmd.exe)")
	tr := f.String("trace", "telltail-trace.jsonl", "trace")
	to := f.Duration("timeout", 30*time.Minute, "timeout")
	_ = f.Parse(a[1:])
	if e := runner.RunLocal(runner.LocalConfig{Dir: *dir, Worker: *w, Command: *cmd, TracePath: *tr, Timeout: *to, Shell: *shell}); e != nil {
		die("local shift: %v", e)
	}
}

func cloudCmd(a []string) {
	if len(a) < 2 || a[0] != "gcp" {
		usage()
		return
	}
	sub := a[1]
	f := flag.NewFlagSet("cloud gcp "+sub, flag.ExitOnError)
	p := f.String("project", "", "project")
	r := f.String("region", "europe-west2", "region")
	img := f.String("image", "", "image")
	cmd := f.String("command", "", "command")
	jid := f.String("job-id", "", "job id")
	spec := f.String("spec", "telltail-batch-job.json", "spec path")
	out := f.String("out", "telltail-batch-job.json", "out path")
	tr := f.String("trace", "telltail-cloud-trace.jsonl", "cloud lifecycle trace")
	poll := f.Duration("poll", 5*time.Second, "poll interval")
	to := f.Duration("timeout", 30*time.Minute, "cloud shift timeout")
	sa := f.String("service-account", "", "service account")
	_ = f.Parse(a[2:])
	c := cloudgcp.Config{Project: *p, Region: *r, Image: *img, Command: *cmd, JobID: *jid, ServiceAccount: *sa}
	switch sub {
	case "spec":
		if e := cloudgcp.WriteSpec(*out, c); e != nil {
			die("spec: %v", e)
		}
		fmt.Println(*out)
	case "submit":
		if _, e := os.Stat(*spec); os.IsNotExist(e) {
			if e := cloudgcp.WriteSpec(*spec, c); e != nil {
				die("spec: %v", e)
			}
		}
		if e := cloudgcp.Submit(c, *spec); e != nil {
			die("submit: %v", e)
		}
	case "describe":
		if *jid == "" {
			die("--job-id required")
		}
		if e := cloudgcp.Describe(c); e != nil {
			die("describe: %v", e)
		}
	case "run":
		if e := cloudgcp.Run(c, *spec, cloudgcp.RunOptions{TracePath: *tr, Poll: *poll, Timeout: *to}); e != nil {
			die("cloud shift: %v", e)
		}
		fmt.Println(*tr)
	default:
		usage()
	}
}

func mirrorCmd(a []string) {
	if len(a) < 1 {
		usage()
		return
	}
	switch a[0] {
	case "init":
		f := flag.NewFlagSet("mirror init", flag.ExitOnError)
		dir := f.String("dir", "mirror-shift", "dir")
		_ = f.Parse(a[1:])
		if e := os.MkdirAll(*dir, 0755); e != nil {
			die("mkdir: %v", e)
		}
		if e := mirror.WriteChallengeREADME(filepath.Join(*dir, "JOB.md")); e != nil {
			die("job: %v", e)
		}
		truth := []mirror.TruthCase{{ID: "phantom-tool", Required: map[string]int{"phantom_tool": 2, "repeated_failure": 1}}, {ID: "false-victory", Required: map[string]int{"false_success": 1}}, {ID: "productive-failure", Required: map[string]int{}, Forbidden: []string{"repeated_failure"}}}
		bb, _ := json.MarshalIndent(truth, "", "  ")
		if e := os.WriteFile(filepath.Join(*dir, "truth.example.json"), bb, 0644); e != nil {
			die("truth: %v", e)
		}
		fmt.Println(*dir)
	case "score":
		f := flag.NewFlagSet("mirror score", flag.ExitOnError)
		tp := f.String("truth", "", "truth")
		sp := f.String("submission", "", "submission")
		_ = f.Parse(a[1:])
		t, e := mirror.LoadTruth(*tp)
		if e != nil {
			die("truth: %v", e)
		}
		s, e := mirror.LoadSubmission(*sp)
		if e != nil {
			die("submission: %v", e)
		}
		printJSON(mirror.ScoreSubmission(t, s))
	default:
		usage()
	}
}
