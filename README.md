# Telltail

**A test lab for AI workers.**

Give an AI agent a real job. Record what it actually does. Verify the evidence. Find out whether it is competent, wasteful, stuck, unsafe, or merely *saying* it finished.

Telltail is a deterministic behavioural evaluation toolkit for autonomous and semi-autonomous AI workers. It turns agent runs into inspectable evidence instead of anecdotes.

It is useful when the question is not simply:

> Can this model answer the prompt?

but:

> Can I trust this worker to do this kind of work, with these tools and permissions, without wasting time or quietly getting it wrong?

Telltail does **not** read hidden chain-of-thought and does not ask an LLM to grade another LLM. It evaluates observable work: commands, tool calls, results, file changes, progress, acceptance evidence, permissions, timing, cost and claims.

Current line: **0.3.0-mirror**.

---

## Why Telltail exists

AI coding agents and autonomous workers are increasingly capable of doing substantial jobs, but ordinary pass/fail testing misses a large part of the problem.

Two agents can both produce a passing result while behaving very differently:

- one discovers the environment, uses the available tools and verifies its work;
- one calls tools that do not exist, repeats failed commands, wanders through several approaches, edits outside its authority and finally declares success because something looked plausible.

A normal test suite may only see the final green check.

Telltail watches the **work trajectory** as well as the outcome.

That makes it useful for agent development, regression testing, model comparison, prompt and tool evaluation, autonomy decisions, supervision experiments and post-mortems on difficult agent runs.

---

## The idea in one line

```text
scenario -> worker shift -> evidence trace -> verify -> analyse -> dossier
```

A **shift** is one worker attempting one job under one defined set of conditions.

Telltail gives that shift structure:

1. A **scenario** describes the job, available tools, traps, acceptance conditions and authority boundaries.
2. A **runner** launches the worker locally or through Google Cloud Batch.
3. Observable events are written to a JSONL **trace**.
4. The trace is protected by an append-only SHA-256 hash chain.
5. Telltail verifies the evidence before analysing it.
6. The **analyser** compares the trace with the scenario and produces a structured result.
7. Several results can be combined into a **worker dossier** so repeated behavioural patterns become visible.

The worker never gets to mark its own homework.

---

## What Telltail can tell you

Telltail currently detects or reports evidence for problems including:

- unavailable or phantom tool calls;
- repeated phantom-tool use;
- materially repeated failed actions;
- missing-path mistakes;
- no-progress and stuck-loop behaviour;
- thrashing across multiple failed approaches;
- mechanically provable false blockers;
- authority violations;
- success claims without acceptance evidence;
- clean, recovered, messy, false-success and useful-failure outcomes;
- model, tool, command and acceptance timing;
- cloud queue and provisioning timing;
- token and cost totals when supplied by the trace; and
- recurring patterns across multiple shifts.

The point is not to create a giant score for “which AI is best”. The useful question is usually narrower:

**Which worker is reliable for this kind of job, under these conditions, at this level of autonomy?**

---

## What do you use it for?

### Test an AI coding agent before trusting it

Create a small sacrificial repository and give the agent a believable maintenance job. Include the same sort of awkwardness it will meet in the real world: a missing tool, stale documentation, a tempting but forbidden shortcut, partial work from yesterday, or a slow test.

Telltail records whether the agent solves the problem cleanly or reaches the answer by driving through three hedges and a greenhouse.

### Regression-test an agent workflow

When you change a system prompt, model, toolset, permissions, orchestration layer or supervisor, rerun the same scenarios.

Instead of saying “it feels better”, compare evidence:

- Did completion improve?
- Did false success fall?
- Did tool misuse increase?
- Did the worker get faster only because it stopped verifying?
- Did a new tool remove wasted work or merely create a new failure mode?

### Compare models or agent configurations

Run matched jobs with different workers and build dossiers from repeated shifts.

Telltail is deliberately better suited to **job-conditioned comparison** than a universal leaderboard. A worker may be excellent at bounded repository repair and poor at tool discovery; another may show the opposite pattern.

### Measure how much autonomy a worker deserves

A worker that succeeds only when heavily supervised is different from one that succeeds cleanly on its own.

Telltail provides the evidence needed to explore questions such as:

- Which tasks can this worker safely own end-to-end?
- Where should approval gates remain?
- When does supervision improve the result enough to justify its cost?
- At what job complexity does reliability start to collapse?

### Diagnose agent failures

Telltail separates stages of the system so a worker is not blamed for a broken harness, identity failure, cloud provisioning problem or invalid job specification.

For cloud work the useful boundary is roughly:

```text
identity -> submit authority -> job spec -> queue -> provision -> worker -> acceptance
```

That distinction matters. “The agent failed” is often the least informative sentence in the room.

### Test the evaluator itself

**Mirror Shift** turns evaluation into its own challenge. An agent builds an evaluator from public examples, then canonical Telltail scores it against hidden truth.

That allows you to test whether a proposed evaluator can recognise real worker defects without inventing accusations or moving the goalposts.

---

## Telltail is not

Telltail is intentionally small in scope.

It is **not**:

- an agent framework;
- an orchestrator;
- a replacement for Codex, Claude Code, Gemini CLI, OpenCode or another worker;
- an LLM-as-judge service;
- a chain-of-thought inspector;
- a policy learner; or
- a universal model leaderboard.

It sits beside the worker and provides the measuring equipment.

Think less “another agent platform” and more **crash-test rig for autonomous software workers**.

---

## Install

With Go installed:

```bash
go install github.com/matthewjameswatkins1978-cyber/Telltail/cmd/telltail@latest
```

Or build from a checkout:

```bash
go build -o telltail ./cmd/telltail
```

On Windows the binary is normally `telltail.exe`; on macOS and Linux it is `telltail`.

Check the CLI:

```bash
telltail version
```

Current commands:

```text
trace verify
analyze
dossier
local run
cloud gcp spec
cloud gcp submit
cloud gcp describe
cloud gcp run
mirror init
mirror score
```

Telltail is designed to be cross-platform. Shell-specific behaviour belongs at an explicit boundary rather than being baked into ordinary tests or scenarios.

---

## Five-minute local example

The repository includes a sample scenario at:

```text
examples/scenarios/tool-friction.json
```

It describes a small job with available and unavailable tools plus acceptance conditions.

### 1. Run a worker

Use a disposable or sacrificial working directory:

```bash
telltail local run \
  --dir ./sacrificial-repo \
  --worker example-worker \
  --command 'your-agent-command-here' \
  --trace shift.jsonl
```

If the command needs a particular shell, use `--shell bash`, `--shell pwsh`, or `--shell cmd.exe`.

The runner records the basic shift lifecycle. Rich worker adapters can append additional events through the `TELLTAIL_TRACE` environment variable so tool calls, commands, model usage, file activity and progress become part of the same evidence stream.

### 2. Verify the evidence

```bash
telltail trace verify shift.jsonl
```

A valid trace reports the event count and final hash:

```text
OK 12 events head=<sha256-hash>
```

If the JSON is malformed, event ordering is broken, or the hash chain no longer matches, verification fails.

### 3. Analyse the shift

```bash
telltail analyze \
  --scenario examples/scenarios/tool-friction.json \
  --trace shift.jsonl \
  --worker example-worker \
  --backend local > result.json
```

The resulting `ShiftResult` contains the outcome, findings, timing, cost/token information when available, and the verified trace head.

### 4. Build a dossier

After several runs:

```bash
telltail dossier \
  --worker example-worker \
  --result result-1.json,result-2.json,result-3.json > dossier.json
```

A dossier reveals patterns that a single lucky or unlucky run cannot: repeated mistakes, false-success rate, completion rate, stuck loops, average duration, cost and finding counts.

---

## Scenario design

A useful Telltail scenario is a small believable workplace, not a synthetic riddle.

Good scenarios contain:

- a visible job with a mechanically checkable outcome;
- enough tooling information for a competent worker to discover what is possible;
- realistic friction or ambiguity;
- explicit authority boundaries;
- positive routes to success, so refusing everything cannot score well; and
- one or more behavioural probes based on failures you actually care about.

Examples of useful scenario families include:

- **Broken Ladder**: the documented route fails but a valid route remains.
- **Missing Spanner**: a required tool is absent but installation is authorised.
- **Phantom Workshop**: tempting tool names are unavailable while the real manifest is discoverable.
- **Yesterday's Worker**: useful partial work already exists and should be preserved.
- **Shiny Rewrite**: a needless redesign is tempting.
- **Liar's Test**: a superficially green command is not sufficient acceptance evidence.
- **Slow Oven**: a legitimate test is slow and premature cancellation is costly.
- **Red Button**: an unauthorised shortcut would solve the visible problem.
- **Poisoned Post-it**: repository content contains an instruction outside the worker's authority.

The strongest Telltail corpus should grow from real failures rather than from an abstract taxonomy of every mistake an agent might conceivably make.

---

## Mirror Shift

Mirror Shift tests evaluators rather than ordinary workers.

Create a challenge directory:

```bash
telltail mirror init --dir mirror-shift
```

This creates public challenge instructions and an example truth format. The real hidden corpus should remain outside the evaluator's workspace.

Score a submission:

```bash
telltail mirror score \
  --truth hidden-truth.json \
  --submission worker-findings.json > mirror-score.json
```

The canonical scorer reports required findings caught, false accusations, forbidden findings, precision, recall and integrity information.

The evaluator can propose findings. It cannot redefine the truth it is being scored against.

---

## Cloud Shift

Telltail can run the same behavioural model across a real Google Cloud Batch execution boundary.

That matters because local success does not prove cloud behaviour. Remote jobs introduce identity, submission authority, queueing, provisioning, provider latency and infrastructure failure modes that should be visible rather than folded into a vague “agent failed”.

Generate a Batch specification:

```bash
telltail cloud gcp spec \
  --project PROJECT \
  --region europe-west2 \
  --image IMAGE \
  --command 'worker-entrypoint ...' \
  --service-account SERVICE_ACCOUNT \
  --out telltail-batch-job.json
```

Run and record a cloud shift:

```bash
telltail cloud gcp run \
  --project PROJECT \
  --region europe-west2 \
  --image IMAGE \
  --command 'worker-entrypoint ...' \
  --service-account SERVICE_ACCOUNT \
  --trace telltail-cloud-trace.jsonl
```

Then verify and analyse that trace exactly as you would a local one.

Telltail does not embed cloud credentials. Authentication remains the responsibility of `gcloud` or the surrounding workload identity.

A real Cloud Shift proof is documented in [`docs/CLOUD-PROOF-2026-09-02.md`](docs/CLOUD-PROOF-2026-09-02.md).

---

## Evidence model

Telltail's core rule is simple:

**Evidence outranks claims.**

A worker saying “done” is an event. It is not acceptance.

Raw traces are intended to be immutable evidence. Reports and dossiers are derived from those traces. This separation makes it possible to improve analysers later without rewriting what actually happened during the shift.

The current trace format is JSONL with a SHA-256 predecessor/hash chain. Verification therefore catches accidental or deliberate mutation of the recorded sequence before analysis begins.

---

## Where the project is now

The current repository contains the **0.3 Mirror Shift foundation**:

- scenario model;
- local runner;
- hash-chained trace format and verifier;
- deterministic behavioural analyser;
- worker dossiers;
- Google Cloud Batch lifecycle backend; and
- Mirror Shift challenge/scorer.

An earlier Telltail 0.2 line contained a broader deterministic certification corpus covering routing, invariance, deltas, state sequences and mutations. That original corpus has not yet been imported into this repository, so 0.3 should be read as an additive behavioural-lab foundation rather than a claim that every historical Telltail test is already present here.

See [`docs/INTEGRATION.md`](docs/INTEGRATION.md) for the integration boundary and [`docs/PLAN.md`](docs/PLAN.md) for the roadmap.

---

## Build and test

```bash
go test ./...
go vet ./...
go build ./cmd/telltail
```

New detectors should not be added because they sound clever. They should correspond to an observable failure mode and arrive with regression or mutation evidence proving that the detector changes when the relevant behaviour changes.

---

## Design principles

1. **Evidence outranks worker claims.**
2. **The worker never grades itself.**
3. **Observable behaviour is enough.** Hidden chain-of-thought is neither required nor desired.
4. **Raw evidence stays immutable.** Analysis is derived.
5. **Attribute failure to the correct layer.** Worker, supervisor, workplace, harness, specification and infrastructure are not the same thing.
6. **Retries must mean something changed.** Identical circling is a defect, not perseverance.
7. **Cloud certification must cross a real cloud boundary.** A generated job spec is not a cloud test.
8. **Do not build a taxonomy cathedral.** Add detectors because reality produced a failure worth catching.
9. **Prefer job-conditioned evidence to universal scores.** Different workers can be good at different kinds of work.
10. **Keep the measuring equipment independent of the thing being measured.**

Telltail exists to make autonomous work less mysterious. Not by asking an AI whether another AI seemed sensible, but by leaving tracks in the snow and measuring where they actually went.
