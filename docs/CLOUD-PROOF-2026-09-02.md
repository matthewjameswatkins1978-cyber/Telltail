# Cloud Shift proof — 2 September 2026

A real Google Cloud Batch canary was run against the existing GARY cloud boundary in project `gary-agent-yard`, region `europe-west2`.

The proving sequence produced three useful observations:

1. A feature-branch GitHub Actions identity was rejected by the configured Workload Identity Federation attribute condition. This is a correct authority boundary and should be classified as an infrastructure/identity failure before worker execution.
2. The first master-branch Batch submission authenticated successfully but attempted to use the default Compute Engine service account. Google Cloud rejected it because the dispatcher was not authorised to act as that service account. GARY's actual Batch implementation uses the dedicated `gary-batch-worker@gary-agent-yard.iam.gserviceaccount.com` identity instead.
3. One subsequent canary job specification was rejected as invalid JSON because of a literal escaped newline introduced while generating the temporary workflow. This is a harness/job-spec failure, not a worker failure.

After correcting those issues, the canary authenticated, submitted a real Batch job under the authorised GARY worker service account, waited for the remote task, and completed successfully.

Successful GitHub Actions run: `33643978741`.
Head commit for the successful proof: `85063d584f4dcb1f452e5a9d0bc64335b9f4c4cc`.

The temporary master canary workflow was removed after proof. The feature-branch proving PR was closed without merge.

These failures are worth retaining conceptually as Telltail regression cases because they demonstrate the attribution boundary:

`identity -> submit authority -> job-spec validity -> queue -> provision -> worker -> acceptance`

A worker should only be scored for stages it actually reached.
