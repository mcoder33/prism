package workflows

import (
	"regexp"
	"strings"
	"testing"
)

// Methodology lint: machine-checkable invariants over the embedded templates.
// Guards the bug class that manual reviews kept finding — dangling cross-references,
// status-glyph drift, invalid command placeholders, and conventions bloat.

var statusGlyphs = []string{"⚪", "🟡", "🟢", "🔵", "✅", "⏸"}

func allBodies() map[string]string {
	m := map[string]string{"conventions": Conventions()}
	for _, w := range All {
		m[w.ID] = Body(w.ID)
	}
	return m
}

func conventionHeadings() []string {
	var hs []string
	for line := range strings.SplitSeq(Conventions(), "\n") {
		for _, prefix := range []string{"### ", "## "} {
			if h, ok := strings.CutPrefix(line, prefix); ok {
				hs = append(hs, h)
				break
			}
		}
	}
	return hs
}

// References like "see conventions, Open tags)" or "(criteria — conventions, Change tiers)"
// must name a section that actually exists in conventions.md.
var sectionRefPattern = regexp.MustCompile(`conventions, ([A-Za-z][A-Za-z .&-]*?)\)`)

func TestConventionSectionRefsResolve(t *testing.T) {
	headings := conventionHeadings()
	if len(headings) == 0 {
		t.Fatal("no headings found in conventions")
	}
	for id, body := range allBodies() {
		for _, m := range sectionRefPattern.FindAllStringSubmatch(body, -1) {
			ref := strings.TrimSpace(m[1])
			found := false
			for _, h := range headings {
				if strings.HasPrefix(h, ref) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s references conventions section %q which does not exist", id, ref)
			}
		}
	}
}

// The status legend, the transition table, and the node.md template must agree on the
// full glyph set — a glyph added to one place and forgotten in another breaks resume.
func TestStatusGlyphsConsistent(t *testing.T) {
	c := Conventions()

	start := strings.Index(c, "## Statuses and transitions")
	end := strings.Index(c, "### README.md")
	if start < 0 || end < 0 || end <= start {
		t.Fatal("Statuses and transitions section not found in conventions")
	}
	section := c[start:end]
	for _, g := range statusGlyphs {
		if !strings.Contains(section, g) {
			t.Errorf("glyph %q missing from the Statuses and transitions section", g)
		}
	}

	_, rest, ok := strings.Cut(c, "### node.md")
	if !ok {
		t.Fatal("node.md template not found in conventions")
	}
	// limit to the node.md template block itself (up to the next heading)
	block, _, _ := strings.Cut(rest, "\n### ")
	for _, g := range statusGlyphs {
		if !strings.Contains(block, g) {
			t.Errorf("node.md template Status line missing glyph %q", g)
		}
	}
}

// Every {{cmd:*}} placeholder must name a real workflow id.
var cmdRefPattern = regexp.MustCompile(`\{\{cmd:([a-z-]+)\}\}`)

func TestCmdPlaceholdersValid(t *testing.T) {
	ids := map[string]bool{}
	for _, w := range All {
		ids[w.ID] = true
	}
	for id, body := range allBodies() {
		for _, m := range cmdRefPattern.FindAllStringSubmatch(body, -1) {
			if !ids[m[1]] {
				t.Errorf("%s uses {{cmd:%s}} which is not a workflow id", id, m[1])
			}
		}
	}
}

// conventions.md is read by EVERY command — it is the per-invocation context tax.
// Trim before adding; raising this budget is a deliberate decision, not a default.
func TestConventionsSizeBudget(t *testing.T) {
	if n := len(Conventions()); n > 14*1024 {
		t.Errorf("conventions.md is %d bytes; budget is 14KB — trim before adding more", n)
	}
	for _, w := range All {
		if n := len(Body(w.ID)); n > 20*1024 {
			t.Errorf("%s body is %d bytes; budget is 20KB", w.ID, n)
		}
	}
}

// Every command leans on the shared conventions — each must point the agent at them.
func TestCommandsPointToConventions(t *testing.T) {
	for _, w := range All {
		if !strings.Contains(Body(w.ID), "conventions") {
			t.Errorf("%s never mentions conventions.md", w.ID)
		}
	}
}

// Every glyph→glyph flip mentioned anywhere in the templates must be a legal transition
// from the conventions table — the table is the state machine; commands may not invent flips.
var transitionPairPattern = regexp.MustCompile(`([⚪🟡🟢🔵✅⏸])\s*→\s*([⚪🟡🟢🔵✅⏸])`)

func TestStatusTransitionsLegal(t *testing.T) {
	c := Conventions()
	tableStart := strings.Index(c, "| Transition | Owner |")
	if tableStart < 0 {
		t.Fatal("transition table not found in conventions")
	}
	table, _, _ := strings.Cut(c[tableStart:], "\n\n")

	legal := map[string]bool{}
	for _, m := range transitionPairPattern.FindAllStringSubmatch(table, -1) {
		legal[m[1]+m[2]] = true
	}
	if strings.Contains(table, "any → ⏸") {
		for _, g := range statusGlyphs {
			legal[g+"⏸"] = true
		}
	}
	if len(legal) < 5 {
		t.Fatalf("parsed only %d transitions from the table — parser or table broken", len(legal))
	}

	for id, body := range allBodies() {
		for _, m := range transitionPairPattern.FindAllStringSubmatch(body, -1) {
			if !legal[m[1]+m[2]] {
				t.Errorf("%s mentions flip %s → %s which is not in the conventions transition table", id, m[1], m[2])
			}
		}
	}
}
