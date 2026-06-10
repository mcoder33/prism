# Detail — 01 json-flag

- Add `--json` bool flag to the `list` cobra command.

- Tag `changeRow` fields: `name`, `current`, `nodes`, `phase` (lowercase keys, stable API).

- In `runList`: after the scan, branch — flag set → `json.MarshalIndent` + `fmt.Println`,
  else the existing table renderer (untouched).

- Edge-cases: empty slice marshals to `[]` (not `null`) — initialize `rows := []changeRow{}`.

- Worked example: two changes, one current →
  `[{"name":"rate-limiter","current":true,"nodes":3,"phase":"apply"}, {"name":"auth-refactor",…}]`.

Open (minor): none.
