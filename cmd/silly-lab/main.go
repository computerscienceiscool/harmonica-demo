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

const barWidth = 36

type tick time.Time

type scenario struct {
	name, subtitle string
	stages         []string
	accent         lipgloss.Color
}

var scenarios = []scenario{
	{"Incident commander's handoff", "Pager says: everything is exciting", []string{"Page", "Assess", "Mitigate", "Hand off", "Debrief"}, "#FF7A90"},
	{"Feature-flag rollout", "Ship boldly; retain the big red rollback button", []string{"Internal", "Canary", "25%", "75%", "Everyone"}, "#69E6B5"},
	{"Restaurant kitchen ticket", "Table 12: one extra-crispy existential crisis", []string{"Ticket", "Prep", "Sizzle", "Plate", "Runner"}, "#FFB86B"},
	{"Bubble tea assembly", "A precise science with tapioca consequences", []string{"Cup", "Ice", "Tea", "Pearls", "Seal"}, "#EE9FC8"},
	{"Lost sock recovery mission", "Case file #404: Left Sock Not Found", []string{"Hamper", "Dryer", "Couch", "Match", "Reunited"}, "#BFA5FF"},
	{"Space raccoon launch", "Mission control: do not let him press every button", []string{"Snacks", "Helmet", "Fuel", "Count down", "Liftoff"}, "#86D8FF"},
}

type model struct {
	scenario, stage    int
	progress, velocity float64
	spring             harmonica.Spring
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6":
			m.scenario = int(msg.String()[0] - '1')
			m.stage, m.progress, m.velocity = 0, 0, 0
		case "l", "right":
			if m.stage < len(scenarios[m.scenario].stages)-1 {
				m.stage++
			}
		case "h", "left":
			if m.stage > 0 {
				m.stage--
			}
		case "r":
			m.stage = 0
		}
	case tick:
		target := float64(m.stage) / float64(len(scenarios[m.scenario].stages)-1)
		m.progress, m.velocity = m.spring.Update(m.progress, m.velocity, target)
		return m, m.Init()
	}
	return m, nil
}

func (m model) View() string {
	s := scenarios[m.scenario]
	accent := lipgloss.NewStyle().Foreground(s.accent)
	var choices strings.Builder
	for i, option := range scenarios {
		prefix := "  "
		if i == m.scenario {
			prefix = "> "
		}
		fmt.Fprintf(&choices, "%s%d. %s\n", prefix, i+1, option.name)
	}
	filled := int(math.Round(m.progress * barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}
	bar := accent.Render(strings.Repeat("█", filled)) + strings.Repeat("·", barWidth-filled)
	return fmt.Sprintf("\n  HARMONICA SILLY LAB\n\n%s\n  %s\n  %s\n\n  %s  %3.0f%%\n\n%s\n\n  %s\n",
		choices.String(), accent.Render(s.name), s.subtitle, bar, m.progress*100, component(m),
		"1–6 choose  •  ←/h and →/l move  •  r reset  •  q quit")
}

func component(m model) string {
	stage := scenarios[m.scenario].stages[m.stage]
	switch m.scenario {
	case 0:
		return fmt.Sprintf("  ON-CALL  ───────▶  INCIDENT LEAD\n  Status card: %s", stage)
	case 1:
		return fmt.Sprintf("  [internal] [canary] [25%%] [75%%] [everyone]\n  Active flag band: %s", stage)
	case 2:
		return fmt.Sprintf("  ┌ TABLE 12 ──────────────────┐\n  │ %-27s │\n  └─────────────────────────────┘", stage)
	case 3:
		fill := strings.Repeat("▓", m.stage+1) + strings.Repeat("░", 5-m.stage-1)
		return fmt.Sprintf("  ╭─────╮   %s\n  │ %s │   %s", fill, fill, stage)
	case 4:
		return fmt.Sprintf("  [hamper] ── [dryer] ── [couch]\n  Sock radar: %s", stage)
	default:
		return fmt.Sprintf("  FUEL  [%s]  COUNTDOWN  T-%02d\n  RACCOON STATUS: %s", strings.Repeat("▰", m.stage+1), 10-m.stage*2, stage)
	}
}

func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 6, .58)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
