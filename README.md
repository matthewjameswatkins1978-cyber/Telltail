# Telltail

**Deterministic behavioural evaluation for autonomous AI workers.**

Current repository line: **0.3 Mirror Shift foundation**.

Telltail gives an AI worker a believable job, records what the worker visibly
does, checks the resulting evidence mechanically, and reports where the shift
was clean, wasteful, stuck, unsafe, or falsely declared complete. It is a
small evaluation and evidence tool: it is not an agent framework,
orchestrator, LLM judge, policy learner, or replacement for the worker itself.

## Repository status

This repository contains the verified Telltail 0.3 behavioural-lab foundation:
the CLI, hash-chained trace model, analyser, worker dossiers, local runner,
Google Cloud Batch backend, and Mirror Shift scorer.

Telltail 0.2 was the earlier deterministic certification harness. Its original
routing, invariance, delta, state-sequence, and mutation corpus is **not yet
imported** here. Read 0.3 as an additive foundation, not as a claim that the
full 0.2 corpus has already been merged. See
[`docs/INTEGRATION.md`](docs/INTEGRATION.md) for the intended integration
boundary.

The current code is verified with:

```bash
go test ./...
go vet ./...
```

## The mental model

The basic unit is a **shift**: one worker attempting one job under one set of
conditions.

1. A **scenario** describes the visible job, available tools, traps,
   acceptance conditions, and allowed or forbidden paths.
2. A **runner** gives the job to a worker locally or submits it to Google Cloud
   Batch.
3. The worker and runner produce a JSONL **trace** of observable events such
   as commands, tool calls, results, file changes, progress, claims,
   acceptance, queue time, and provisioning.
4. Telltail verifies the trace's append-only SHA-256 hash chain. If the trace
   was altered or is malformed, verification fails.
5. The **analyser** checks the verified evidence against the scenario and
   produces a shift result: acceptance, outcome class, findings, timing,
   cost, token counts, and the trace head.
6. **Dossiers** combine results from several shifts for one worker. This
   reveals patterns that a single run cannot: repeated mistakes, false
   successes, average cost, and completion rate.

In short:

```text
scenario -> worker shift -> hash-chained trace -> verify -> analyse -> dossier
```

Telltail judges observable behaviour, not hidden chain-of-thought. Useful
evidence includes tool calls, commands, results, file changes, acceptance
state, permissions, timing, cost, tokens, and externally visible claims.

## Install or build

With Go installed, install the published command:

```bash
go install github.com/matthewjameswatkins1978-cyber/Telltail/cmd/telltail@latest
```

Or clone the repository and build a local binary:

```bash
go build -o telltail ./cmd/telltail
```

On Windows the output is `telltail.exe`; on macOS or Linux it is `telltail`.
The examples below use `./telltail`. Substitute `telltail.exe` or the full
path to the binary on Windows.

The multi-line examples use Bash-style `\` continuation. In PowerShell,
run the same command on one line or replace the continuations with PowerShell
backticks.

To see the supported commands:

```bash
./telltail
./telltail version
```

## A complete local workflow

The checked-in example scenario is a useful starting point:
[`examples/scenarios/tool-friction.json`](examples/scenarios/tool-friction.json).
It asks a worker to update `app.txt`, describes the tools it may use, and
defines a tempting unavailable tool plus acceptance conditions.

### 1. Run a worker locally

Run the worker command inside a disposable or sacrificial directory:

```bash
./telltail local run \
  --dir ./sacrificial-repo \
  --worker example-worker \
  --command 'your-agent-command-here' \
  --trace shift.jsonl
```

The local runner starts the command in `--dir`, records the worker name, and
writes a JSONL trace to `shift.jsonl`. The trace includes the shift start,
the worker process result, and the shift end. A richer worker adapter can use
the `TELLTAIL_TRACE` environment variable to append its own tool, model,
file, and progress events to the same trace.

Use `--shell bash`, `--shell pwsh`, or `--shell cmd.exe` when the worker
command needs a particular shell. Use `--timeout 10m` to bound a run.

The command exits non-zero when the worker times out or fails. That is useful
operational feedback, but it is not the same thing as a completed behavioural
evaluation: analyse the resulting trace when you need the structured result.

### 2. Verify the trace

Before trusting a trace, verify its JSON shape and hash chain:

```bash
./telltail trace verify shift.jsonl
```

Successful output looks like:

```text
OK 3 events head=<sha256-hash>
```

The event count and final hash identify the evidence surface that was checked.
Malformed JSON, broken sequence or predecessor hashes, or any other chain
problem causes verification to fail. Do not treat an unverified trace as
certification evidence.

### 3. Analyse the shift

Give the analyser the scenario and the verified trace:

```bash
./telltail analyze \
  --scenario examples/scenarios/tool-friction.json \
  --trace shift.jsonl \
  --worker example-worker \
  --backend local > result.json
```

This writes a structured `ShiftResult` JSON file. It contains:

- whether the job was accepted and its outcome class;
- findings such as unavailable-tool calls, repeated failures, stuckness,
  false success, false blockers, and authority violations;
- time split across model, tools, commands, acceptance, and other phases;
- cost and token totals when the trace provides them; and
- the verified trace's final hash, `trace_head`.

Use `--backend gcp-batch` when analysing a trace whose backend is Google Cloud
Batch. Analysis is deterministic over the scenario and trace; it does not
launch another worker.

### 4. Build a worker dossier

After several shifts, aggregate their result files:

```bash
./telltail dossier \
  --worker example-worker \
  --result result-1.json,result-2.json,result-3.json > dossier.json
```

The dossier is a JSON summary across those shifts: completion rate, accepted
shifts, clean/recovered/messy successes, useful failures, repeated mistakes,
phantom tools, false successes, stuck loops, average duration, average cost,
and counts by finding code. Use it when choosing between workers or looking
for a stable behavioural pattern rather than judging one lucky run.

## Mirror Shift

Mirror Shift evaluates the evaluators. An agent receives public traces and
builds a small evaluator that reports findings; canonical Telltail then scores
that evaluator against hidden ground truth. The hidden traces and truth should
stay outside the evaluator's workspace. The scorer, not the evaluator, remains
the authority.

### 1. Create a challenge directory

```bash
./telltail mirror init --dir mirror-shift
```

This creates:

- `mirror-shift/JOB.md`, the challenge instructions; and
- `mirror-shift/truth.example.json`, a safe example of the truth format.

It does not create or expose the real hidden corpus.

### 2. Submit evaluator findings

An evaluator submission has one case per trace and a count for each finding it
believes it detected. For example:

```json
{
  "evaluator": "worker-x",
  "cases": [
    {"id": "case-1", "findings": {"phantom_tool": 1}}
  ]
}
```

Score it against the hidden truth:

```bash
./telltail mirror score \
  --truth hidden-truth.json \
  --submission worker-findings.json > mirror-score.json
```

The output reports required findings caught, false accusations, forbidden
findings raised, recall, precision, and an integrity score. Run multiple
evaluators against the same hidden truth when comparing evaluator quality.
Case IDs in the submission must correspond to the IDs in the truth set.

## Cloud Shift: Google Cloud Batch

The cloud backend uses the same trace format as local runs and records the
observable remote lifecycle: submission, queueing, provisioning, completion,
failure, and timing. You need an authenticated `gcloud` CLI, a Google Cloud
project and region, a runnable container image, and permission to submit the
Batch job under the selected service account.

### 1. Generate a Batch job specification

```bash
./telltail cloud gcp spec \
  --project gary-agent-yard \
  --region europe-west2 \
  --image europe-west2-docker.pkg.dev/PROJECT/telltail/worker:0.3 \
  --command 'telltail worker-entrypoint ...' \
  --service-account gary-batch-worker@gary-agent-yard.iam.gserviceaccount.com \
  --out telltail-batch-job.json
```

This writes the JSON Batch specification to `telltail-batch-job.json`. It is
useful for reviewing or reusing the exact container, command, resources,
retry settings, logging destination, and service-account choice before
submission.

### 2. Run and record a complete cloud shift

```bash
./telltail cloud gcp run \
  --project gary-agent-yard \
  --region europe-west2 \
  --image europe-west2-docker.pkg.dev/PROJECT/telltail/worker:0.3 \
  --command 'telltail worker-entrypoint ...' \
  --service-account gary-batch-worker@gary-agent-yard.iam.gserviceaccount.com \
  --spec telltail-batch-job.json \
  --trace telltail-cloud-trace.jsonl \
  --job-id telltail-shift-001
```

This submits the job, polls its remote state, and writes the cloud lifecycle
trace to `telltail-cloud-trace.jsonl`. On success it prints the trace path.
Use the trace verification and analysis commands above afterwards:

```bash
./telltail trace verify telltail-cloud-trace.jsonl
./telltail analyze \
  --scenario scenario.json \
  --trace telltail-cloud-trace.jsonl \
  --worker example-worker \
  --backend gcp-batch > cloud-result.json
```

`cloud gcp submit` and `cloud gcp describe` are also available when you need
lower-level control over submission or inspection. `cloud gcp run` is the
certification-oriented path because it records the remote lifecycle in a
Telltail trace.

Telltail does not embed cloud credentials: the `gcloud` CLI or workload
identity owns authentication. For a production Cloud Shift, the container
should write its hash-chained worker trace and final result to durable object
storage, while stdout and stderr go to Cloud Logging. Those worker-internal
events can then be merged with the lifecycle evidence; the current `cloud gcp
run` command itself records the Batch lifecycle and does not inspect the
container's private files.

## What Telltail currently detects

- unavailable or phantom tool calls, including repeated phantom-tool use;
- materially equivalent repeated failed actions;
- missing-path mistakes;
- no-progress and stuck-loop signals;
- multi-approach thrash without observable progress;
- mechanically provable false blockers;
- authority violations;
- success claims without acceptance;
- clean, recovered, messy, false-success, and useful-failure outcomes; and
- phase timing, cloud queue/provision timing, cost, and worker/tool usage.

The 0.3 line is intentionally a foundation, not a taxonomy cathedral. New
detectors should be evidence-backed and mutation-tested before being trusted.

## Build and test from a checkout

```bash
go test ./...
go vet ./...
go build ./cmd/telltail
```

The original Telltail 0.2 certification corpus remains outside this checkout
until it can be imported without silently changing its baseline. That is the
reason this repository describes itself as an additive 0.3 foundation.
