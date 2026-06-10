# Spec — 01 json-flag

## Requirement: JSON output

`prism list --json` SHALL print all active changes as a JSON array on stdout and nothing else.

### Scenario: lists changes as JSON
- WHEN `prism list --json` runs in a repo with two active changes
- THEN stdout is a JSON array of two objects with `name`, `current`, `nodes`, `phase`

### Scenario: empty repo yields empty array
- WHEN `prism list --json` runs with no `.prism/` changes
- THEN stdout is `[]` and the exit code is 0

## Coverage
(appended by prism:apply as tasks are checked off; read and updated by prism:verify)
- Scenario: lists changes as JSON → test `internal/cli/list_test.go::TestListJSON`
- Scenario: empty repo yields empty array → test `internal/cli/list_test.go::TestListJSONEmpty`
