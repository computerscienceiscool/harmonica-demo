package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps       = 60
	barWidth  = 48
	frequency = 6.0
	damping   = 0.55
)

var (
	filledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#575BD8"))
	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#575BD8"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "246", Dark: "241"})
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

type model struct {
	progress float64
	velocity float64
	target   float64
	spring   harmonica.Spring
	scenario int
	stage    int
}

type scenario struct {
	name        string
	description string
	stages      []string
	actions     string
}

var scenarios = []scenario{
	{
		name:        "File transfer",
		description: "Copy a backup or Grid/CAS bundle.",
		stages:      []string{"Choose files", "Copy", "Verify checksum", "Complete"},
		actions:     "start • pause • cancel • retry",
	},
	{
		name:        "Test suite",
		description: "Run selected project checks.",
		stages:      []string{"Select checks", "Unit tests", "Integration tests", "Report"},
		actions:     "run selected • show failures • rerun failed",
	},
	{
		name:        "Order shipment",
		description: "Prepare and complete one shipment safely.",
		stages:      []string{"Invoice found", "Review", "Label purchased", "Printed"},
		actions:     "approve • reject • retry safely",
	},
	{
		name:        "Agent work queue",
		description: "Move one work item through its delivery stages.",
		stages:      []string{"Research", "Implement", "Test", "Hand off"},
		actions:     "change task • mark blocked • request review",
	},
	{
		name:        "Media upload",
		description: "Prepare, upload, and verify one photo or video.",
		stages:      []string{"Choose file", "Prepare", "Upload", "Verify"},
		actions:     "choose file • pause/resume • retry verification",
	},
}

func (m model) Init() tea.Cmd {
	return animate()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5":
			m.scenario = int(msg.String()[0] - '1')
			m.stage = 0
			m.target = 0
		case "right", "l", "+":
			if m.stage < len(scenarios[m.scenario].stages)-1 {
				m.stage++
			}
			m.target = stageProgress(m.scenario, m.stage)
		case "left", "h", "-":
			if m.stage > 0 {
				m.stage--
			}
			m.target = stageProgress(m.scenario, m.stage)
		case "r":
			m.stage = 0
			m.target = 0
		}
	case frameMsg:
		m.progress, m.velocity = m.spring.Update(m.progress, m.velocity, m.target)
		return m, animate()
	}

	return m, nil
}

func (m model) View() string {
	selected := scenarios[m.scenario]
	filled := int(math.Round(m.progress * barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}

	bar := filledStyle.Render(strings.Repeat(" ", filled)) +
		emptyStyle.Render(strings.Repeat("·", barWidth-filled))

	var choices strings.Builder
	for i, option := range scenarios {
		prefix := "  "
		if i == m.scenario {
			prefix = "> "
		}
		fmt.Fprintf(&choices, "%s%d. %s\n", prefix, i+1, option.name)
	}

	var stages strings.Builder
	for i, stage := range selected.stages {
		marker := "○"
		if i < m.stage {
			marker = "✓"
		}
		if i == m.stage {
			marker = "●"
		}
		fmt.Fprintf(&stages, "  %s %s\n", marker, stage)
	}

	return fmt.Sprintf(
		"\n  Harmonica animated progress bar\n\n%s\n  %s\n  %s  %3.0f%%\n\n%s\n  Actions: %s\n\n  %s\n",
		choices.String(),
		selected.name,
		bar,
		m.progress*100,
		stages.String(),
		selected.actions,
		helpStyle.Render("1–5 choose a scenario • ←/h and →/l move stages • r resets • q quits"),
	)
}

func stageProgress(scenarioIndex, stage int) float64 {
	stageCount := len(scenarios[scenarioIndex].stages)
	if stageCount < 2 {
		return 0
	}

	return float64(stage) / float64(stageCount-1)
}

func main() {
	m := model{
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running animated progress bar:", err)
		os.Exit(1)
	}
}
