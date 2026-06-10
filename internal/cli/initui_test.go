package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

func key(t tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: t} }

func space() tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}} }

func update(t *testing.T, m initModel, msg tea.Msg) (initModel, tea.Cmd) {
	t.Helper()
	model, cmd := m.Update(msg)
	next, ok := model.(initModel)
	if !ok {
		t.Fatalf("Update returned %T, want initModel", model)
	}
	return next, cmd
}

func TestNewInitModelPreselectsDetectedTools(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	m := newInitModel(root)
	for i, tl := range adapters.All {
		want := tl.ID == "claude"
		if m.selected[i] != want {
			t.Errorf("selected[%s] = %v, want %v", tl.ID, m.selected[i], want)
		}
	}
}

func TestPickNavigationToggleAndCancel(t *testing.T) {
	m := newInitModel(t.TempDir())

	m, _ = update(t, m, key(tea.KeyDown))
	m, _ = update(t, m, key(tea.KeyDown))
	if m.cursor != 2 {
		t.Fatalf("cursor = %d, want 2", m.cursor)
	}
	m, _ = update(t, m, key(tea.KeyUp))
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", m.cursor)
	}

	m, _ = update(t, m, space())
	if !m.selected[1] {
		t.Fatal("space should toggle selection on")
	}
	m, _ = update(t, m, space())
	if m.selected[1] {
		t.Fatal("space should toggle selection off")
	}

	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	for i := range adapters.All {
		if !m.selected[i] {
			t.Fatal("'a' should select all")
		}
	}
	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	for i := range adapters.All {
		if m.selected[i] {
			t.Fatal("'a' twice should deselect all")
		}
	}

	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if !m.canceled {
		t.Fatal("q should cancel")
	}
}

func TestEnterBuildsInstallQueue(t *testing.T) {
	m := newInitModel(t.TempDir())
	m.selected[0] = true // Claude

	m, cmd := update(t, m, key(tea.KeyEnter))
	if m.phase != phaseInstall {
		t.Fatalf("phase = %d, want phaseInstall", m.phase)
	}
	if cmd == nil {
		t.Fatal("enter must kick off the first install step")
	}
	wantSteps := 1 + len(workflows.All) // shared + one tool
	if len(m.steps) != wantSteps {
		t.Fatalf("steps = %d, want %d", len(m.steps), wantSteps)
	}
	if m.steps[0].toolIdx != -1 {
		t.Fatal("first step must be the shared conventions")
	}
}

func TestEnterWithNothingSelectedQuits(t *testing.T) {
	m := newInitModel(t.TempDir())
	m, _ = update(t, m, key(tea.KeyEnter))
	if m.phase != phasePick {
		t.Fatal("phase must stay phasePick when nothing is selected")
	}
}

func TestStepDoneAdvancesToDone(t *testing.T) {
	m := newInitModel(t.TempDir())
	m.selected[0] = true
	m, _ = update(t, m, key(tea.KeyEnter))

	for m.phase == phaseInstall {
		m, _ = update(t, m, stepDoneMsg{step: m.steps[m.stepIdx]})
	}

	if m.phase != phaseDone {
		t.Fatalf("phase = %d, want phaseDone", m.phase)
	}
	if !m.sharedDone {
		t.Fatal("sharedDone must be set")
	}
	if m.toolDone[0] != len(workflows.All) {
		t.Fatalf("toolDone = %d, want %d", m.toolDone[0], len(workflows.All))
	}
}

func TestStepDoneErrorAborts(t *testing.T) {
	m := newInitModel(t.TempDir())
	m.selected[0] = true
	m, _ = update(t, m, key(tea.KeyEnter))

	m, _ = update(t, m, stepDoneMsg{step: m.steps[0], err: os.ErrPermission})
	if m.err == nil {
		t.Fatal("step error must be recorded")
	}
}

func TestViewByPhase(t *testing.T) {
	m := newInitModel(t.TempDir())
	pick := m.View()
	for _, want := range []string{"Select tools", "Claude Code", "space toggle"} {
		if !strings.Contains(pick, want) {
			t.Errorf("pick view missing %q", want)
		}
	}

	m.selected[0] = true
	m, _ = update(t, m, key(tea.KeyEnter))
	install := m.View()
	for _, want := range []string{"Installing", "conventions.md"} {
		if !strings.Contains(install, want) {
			t.Errorf("install view missing %q", want)
		}
	}

	for m.phase == phaseInstall {
		m, _ = update(t, m, stepDoneMsg{step: m.steps[m.stepIdx]})
	}
	done := m.View()
	for _, want := range []string{"Done!", "propose"} {
		if !strings.Contains(done, want) {
			t.Errorf("done view missing %q", want)
		}
	}
}

func TestProgressBarFill(t *testing.T) {
	half := progressBar(4, 8)
	if got := strings.Count(half, "█"); got != barWidth/2 {
		t.Errorf("filled = %d, want %d", got, barWidth/2)
	}
	if got := strings.Count(half, "░"); got != barWidth/2 {
		t.Errorf("empty = %d, want %d", got, barWidth/2)
	}
	full := progressBar(8, 8)
	if got := strings.Count(full, "█"); got != barWidth {
		t.Errorf("full = %d, want %d", got, barWidth)
	}
}
