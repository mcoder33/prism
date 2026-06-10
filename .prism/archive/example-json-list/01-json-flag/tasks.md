# Tasks — 01 json-flag

## 1. Implementation
- [x] 1.1 add `--json` flag to the list command
- [x] 1.2 add json tags to `changeRow`
- [x] 1.3 `renderJSON` + branch in `runList` (empty → `[]`)

## 2. Tests
- [x] 2.1 `TestListJSON` — two changes, fields match the table
- [x] 2.2 `TestListJSONEmpty` — `[]`, exit 0

## 3. Checks
- [x] 3.1 `make test` green, table output unchanged (golden copy)
