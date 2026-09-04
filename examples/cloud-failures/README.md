# Cloud failure attribution examples

These are scenario ideas for Telltail's Cloud Shift corpus.

A cloud run is not one undifferentiated event. Attribute failures to the stage that actually failed:

```text
identity -> submit authority -> job-spec validity -> queue -> provision -> worker -> acceptance
```

Examples captured during the September 2026 canary:

- feature-branch identity rejected by WIF attribute policy;
- dispatcher lacked `actAs` permission for the wrong/default worker service account;
- malformed Batch job JSON rejected before scheduling;
- corrected authorised submission reached and completed a remote Batch task.

Do not score the worker for stages it never reached.
