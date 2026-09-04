# Integrating this 0.3 foundation into the existing Telltail 0.2 tree

The existing Telltail 0.2 checkout is the canonical product. This packet is deliberately additive because the live Windows checkout was not mounted in the build environment.

Recommended integration order:

1. Create a branch from the verified Telltail 0.2.0 baseline/tag.
2. Copy the 0.3 packages under temporary names or map them into the existing module namespace rather than replacing 0.2 packages wholesale.
3. Preserve every existing 0.2 test and corpus fixture unchanged on the first integration commit.
4. Add the 0.3 event/scenario/trace models and get all old + new tests green.
5. Adapt current GARY `telltail.adapter.v1` decision events into trajectory events.
6. Add `local run`, `cloud gcp run`, `analyze`, `dossier`, and `mirror` commands without breaking current CLI commands.
7. Run the original 0.2 corpus and hardening corpus. Any baseline change requires an explicit policy decision, not an automatic re-record.
8. Run the new local sacrificial scenario.
9. Run one real Google Cloud Batch shift and verify its trace hash chain plus acceptance result.
10. Tag only after Windows + Linux builds, `go test ./...`, `go vet ./...`, existing corpus, new behavioural tests and real cloud canary all pass.

Do not copy the temporary GARY cloud-canary workflow. The lasting implementation belongs in Telltail's own backend/CI once its repository/cloud identity is established.
