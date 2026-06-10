// Bubbletea TUI for `prism init`: logo → tool multi-select → live install progress.
// Used only on a TTY without --tools; plain output paths live in init.go.
package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

const logoArt = `██████╗ ██████╗ ██╗███████╗███╗   ███╗
██╔══██╗██╔══██╗██║██╔════╝████╗ ████║
██████╔╝██████╔╝██║███████╗██╔████╔██║
██╔═══╝ ██╔══██╗██║╚════██║██║╚██╔╝██║
██║     ██║  ██║██║███████║██║ ╚═╝ ██║
╚═╝     ╚═╝  ╚═╝╚═╝╚══════╝╚═╝     ╚═╝`

const (
	barWidth = 16
	// stepDelay keeps the install animation perceivable; file writes alone are instant.
	stepDelay = 20 * time.Millisecond
)

var (
	logoRamp = []string{"#B294F9", "#A887F8", "#9D79F7", "#926BF6", "#875DF5", "#7C4FF4"}

	stSubtle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	stSection  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	stCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	stChecked  = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	stOK       = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	stAccent   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	stBarFill  = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	stBarEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

type initPhase int

const (
	phasePick initPhase = iota
	phaseInstall
	phaseDone
)

// installStep is one unit of work; toolIdx == -1 means the shared conventions file.
type installStep struct {
	toolIdx int
	wfIdx   int
}

type stepDoneMsg struct {
	step installStep
	err  error
}

type initModel struct {
	projectRoot string
	detected    map[string]bool
	configured  map[string]bool

	phase    initPhase
	cursor   int
	selected map[int]bool
	canceled bool
	err      error

	spinner    spinner.Model
	queue      []adapters.Tool
	steps      []installStep
	stepIdx    int
	sharedDone bool
	toolDone   []int
}

func newInitModel(projectRoot string) initModel {
	detected := map[string]bool{}
	for _, t := range installer.DetectTools(projectRoot) {
		detected[t.ID] = true
	}
	configured := map[string]bool{}
	for _, t := range installer.ConfiguredTools(projectRoot) {
		configured[t.ID] = true
	}
	selected := map[int]bool{}
	for i, t := range adapters.All {
		selected[i] = detected[t.ID] || configured[t.ID]
	}
	sp := spinner.New(spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))))
	return initModel{
		projectRoot: projectRoot,
		detected:    detected,
		configured:  configured,
		selected:    selected,
		spinner:     sp,
	}
}

func (m initModel) Init() tea.Cmd { return nil }

func (m initModel) pickedTools() []adapters.Tool {
	var picked []adapters.Tool
	for i, t := range adapters.All {
		if m.selected[i] {
			picked = append(picked, t)
		}
	}
	return picked
}

func (m initModel) runStep(i int) tea.Cmd {
	root, queue, steps := m.projectRoot, m.queue, m.steps
	return tea.Tick(stepDelay, func(time.Time) tea.Msg {
		st := steps[i]
		var err error
		if st.toolIdx < 0 {
			_, err = installer.InstallShared(root)
		} else {
			_, err = installer.InstallCommand(root, queue[st.toolIdx], workflows.All[st.wfIdx])
		}
		return stepDoneMsg{step: st, err: err}
	})
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.phase {
		case phasePick:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.canceled = true
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(adapters.All)-1 {
					m.cursor++
				}
			case " ":
				m.selected[m.cursor] = !m.selected[m.cursor]
			case "a":
				all := true
				for i := range adapters.All {
					if !m.selected[i] {
						all = false
						break
					}
				}
				for i := range adapters.All {
					m.selected[i] = !all
				}
			case "enter":
				m.queue = m.pickedTools()
				if len(m.queue) == 0 {
					return m, tea.Quit
				}
				m.steps = []installStep{{toolIdx: -1}}
				for ti := range m.queue {
					for wi := range workflows.All {
						m.steps = append(m.steps, installStep{toolIdx: ti, wfIdx: wi})
					}
				}
				m.toolDone = make([]int, len(m.queue))
				m.phase = phaseInstall
				return m, tea.Batch(m.spinner.Tick, m.runStep(0))
			}
		case phaseInstall:
			if msg.String() == "ctrl+c" {
				m.canceled = true
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stepDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		if msg.step.toolIdx < 0 {
			m.sharedDone = true
		} else {
			m.toolDone[msg.step.toolIdx]++
		}
		m.stepIdx++
		if m.stepIdx >= len(m.steps) {
			m.phase = phaseDone
			return m, tea.Quit
		}
		return m, m.runStep(m.stepIdx)
	}
	return m, nil
}

func logoView() string {
	lines := strings.Split(logoArt, "\n")
	var b strings.Builder
	for i, line := range lines {
		color := logoRamp[i%len(logoRamp)]
		b.WriteString("  " + lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(line) + "\n")
	}
	logoWidth := lipgloss.Width(lines[0])
	sub := fmt.Sprintf("v%s · workflow installer", workflows.Version)
	b.WriteString("  " + stSubtle.Render(lipgloss.PlaceHorizontal(logoWidth, lipgloss.Center, sub)) + "\n")
	return b.String()
}

func sectionTitle(title string, width int) string {
	rule := strings.Repeat("─", max(width-lipgloss.Width(title)-1, 3))
	return "  " + stSection.Render(title) + " " + stSubtle.Render(rule)
}

func progressBar(done, total int) string {
	filled := done * barWidth / total
	return stBarFill.Render(strings.Repeat("█", filled)) +
		stBarEmpty.Render(strings.Repeat("░", barWidth-filled))
}

func (m initModel) maxToolNameWidth() int {
	w := 0
	for _, t := range adapters.All {
		w = max(w, lipgloss.Width(t.Name))
	}
	return w
}

func (m initModel) View() string {
	var b strings.Builder
	b.WriteString("\n" + logoView() + "\n")

	nameW := m.maxToolNameWidth()
	switch m.phase {
	case phasePick:
		b.WriteString(sectionTitle("Select tools", 42) + "\n")
		for i, t := range adapters.All {
			cursor := "  "
			if i == m.cursor {
				cursor = stCursor.Render("▸ ")
			}
			box := stSubtle.Render("◯")
			if m.selected[i] {
				box = stChecked.Render("◉")
			}
			note := ""
			switch {
			case m.configured[t.ID]:
				note = stSubtle.Render("installed")
			case m.detected[t.ID]:
				note = stSubtle.Render("detected")
			}
			name := fmt.Sprintf("%-*s", nameW+2, t.Name)
			if i == m.cursor {
				name = lipgloss.NewStyle().Bold(true).Render(name)
			}
			fmt.Fprintf(&b, "  %s%s %s %s\n", cursor, box, name, note)
		}
		b.WriteString("\n  " + stSubtle.Render("space toggle · a all · enter install · q quit") + "\n")

	case phaseInstall, phaseDone:
		b.WriteString(sectionTitle("Installing", 42) + "\n")
		shared := stSubtle.Render("◦") + " " + stSubtle.Render("conventions.md")
		if m.sharedDone {
			shared = stOK.Render("✔") + " conventions.md " + stSubtle.Render("→ .prism/")
		}
		b.WriteString("  " + shared + "\n")

		total := len(workflows.All)
		curTool := -1
		if m.phase == phaseInstall && m.stepIdx < len(m.steps) {
			curTool = m.steps[m.stepIdx].toolIdx
		}
		for ti, t := range m.queue {
			name := fmt.Sprintf("%-*s", nameW+2, t.Name)
			done := m.toolDone[ti]
			switch {
			case done == total:
				dir := filepath.Dir(t.CommandFile(workflows.All[0].ID)) + "/"
				fmt.Fprintf(&b, "  %s %s %s\n",
					stOK.Render("✔"), name, stSubtle.Render(fmt.Sprintf("%d commands → %s", total, dir)))
			case ti == curTool:
				fmt.Fprintf(&b, "  %s%s %s %s\n",
					m.spinner.View(), name, progressBar(done, total),
					stSubtle.Render(fmt.Sprintf("%d/%d", done, total)))
			default:
				fmt.Fprintf(&b, "  %s %s %s\n",
					stSubtle.Render("◦"), stSubtle.Render(name), stSubtle.Render("waiting"))
			}
		}

		if m.phase == phaseDone {
			first := m.queue[0]
			b.WriteString("\n  " + stOK.Render("✔ Done!") + " Try " +
				stAccent.Render(first.CommandRef("propose")) + " in " + first.Name + " to start a change.\n")
			b.WriteString("  " + stSubtle.Render("Restart your IDE/agent if slash commands do not show up.") + "\n")
		}
	}
	return b.String()
}

func runInitTUI(projectRoot string) error {
	res, err := tea.NewProgram(newInitModel(projectRoot)).Run()
	if err != nil {
		return err
	}
	final := res.(initModel)
	switch {
	case final.err != nil:
		return final.err
	case final.canceled:
		fmt.Println(yellow("Canceled."))
	case final.phase == phasePick:
		fmt.Println(yellow("No tools selected — nothing to do."))
	}
	return nil
}
