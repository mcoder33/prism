// Package change reads a prism change from disk — README state table, nodes, tasks, open
// items, specs, diagrams — and reports deterministic findings. It backs `prism status` and
// `prism check`: the mechanical cross-checks the agent used to do by reading Markdown.
package change

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// PrismDir is the project-relative directory holding all changes.
const PrismDir = ".prism"

// Glyphs is the full status vocabulary, in lifecycle order.
var Glyphs = []string{"⚪", "🟡", "🟢", "🔵", "✅", "⏸"}

// Severity of a finding.
type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

// Finding is one deterministic observation: the envelope shared by status and check.
type Finding struct {
	Severity Severity `json:"severity"`
	Code     string   `json:"code"`
	Path     string   `json:"path,omitempty"`
	Message  string   `json:"message"`
}

// Count is a done/total pair (tasks checkboxes).
type Count struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

// OpenCounts tallies a node's Open items by tag.
type OpenCounts struct {
	Blocking int `json:"blocking"`
	Minor    int `json:"minor"`
	Deferred int `json:"deferred"`
	Untagged int `json:"untagged"`
}

// Node is one row of the README status table joined with its directory.
type Node struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Status     string     `json:"status"`
	NodeStatus string     `json:"nodeStatus,omitempty"`
	Dir        string     `json:"dir,omitempty"`
	Open       OpenCounts `json:"open"`
	Tasks      Count      `json:"tasks"`
	Artifacts  []string   `json:"artifacts"`
	Escalated  bool       `json:"escalated,omitempty"`
}

// Change is the parsed state of one change directory.
type Change struct {
	Name      string    `json:"change"`
	Dir       string    `json:"dir"`
	Current   bool      `json:"current"`
	Phase     string    `json:"phase"`
	Tier      string    `json:"tier"`
	Branch    string    `json:"branch,omitempty"`
	GitBranch string    `json:"gitBranch,omitempty"`
	Nodes     []Node    `json:"nodes"`
	RootTasks *Count    `json:"rootTasks,omitempty"`
	Findings  []Finding `json:"findings"`
}

var (
	phaseLine  = regexp.MustCompile(`^\*\*Phase:\*\*(.*)$`)
	boldWord   = regexp.MustCompile(`\*\*([a-z]+)\*\*`)
	tierLine   = regexp.MustCompile(`^\*\*Tier:\*\*\s*([a-z]+)`)
	branchLine = regexp.MustCompile("^\\*\\*Branch:\\*\\*\\s*`?([^`\\s]+)`?")
	nodeDirRe  = regexp.MustCompile(`^\d+[a-z]*-`)
	checkbox   = regexp.MustCompile(`^\s*[-*]\s+\[([ xX])\]`)
)

// Current returns the slug in .prism/CURRENT ("" if absent).
func Current(projectRoot string) string {
	b, err := os.ReadFile(filepath.Join(projectRoot, PrismDir, "CURRENT"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// Active lists active change slugs (sub-dirs of .prism/ except archive/), sorted.
func Active(projectRoot string) []string {
	entries, err := os.ReadDir(filepath.Join(projectRoot, PrismDir))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

// ErrNoChange means no change could be resolved without guessing.
var ErrNoChange = errors.New("no change selected")

// Resolve picks the change to report on: explicit name → CURRENT → the only active change.
// It never guesses between several active changes. A name may be a path under .prism/,
// e.g. "archive/<change>".
func Resolve(projectRoot, explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if cur := Current(projectRoot); cur != "" {
		return cur, nil
	}
	active := Active(projectRoot)
	switch len(active) {
	case 0:
		return "", fmt.Errorf("%w: no active changes in .prism/", ErrNoChange)
	case 1:
		return active[0], nil
	default:
		return "", fmt.Errorf("%w: .prism/CURRENT is empty and there are %d active changes (%s) — pass --change",
			ErrNoChange, len(active), strings.Join(active, ", "))
	}
}

// Load parses the change named name (relative to .prism/) and computes its findings.
func Load(projectRoot, name string) (*Change, error) {
	dir := filepath.Join(projectRoot, PrismDir, filepath.FromSlash(name))
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		msg := fmt.Sprintf("no change directory .prism/%s/", name)
		if Current(projectRoot) == name {
			msg += " — .prism/CURRENT is stale; repair it with the use command"
		}
		return nil, errors.New(msg)
	}
	c := &Change{
		Name:    name,
		Dir:     filepath.ToSlash(filepath.Join(PrismDir, name)),
		Current: Current(projectRoot) == name,
		Tier:    "standard",
		Nodes:   []Node{},
	}
	c.GitBranch = gitBranch(projectRoot)
	if err := c.readREADME(dir); err != nil {
		c.add(Error, "readme_missing", "README.md", "README.md not found — the change has no state table yet")
	}
	c.joinNodes(dir)
	if n, ok := countTasks(filepath.Join(dir, "tasks.md")); ok {
		c.RootTasks = &n
	}
	// An archived change is history: its recorded branch is expected to differ.
	archived := strings.HasPrefix(filepath.ToSlash(name), "archive/")
	if !archived && c.Branch != "" && c.GitBranch != "" && c.Branch != c.GitBranch {
		c.add(Warning, "branch_mismatch", "README.md",
			fmt.Sprintf("on branch %s, change expects %s — switch branches yourself", c.GitBranch, c.Branch))
	}
	return c, nil
}

func (c *Change) add(sev Severity, code, path, msg string) {
	c.Findings = append(c.Findings, Finding{Severity: sev, Code: code, Path: path, Message: msg})
}

// Errors counts error-severity findings.
func (c *Change) Errors() int {
	n := 0
	for _, f := range c.Findings {
		if f.Severity == Error {
			n++
		}
	}
	return n
}

func (c *Change) readREADME(dir string) error {
	f, err := os.Open(filepath.Join(dir, "README.md"))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	var header []string
	inTable := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if m := phaseLine.FindStringSubmatch(line); m != nil {
			if b := boldWord.FindStringSubmatch(m[1]); b != nil {
				c.Phase = b[1]
			}
			continue
		}
		if m := tierLine.FindStringSubmatch(line); m != nil {
			c.Tier = m[1]
			continue
		}
		if m := branchLine.FindStringSubmatch(line); m != nil {
			c.Branch = m[1]
			continue
		}
		if !strings.HasPrefix(line, "|") {
			inTable = false
			continue
		}
		cells := splitRow(line)
		if !inTable {
			if len(cells) > 0 && strings.EqualFold(cells[0], "Node") {
				header = cells
				inTable = true
			}
			continue
		}
		if isSeparator(cells) {
			continue
		}
		c.Nodes = append(c.Nodes, rowNode(header, cells))
	}
	return sc.Err()
}

func splitRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isSeparator(cells []string) bool {
	for _, c := range cells {
		if strings.Trim(c, "-: ") != "" {
			return false
		}
	}
	return true
}

func rowNode(header, cells []string) Node {
	n := Node{Artifacts: []string{}}
	for i, h := range header {
		if i >= len(cells) {
			break
		}
		switch strings.ToLower(h) {
		case "node":
			n.ID = strings.Trim(cells[i], "`")
		case "title":
			n.Title = cells[i]
		case "status":
			n.Status = glyphIn(cells[i])
		}
	}
	return n
}

func glyphIn(s string) string {
	for _, g := range Glyphs {
		if strings.Contains(s, g) {
			return g
		}
	}
	return ""
}

// nodeDirs maps node basenames (NN-name, NNa-name) to their change-relative paths.
func nodeDirs(dir string) map[string]string {
	out := map[string]string{}
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || p == dir {
			return nil
		}
		if !nodeDirRe.MatchString(d.Name()) {
			return filepath.SkipDir
		}
		rel, _ := filepath.Rel(dir, p)
		out[d.Name()] = filepath.ToSlash(rel)
		return nil
	})
	return out
}

func (c *Change) joinNodes(dir string) {
	dirs := nodeDirs(dir)
	seen := map[string]bool{}
	for i := range c.Nodes {
		n := &c.Nodes[i]
		key := n.ID
		if strings.Contains(key, "/") {
			key = key[strings.LastIndex(key, "/")+1:]
		}
		rel, ok := dirs[key]
		if !ok {
			c.add(Error, "row_without_dir", "README.md",
				fmt.Sprintf("table row %s has no node directory", n.ID))
			continue
		}
		seen[key] = true
		n.Dir = rel
		c.readNode(dir, n)
	}
	var orphans []string
	for name, rel := range dirs {
		if !seen[name] {
			orphans = append(orphans, rel)
		}
	}
	sort.Strings(orphans)
	for _, rel := range orphans {
		c.add(Error, "dir_without_row", rel, fmt.Sprintf("node directory %s/ has no row in README.md", rel))
	}
}

func (c *Change) readNode(changeDir string, n *Node) {
	ndir := filepath.Join(changeDir, filepath.FromSlash(n.Dir))
	entries, _ := os.ReadDir(ndir)
	for _, e := range entries {
		if !e.IsDir() {
			n.Artifacts = append(n.Artifacts, e.Name())
		}
	}
	rel := func(f string) string { return n.Dir + "/" + f }

	nodeMD, err := os.ReadFile(filepath.Join(ndir, "node.md"))
	if err != nil {
		c.add(Error, "missing_artifact", rel("node.md"), "node.md is missing (mandatory for every node)")
	} else {
		text := string(nodeMD)
		n.NodeStatus = statusOf(text)
		n.Open, n.Escalated = openItems(text)
		if n.NodeStatus != "" && n.Status != "" && n.NodeStatus != n.Status {
			c.add(Error, "status_mismatch", rel("node.md"),
				fmt.Sprintf("node.md says %s but the README table says %s", n.NodeStatus, n.Status))
		}
		if n.Open.Untagged > 0 {
			c.add(Error, "open_untagged", rel("node.md"),
				fmt.Sprintf("%d Open item(s) without a [blocking]/[minor]/[deferred: …] tag", n.Open.Untagged))
		}
		if n.Open.Blocking > 0 && (n.Status == "🟢" || n.Status == "🔵" || n.Status == "✅") {
			c.add(Warning, "blocking_open_past_drill", rel("node.md"),
				fmt.Sprintf("node is %s but still has %d [blocking] Open item(s)", n.Status, n.Open.Blocking))
		}
		for _, label := range []string{"**What:**", "**Status:**"} {
			if !strings.Contains(text, label) {
				c.add(Warning, "node_format", rel("node.md"), "node.md lacks the "+label+" line")
			}
		}
	}

	tasks, hasTasks := countTasks(filepath.Join(ndir, "tasks.md"))
	n.Tasks = tasks
	switch n.Status {
	case "🟢", "🔵", "✅":
		if !hasTasks {
			c.add(Error, "missing_artifact", rel("tasks.md"),
				fmt.Sprintf("node is %s but tasks.md is missing", n.Status))
		}
	}
	if (n.Status == "🔵" || n.Status == "✅") && tasks.Done < tasks.Total {
		c.add(Warning, "unchecked_tasks", rel("tasks.md"),
			fmt.Sprintf("node is %s but %d of %d tasks are unchecked", n.Status, tasks.Total-tasks.Done, tasks.Total))
	}
	if tasks.Total > 12 {
		c.add(Warning, "atomicity_tasks", rel("tasks.md"),
			fmt.Sprintf("%d task boxes (bound ~12) — the node may not be atomic", tasks.Total))
	}
}

func statusOf(nodeMD string) string {
	for line := range strings.SplitSeq(nodeMD, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "**Status:**"); ok {
			return glyphIn(rest)
		}
	}
	return ""
}

// openItems tallies the Open block of node.md: the text after **Open:** plus the list
// items that follow it, up to the next heading or bold label.
func openItems(nodeMD string) (OpenCounts, bool) {
	var oc OpenCounts
	escalated := false
	tally := func(item string) {
		item = strings.TrimSpace(item)
		switch strings.ToLower(strings.Trim(item, " .")) {
		case "", "none", "—", "-", "…", "n/a":
			return
		}
		if strings.Contains(item, "verify escalation") || strings.Contains(item, "reverted") {
			escalated = true
		}
		switch {
		case strings.Contains(item, "[blocking]"):
			oc.Blocking++
		case strings.Contains(item, "[minor]"):
			oc.Minor++
		case strings.Contains(item, "[deferred:"):
			oc.Deferred++
		default:
			oc.Untagged++
		}
	}
	in := false
	for line := range strings.SplitSeq(nodeMD, "\n") {
		t := strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(t, "**Open:**"); ok {
			in = true
			tally(rest)
			continue
		}
		if !in {
			continue
		}
		switch {
		case strings.HasPrefix(t, "#"), strings.HasPrefix(t, "**"):
			in = false
		case strings.HasPrefix(t, "- "), strings.HasPrefix(t, "* "):
			tally(t[2:])
		}
	}
	return oc, escalated
}

func countTasks(path string) (Count, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Count{}, false
	}
	var n Count
	for line := range strings.SplitSeq(string(b), "\n") {
		if m := checkbox.FindStringSubmatch(line); m != nil {
			n.Total++
			if m[1] != " " {
				n.Done++
			}
		}
	}
	return n, true
}

func gitBranch(projectRoot string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	b := strings.TrimSpace(string(out))
	if b == "HEAD" {
		return ""
	}
	return b
}

// Check runs artifact validation on top of the state findings: well-formed drawio XML,
// spec scenarios with WHEN/THEN, and the spec atomicity bound.
func (c *Change) Check(projectRoot string) {
	dir := filepath.Join(projectRoot, filepath.FromSlash(c.Dir))
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		rel = filepath.ToSlash(rel)
		switch {
		case strings.HasSuffix(d.Name(), ".drawio"):
			if msg := ValidateDrawio(p); msg != "" {
				c.add(Error, "drawio_invalid", rel, msg)
			}
		case d.Name() == "spec.md":
			b, err := os.ReadFile(p)
			if err == nil {
				for _, f := range CheckSpec(string(b)) {
					f.Path = rel
					c.Findings = append(c.Findings, f)
				}
			}
		}
		return nil
	})
}

// ValidateDrawio returns "" when the file is well-formed mxGraph XML, else the problem.
func ValidateDrawio(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return err.Error()
	}
	defer func() { _ = f.Close() }()
	dec := xml.NewDecoder(f)
	root := ""
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "not well-formed XML: " + err.Error()
		}
		if se, ok := tok.(xml.StartElement); ok && root == "" {
			root = se.Name.Local
		}
	}
	switch root {
	case "mxfile", "mxGraphModel":
		return ""
	case "":
		return "empty document — no root element"
	default:
		return fmt.Sprintf("root element is <%s>, want <mxfile>", root)
	}
}

// CheckSpec validates spec.md: every Scenario has a WHEN (or GIVEN) and a THEN step, and the
// node stays within the atomicity bound of ~3 Requirements.
func CheckSpec(text string) []Finding {
	var out []Finding
	requirements := 0
	scenario := ""
	hasWhen, hasThen := false, false
	flush := func() {
		if scenario == "" {
			return
		}
		var missing []string
		if !hasWhen {
			missing = append(missing, "WHEN")
		}
		if !hasThen {
			missing = append(missing, "THEN")
		}
		if len(missing) > 0 {
			out = append(out, Finding{Severity: Error, Code: "scenario_incomplete",
				Message: fmt.Sprintf("Scenario %q has no %s step", scenario, strings.Join(missing, "/"))})
		}
		scenario = ""
	}
	for line := range strings.SplitSeq(text, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "#") {
			flush()
			h := strings.TrimSpace(strings.TrimLeft(t, "#"))
			if strings.HasPrefix(h, "Requirement:") {
				requirements++
			}
			if name, ok := strings.CutPrefix(h, "Scenario:"); ok {
				scenario = strings.TrimSpace(name)
				hasWhen, hasThen = false, false
			}
			continue
		}
		if scenario == "" {
			continue
		}
		step := strings.ToUpper(strings.TrimLeft(t, "-*> "))
		step = strings.TrimLeft(step, "*")
		switch {
		case strings.HasPrefix(step, "WHEN"), strings.HasPrefix(step, "GIVEN"):
			hasWhen = true
		case strings.HasPrefix(step, "THEN"):
			hasThen = true
		}
	}
	flush()
	if requirements > 3 {
		out = append(out, Finding{Severity: Warning, Code: "atomicity_spec",
			Message: fmt.Sprintf("%d Requirements (bound ~3) — the node may not be atomic", requirements)})
	}
	return out
}
