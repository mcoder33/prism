# Signatures — 01 json-flag

```go
// changeRow is the existing list model (CHANGED: json tags added, fields untouched).
type changeRow struct {
    Name    string `json:"name"`
    Current bool   `json:"current"`
    Nodes   int    `json:"nodes"`
    Phase   string `json:"phase"`
}

// runList renders the change list (CHANGED: branches on --json before the table renderer).
func runList(cmd *cobra.Command, args []string) error

// renderJSON marshals rows as an indented JSON array to w (NEW).
// Empty input prints "[]", never "null".
func renderJSON(w io.Writer, rows []changeRow) error
```
