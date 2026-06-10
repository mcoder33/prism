# Concept — example-json-list

## Best practices

> Skipped — user opted out (routine flag addition).

## Candidate strategies

- **A. Marshal the existing list model** — RECOMMENDED — the table already renders from a
  `[]changeRow`; tag it for JSON and marshal the same slice. One source of truth, no drift
  between table and JSON.
- **B. Separate JSON walker** — re-scan `.prism/` independently for JSON — duplicate logic, drift risk.
- **C. `--format` enum now** — generalize early — speculative, no second format requested.

## Chosen strategy

A — because the slice the table renders from is already the complete data set.
Rejected: B duplicates the scan · C is speculation.

## Data flow

`.prism/` dirs → scan (existing) → `[]changeRow` → branch: table renderer (unchanged) | `json.Marshal` → stdout.
