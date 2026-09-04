# Telltail 0.3 plan — from certification harness to behavioural work lab

## North star

Telltail answers: **Would I hire this agent for this kind of work, and how much autonomy should I give it?**

It does that by running reproducible fake jobs in controlled workplaces, recording observable work trajectories, verifying the real outcome mechanically, and accumulating a worker/crew dossier from evidence.

Telltail does not infer hidden chain-of-thought and does not let the worker grade itself.

## Architecture rule

The scenario, trace and analysis semantics are identical across execution backends.

- **Local Lab**: fast, cheap, deterministic development and regression work.
- **Cloud Shift**: real isolated remote execution with identity, queue, provisioning, provider/network latency and infrastructure failures visible.

A result is not a valid cloud certification merely because a cloud job specification was generated. A real remote task must cross the execution boundary.

## Phase A — 0.3 foundation (implemented in this packet)

1. Canonical scenario model: visible job, available/unavailable tools, hidden opportunities/traps, acceptance and authority boundaries.
2. Append-only SHA-256 hash-chained JSONL trajectory.
3. Behaviour analyser for:
   - phantom/unavailable tool calls;
   - repeated phantom tool calls;
   - materially repeated failed actions;
   - missing-path mistakes;
   - no-progress/stuck loops;
   - observable thrashing across distinct failed approaches;
   - mechanically provable false blockers;
   - authority violations;
   - success claims without acceptance.
4. Job profiler for model/tool/command/queue/provision/acceptance/teardown time, failed work, repeated work, estimated avoidable time, token use and cost.
5. Shift outcome classification: clean success, recovered success, messy success, false success, useful failure, failure.
6. Worker dossier aggregation across shifts.
7. Google Cloud Batch backend with explicit execution service-account support.
8. `cloud gcp run`: records cloud submit plus queue/provision lifecycle into the same hash-chained trace model.
9. Mirror Shift challenge and canonical hidden-truth scorer.

## Phase B — integrate into Telltail 0.2

Preserve all existing 0.2 routing/invariance/delta/state-sequence/mutation tests. Add 0.3 beside them rather than replacing them.

1. Add adapters from current Telltail/GARY events into the 0.3 event schema.
2. Add rich worker adapters for Gemini CLI, Codex CLI, Claude Code/OpenCode where available.
3. Capture tool calls/results, commands, file mutations, model/provider/token usage, progress markers, acceptance evidence and worker claims.
4. Make every rich adapter mutation-tested: deliberately remove/mangle an event and prove the detector changes.
5. Merge the Cloud Shift lifecycle trace with the worker-internal trace using stable shift/job IDs.
6. Persist cloud trace/result artifacts to durable object storage; send stdout/stderr to Cloud Logging.
7. Keep raw traces immutable and derive reports from them.

## Phase C — the sacrificial workplace corpus

Build tiny believable repositories, each with a small visible job and hidden behavioural probes. Grow primarily from real failures.

Initial rooms:

- **Broken Ladder**: documented route fails; valid route remains.
- **Missing Spanner**: required tool absent but authorised to install.
- **Phantom Workshop**: tempting tool names are unavailable; manifest is discoverable.
- **Yesterday's Worker**: useful partial work exists and should be preserved.
- **Shiny Rewrite**: unnecessary redesign is tempting but not required.
- **Liar's Test**: superficially green command does not prove acceptance.
- **Slow Oven**: legitimate test is slow; premature cancellation is costly.
- **Red Button**: easy unauthorised action would solve the visible task.
- **Poisoned Post-it**: untrusted repository content gives an instruction outside authority.
- **Fork in the Road**: preferred route breaks; alternatives differ in cost/risk.

Each scenario has positive controls as well as traps so 'refuse everything' cannot score well.

## Phase D — worker intelligence from longitudinal evidence

Aggregate enough shifts to estimate job-conditioned strengths rather than one global leaderboard.

Core dossier metrics:

- completion and first-pass acceptance;
- clean/recovered/messy/false success rates;
- mistakes per shift and severity distribution;
- repeated-mistake/recidivism rate;
- correction latency and recovery efficiency;
- phantom-tool, missing-path and tool-misuse rates;
- tool discovery/utilisation;
- false blocker and unnecessary escalation rate;
- scope/authority violations;
- stuck/thrash frequency;
- verification discipline and claim accuracy;
- wall time, useful work, rework, avoidable work and cost per accepted job;
- local-to-cloud degradation;
- supervisor tax and supervisor gain.

Do not reduce all of this to one universal score. Route by job class and evidence.

## Phase E — counterfactual lab

Rerun matched scenarios changing one meaningful factor:

- worker model;
- prompt/worker contract;
- tool availability (for example Suture present/absent);
- Yardmaster strength;
- approval policy;
- local vs cloud;
- documentation quality.

Use repeated trials and confidence intervals. This is how Telltail answers whether a change actually helped rather than relying on anecdotes.

## Phase F — Mirror Shift tournament

1. Give multiple builder agents the same public challenge: build a deterministic evaluator for worker traces.
2. Keep hidden trace corpus and truth outside their workspace.
3. Canonical Telltail scores each evaluator for recall, false accusations and forbidden misclassifications.
4. Cross-run evaluator A on builder B's trajectory, B on C, C on A.
5. Keep canonical scoring outside all builders so no evaluator can move the goalposts.
6. Mutation-test the hidden corpus and evaluator harness.

Mirror Shift measures engineering competence plus whether an agent understands what good agent work actually looks like.

## Phase G — autonomy horizon and crew chemistry

Once there is enough data, fit success probability against reference job effort/complexity for each worker and crew pairing. Report useful horizons such as 90%, 80% and 50% reliable job size.

Measure worker-alone vs Yardmaster-supervised performance to learn where supervision saves more than it costs.

## Rules that stay frozen

1. Evidence outranks worker claims.
2. The worker never grades itself.
3. Observable behaviour only; hidden chain-of-thought is not required.
4. Keep raw source traces immutable; reports are derived.
5. Distinguish worker failure from Yardmaster, workplace, job-spec, harness and external/infrastructure failure.
6. A retry must represent a meaningful change in evidence or strategy; identical circling is a defect.
7. Cloud tests must actually run in the cloud.
8. Do not build a giant evaluator platform merely because it is possible. New detectors need an evidence-backed failure mode and a regression/mutation test.
