# Integrating Telltail

Telltail is designed to sit **beside** an AI worker, not inside it.

The worker does the job. Telltail provides the scenario, records observable evidence, verifies that evidence, analyses behaviour and aggregates results over time.

That separation is deliberate: the thing being evaluated should not be the authority that decides whether it behaved correctly.

## Integration shape

The simplest integration looks like this:

```text
scenario
   |
   v
worker process  --->  observable events  --->  Telltail trace
   |                                           |
   v                                           v
real files / tools                         verify hash chain
                                               |
                                               v
                                           analyse
                                               |
                                               v
                                           result / dossier
```

A worker can be integrated at several levels.

## Level 1: process-only integration

Use `telltail local run` around an existing agent command.

```bash
telltail local run \
  --dir ./sacrificial-repo \
  --worker codex-example \
  --command 'your-agent-command' \
  --trace shift.jsonl
```

This gives you a basic shift lifecycle without changing the worker itself. It is useful for smoke tests, timeout behaviour and process-level comparison.

It is intentionally limited: Telltail cannot infer tool calls or file intent that were never emitted as observable events.

## Level 2: rich worker adapter

For serious behavioural evaluation, adapt the worker so it appends structured events to the trace exposed through the `TELLTAIL_TRACE` environment variable.

Useful events include:

- tool call and tool result;
- shell command and result;
- file read/write/mutation;
- model/provider selection;
- token and cost information;
- progress markers;
- worker claims such as “complete” or “blocked”;
- acceptance checks;
- permission or approval decisions; and
- externally visible errors.

The goal is **not** to capture private reasoning. Capture the actions and claims another system could actually observe.

## Level 3: scenario-aware integration

A scenario gives those events meaning.

A useful scenario defines:

- the visible job;
- available and unavailable tools;
- relevant paths;
- acceptance conditions;
- authority boundaries;
- known traps or hidden opportunities; and
- enough environmental information for a competent worker to discover the correct route.

This lets Telltail distinguish, for example, a legitimate retry from pointless repetition or a real blocker from a provably false one.

## Level 4: longitudinal evaluation

A single shift tells you what happened once.

A dossier tells you what a worker tends to do.

Run the same worker across multiple scenarios and repeated trials, then aggregate results with:

```bash
telltail dossier \
  --worker worker-name \
  --result result-1.json,result-2.json,result-3.json
```

This is the level at which Telltail becomes useful for model routing and autonomy decisions. Recurring phantom tools, repeated failure loops, false-success claims and expensive recovery patterns are much more informative than one memorable run.

## Integrating specific agent systems

Telltail does not require a particular worker architecture. The same pattern can wrap CLI agents, local orchestration systems, cloud workers or custom autonomous software.

For systems such as Codex CLI, Claude Code, Gemini CLI, OpenCode or a bespoke worker, prefer a small adapter layer that maps the worker's native events into Telltail events.

Keep the adapter boring:

```text
native event -> explicit mapping -> Telltail event
```

Avoid hiding policy inside the adapter. The adapter should describe what happened, not decide whether it was good.

## Cross-platform rule

Unless a scenario is explicitly testing one operating system, integration code should remain cross-platform.

Do not hard-code PowerShell, Windows drive letters, POSIX paths or shell semantics into ordinary tests. Put platform-specific launch behaviour behind an explicit boundary and keep scenario semantics portable.

A useful test question is:

> If the same worker and scenario ran on Windows, Linux and macOS, which differences are genuinely part of the thing being tested?

Only those differences should survive into the scenario or adapter.

## Cloud integration

The Google Cloud Batch backend records the remote lifecycle using the same evidence model as local shifts.

The useful attribution chain is:

```text
identity
  -> submit authority
  -> job-spec validity
  -> queue
  -> provision
  -> worker execution
  -> acceptance
```

Do not score a worker for stages it never reached.

A production cloud adapter should ideally persist:

- the worker's internal Telltail trace;
- the cloud lifecycle trace;
- final acceptance evidence; and
- immutable raw artifacts needed to reproduce the analysis.

Those evidence streams can then be joined with stable shift/job identifiers.

## Integrating the historical 0.2 corpus

The current repository is the 0.3 behavioural-lab foundation. An earlier 0.2 Telltail line contained a broader deterministic certification corpus covering routing, invariance, deltas, state sequences and mutations.

That corpus should be imported additively rather than silently replaced.

Recommended order:

1. preserve the historical fixtures and expected outcomes unchanged;
2. map old events into the current trace model where possible;
3. run old and new tests side-by-side;
4. treat any changed baseline as an explicit policy decision;
5. add mutation tests around every adapter that could alter meaning; and
6. only retire old paths once equivalent evidence exists in the current model.

The objective is not version archaeology. It is to keep previously proven behaviour while gaining the richer shift/trace/dossier model.

## What a good integration preserves

A strong Telltail integration has five properties:

1. **The worker cannot rewrite the grading rules.**
2. **Raw evidence is separable from derived judgement.**
3. **Acceptance is mechanical wherever possible.**
4. **Failures are attributed to the correct layer.**
5. **The same scenario can be rerun after a model, prompt, tool or policy change.**

If those properties hold, Telltail can answer a much more useful question than “did the demo work?”

It can answer: **what changed in the worker's behaviour, and do we have evidence that the change is actually better?**
