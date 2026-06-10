# 01 — json-flag

- **What:** `prism list --json` prints the change list as a JSON array.

- **Logic:** tag the existing `changeRow` slice for JSON; on the flag, marshal it instead of rendering the table.

- **Guarantees:** default output byte-identical; JSON mode emits only valid JSON on stdout.

- **Input → output:** `[]changeRow` → JSON array (or the unchanged table).

**Status:** ✅

**Open:** none
