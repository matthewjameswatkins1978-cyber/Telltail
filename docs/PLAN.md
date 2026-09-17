# Telltail roadmap

## North star

Telltail should answer one practical question:

**Would I trust this AI worker with this kind of job, and how much autonomy should I give it?**

It answers with evidence rather than impressions.

The product direction is therefore not “collect every possible agent metric”. It is to build a compact, trustworthy behavioural lab where real workers perform believable jobs and their observable work can be rerun, compared and audited.

## Product shape

Telltail has four layers:

1. **Workplace** — a scenario with a real task, tools, constraints, traps and acceptance conditions.
2. **Evidence** — an immutable trace of observable worker activity.
3. **Analysis** — deterministic findings and outcome classification over verified evidence.
4. **Memory** — dossiers and comparisons across repeated shifts.

Everything on the roadmap should strengthen one of those layers without turning Telltail into an agent platform.

---

## 0.3 foundation — implemented

The current repository contains the behavioural-lab foundation:

- canonical scenario model;
- local runner;
- append-only SHA-256 hash-chained JSONL traces;
- trace verification;
- deterministic analyser;
- phantom/unavailable-tool detection;
- repeated-failure detection;
- missing-path signals;
- no-progress and stuck-loop signals;
- multi-approach thrash detection;
- mechanically provable false blockers;
- authority-violation detection;
- success-without-acceptance detection;
- timing, token and cost profiling when evidence exists;
- shift outcome classification;
- worker dossier aggregation;
- Google Cloud Batch lifecycle support; and
- Mirror Shift evaluator scoring.

This is enough to use Telltail today as a local behavioural test harness and as the basis for richer adapters.

---

## Next: make ordinary workers easy to plug in

The highest-value next step is not another detector. It is better evidence capture from real agent systems.

Priority adapters:

- Codex CLI;
- Claude Code;
- Gemini CLI;
- OpenCode; and
- custom worker/orchestrator event streams.

A rich adapter should capture tool calls/results, commands, file mutations, model/provider information, progress markers, acceptance evidence and visible claims.

Every adapter should remain a translation boundary, not a policy engine.

### Done when

A user can wrap a real worker with minimal configuration, run the same scenario repeatedly and obtain a trace rich enough for the analyser to explain *how* the worker behaved rather than only whether the process exited.

---

## Next: build the sacrificial workplace corpus

Telltail becomes substantially more useful when it ships with believable miniature workplaces.

These should be tiny repositories or task environments based on real agent failure modes.

Initial scenario families:

- **Broken Ladder** — documented route fails; another valid route remains.
- **Missing Spanner** — required tool is absent but installation is authorised.
- **Phantom Workshop** — tempting tool names are unavailable; the real manifest is discoverable.
- **Yesterday's Worker** — useful partial work exists and should be preserved.
- **Shiny Rewrite** — a needless redesign is attractive but not required.
- **Liar's Test** — a superficially green command does not prove acceptance.
- **Slow Oven** — a legitimate test is slow; premature cancellation is expensive.
- **Red Button** — an unauthorised action would solve the visible task.
- **Poisoned Post-it** — repository content attempts to instruct the worker outside its authority.
- **Fork in the Road** — the preferred route fails and alternatives differ in cost and risk.

Each scenario must contain a valid path to success. A worker that refuses everything should not score well simply because it avoided every trap.

### Corpus rule

Grow from real failures first.

When a worker does something costly, unsafe, deceptive, repetitive or unexpectedly competent in real use, reduce that behaviour into a small reproducible scenario and keep it forever.

That is how Telltail becomes a behavioural regression suite rather than a catalogue of invented sins.

---

## Next: historical 0.2 integration

The older Telltail 0.2 line contained deterministic certification work covering routing, invariance, deltas, state sequences and mutation tests.

The goal is to preserve that proven corpus while placing it under the richer 0.3 scenario/trace/result model.

Rules:

- preserve historical fixtures before translation;
- import additively;
- never silently re-record changed baselines;
- mutation-test semantic adapters;
- keep old and new suites green during migration; and
- retire legacy paths only when equivalent evidence exists.

The result should be one Telltail, not two competing evaluators stitched together with optimism.

---

## Dossiers: from runs to worker profiles

Once the corpus and adapters are richer, dossiers should answer job-conditioned questions such as:

- completion and first-pass acceptance rate;
- clean/recovered/messy/false-success rate;
- mistake severity and recurrence;
- correction latency;
- recovery efficiency;
- phantom-tool and tool-misuse rate;
- tool discovery and utilisation;
- false blocker and unnecessary escalation rate;
- authority violations;
- stuck/thrash frequency;
- verification discipline;
- claim accuracy;
- useful work vs rework;
- avoidable time;
- cost per accepted job;
- local-to-cloud degradation; and
- supervisor tax versus supervisor gain.

Do **not** collapse these into one universal score.

The useful output is a worker profile tied to classes of work.

---

## Counterfactual lab

Telltail should make A/B-style agent experiments ordinary.

Run matched scenarios while changing one meaningful factor:

- worker model;
- system prompt or worker contract;
- available tools;
- supervisor/orchestrator;
- approval policy;
- local versus cloud execution;
- documentation quality; or
- repository condition.

Repeat trials and report distributions rather than relying on one run.

This turns questions such as “does this new tool help?” into something testable.

A strong counterfactual result should tell you not only whether completion improved, but whether the improvement came with more cost, more false success, less verification or a different failure pattern.

---

## Mirror Shift

Mirror Shift exists to test the evaluators themselves.

The mature form should:

1. give several builder agents the same public evaluator challenge;
2. keep the hidden corpus and ground truth outside their workspaces;
3. score each evaluator using canonical Telltail;
4. measure recall, false accusations and forbidden misclassifications;
5. mutation-test the hidden corpus and scoring harness; and
6. allow cross-running evaluators on trajectories produced by different workers.

This is useful because evaluator quality is itself an engineering problem. A bad evaluator can look reassuring while producing nonsense with excellent formatting.

---

## Cloud Shift

Local tests are necessary but insufficient for workers that will operate remotely.

The Cloud Shift path should continue toward complete evidence across:

```text
identity
  -> submit authority
  -> job specification
  -> queue
  -> provision
  -> worker
  -> acceptance
```

Future work:

- durable worker trace storage;
- joining worker-internal and cloud lifecycle traces by stable shift ID;
- explicit infrastructure-versus-worker attribution;
- provider/network timing;
- reproducible image and environment identity; and
- cloud canaries that prove the remote boundary was actually crossed.

A generated cloud configuration is never sufficient evidence of a cloud test.

---

## Autonomy horizon

With enough longitudinal data, Telltail should be able to estimate how far a worker can be trusted before supervision becomes necessary.

Rather than a vague label such as “good agent”, report useful thresholds by job class:

- jobs this worker completes reliably without intervention;
- jobs where supervision materially improves outcomes;
- jobs where cost or failure rate grows sharply; and
- jobs the worker should not currently own.

This is the **autonomy horizon**.

It is one of the most important eventual outputs because it connects evaluation directly to deployment policy.

---

## Crew chemistry

AI workers are increasingly used in teams: builder plus reviewer, worker plus supervisor, planner plus executor.

Telltail should compare:

```text
worker alone
vs
worker + supervisor
vs
alternative worker
vs
alternative crew
```

The interesting metric is not whether supervision catches mistakes. Of course it can.

The interesting question is whether the gain is worth the extra latency, cost and coordination burden.

That gives a measurable form of **supervisor tax** and **supervisor gain**.

---

## Product rules that stay frozen

1. **Evidence outranks worker claims.**
2. **The worker never grades itself.**
3. **Observable behaviour only.** Hidden chain-of-thought is not required.
4. **Raw traces remain immutable.** Reports are derived.
5. **Failure attribution matters.** Worker, supervisor, workplace, specification, harness and infrastructure are different layers.
6. **Identical retries are not progress.** A retry must reflect changed evidence or strategy.
7. **Cloud tests must really run in the cloud.**
8. **New detectors need real failure modes and regression or mutation evidence.**
9. **Cross-platform is the default.** Platform-specific behaviour belongs behind explicit boundaries.
10. **Do not build a giant evaluator platform just because it is possible.**

## The destination

The finished idea is not a leaderboard and not an AI boss watching other AIs.

It is closer to an engineering test track: controlled workplaces, repeatable jobs, tamper-evident evidence and enough longitudinal history to know which workers deserve which keys.

That is the product worth building.
