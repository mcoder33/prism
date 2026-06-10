# Proposal — example-json-list

## Why

- `prism list` prints a human table only; scripts and CI steps have to scrape it.
- Two internal pipelines already parse the output with `awk` and broke on the last format tweak.

## What

- Add a `--json` flag to `prism list` that emits the same data as a JSON array on stdout.

## Constraints & Invariants

- Default (no flag) output stays byte-for-byte identical — the human table is untouched.
- JSON mode prints nothing but valid JSON (no headers, no styling, no log lines on stdout).

## Decisions

- Marshal the existing list model; no new data source (strategy A in concept.md).
- One object per change: `{"name", "current", "nodes", "phase"}`.

## Non-goals

- No `--format` generalization (yaml/tsv) — until a second consumer asks.
- No JSON for other commands.
