# Telltail

**Deterministic behavioural evaluation for autonomous AI workers.**

Current repository line: **0.3 Mirror Shift foundation**.

Telltail is a deterministic behavioural evaluation laboratory for autonomous AI workers.

It does not ask "did the model sound clever?" It gives workers believable jobs, records observable behaviour, verifies the job mechanically, and builds evidence about completion, repeated mistakes, phantom tools, stuckness, authority discipline, cost, latency and bottlenecks.

## Repository status

This repository currently contains the verified Telltail 0.3 behavioural-lab foundation: the CLI, hash-chained trace model, analyser, worker dossiers, local runner, Google Cloud Batch backend and Mirror Shift scorer.

Telltail 0.2 was the earlier deterministic certification harness. Its original routing/invariance/delta/state-sequence/mutation corpus is not yet imported into this repository, so 0.3 should be read as an additive foundation rather than a claim that the full 0.2 corpus has already been merged. See [`docs/INTEGRATION.md`](docs/INTEGRATION.md).

The code in this repository is verified with `go test ./...` and `go vet ./...`.

## Install

```bash
go install github.com/matthewjameswatkins1978-cyber/Telltail/cmd/telltail@latest
```

Or clone and build locally:

```bash
go build ./cmd/telltail
```

## 0.3 goals

1. Preserve Telltail 0.2's deterministic certification philosophy.
2. Add an append-only hash-chained worker trajectory.
3. Add forensic mistake detection and repeated-failure detection.
4. Add a job profiler that separates worker time from queue/provision/tool/test overhead.
5. Add Worker Dossiers aggregated across shifts.
6. Add **Mirror Shift**: agents build evaluators for AI-worker traces, then canonical Telltail scores those evaluators against hidden ground truth.
7. Make local and cloud execution equal citizens. Google Cloud Batch is the first cloud backend.

## Important boundary

Telltail evaluates **observable behaviour**, not hidden chain-of-thought. Evidence is tool calls, commands, results, file changes, acceptance state, permissions, timing, cost, tokens and externally visible claims.

## Build

```bash
go test ./...
go vet ./...
go build ./cmd/telltail
```

## Local shift

```bash
./telltail local run \
  --dir ./sacrificial-repo \
  --worker gemini-flash \
  --command 'your-agent-command-here' \
  --trace shift.jsonl
```

The generic local runner captures process-level events. Rich adapters should additionally append tool/model/file/progress events using the same event schema.

## Analyse a shift

```bash
./telltail trace verify shift.jsonl
./telltail analyze --scenario scenario.json --trace shift.jsonl --worker gemini-flash --backend local > result.json
./telltail dossier --worker gemini-flash --result result.json
```

## Cloud Shift: Google Cloud Batch

The cloud backend deliberately uses the same Telltail scenario and trace model as local runs.

Generate a Batch container-job spec:

```bash
./telltail cloud gcp spec \
  --project gary-agent-yard \
  --region europe-west2 \
  --image europe-west2-docker.pkg.dev/PROJECT/telltail/worker:0.3 \
  --command 'telltail worker-entrypoint ...' \
  --service-account gary-batch-worker@gary-agent-yard.iam.gserviceaccount.com \
  --out telltail-batch-job.json
```

Run a complete remote lifecycle and record it as a Telltail trace:

```bash
./telltail cloud gcp run \
  --project gary-agent-yard \
  --region europe-west2 \
  --image IMAGE \
  --command 'COMMAND' \
  --service-account gary-batch-worker@gary-agent-yard.iam.gserviceaccount.com \
  --spec telltail-batch-job.json \
  --trace telltail-cloud-trace.jsonl \
  --job-id telltail-mirror-001
```

`cloud gcp submit` and `cloud gcp describe` remain available as lower-level operations, but `run` is the certification path because it records the remote lifecycle.

Telltail intentionally does not embed cloud credentials. The `gcloud` CLI / workload identity owns authentication.

For a production Cloud Shift, the container should write the hash-chained JSONL trace and final result to durable object storage, while stdout/stderr go to Cloud Logging. Queue and provision events belong in the same final result so cloud overhead is not misattributed to the worker.

## Mirror Shift

```bash
./telltail mirror init --dir mirror-shift
```

Give `mirror-shift/JOB.md` plus public training traces to an agent and ask it to build an evaluator. Keep the real hidden trace corpus and truth outside the worker's workspace.

The evaluator emits:

```json
{"evaluator":"worker-x","cases":[{"id":"case-1","findings":{"phantom_tool":1}}]}
```

Then canonical Telltail scores it:

```bash
./telltail mirror score --truth hidden-truth.json --submission worker-findings.json
```

Run multiple builder agents. Cross-evaluation can be added by feeding each evaluator other builders' worker trajectories, but the canonical hidden scorer remains the authority. That prevents evaluators from grading themselves.

## Current detectors

- unavailable / phantom tool calls and repeated phantom-tool recidivism
- materially equivalent repeated failed actions
- missing-path mistakes
- no-progress / stuck-loop signal
- multi-approach thrash signal without observable progress
- mechanically provable false blockers when the hidden harness knows a route remains
- authority violations
- success claims without acceptance
- clean/recovered/messy/false-success and useful-failure outcome classes
- per-phase time, cloud queue/provision timing and cost accounting
- worker/tool usage aggregation

0.3 is intentionally a foundation, not a taxonomy cathedral. New detectors should be evidence-backed and mutation-tested before being trusted.
