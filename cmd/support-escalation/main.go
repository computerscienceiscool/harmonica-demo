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

const width = 38

type tick time.Time
type model struct {
	stage              int
	progress, velocity float64
	spring             harmonica.Spring
}

var stages = []string{"Ticket received", "Reproduce", "Investigate", "Resolve", "Confirm"}
var accent = lipgloss.NewStyle().Foreground(lipgloss.Color("#C4B5FD"))

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "l", "right":
			if m.stage < len(stages)-1 {
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
		target := float64(m.stage) / float64(len(stages)-1)
		m.progress, m.velocity = m.spring.Update(m.progress, m.velocity, target)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	filled := int(math.Round(m.progress * width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	view := "\n  SUPPORT ESCALATION\n  Acme account • reporting export unavailable\n\n  " + accent.Render(strings.Repeat("█", filled)) + strings.Repeat("░", width-filled) + fmt.Sprintf("  %3.0f%%\n\n", m.progress*100)
	for i, stage := range stages {
		marker := "○"
		if i < m.stage {
			marker = "✓"
		}
		if i == m.stage {
			marker = "●"
		}
		view += fmt.Sprintf("  %s %s\n", marker, stage)
	}
	return view + "\n  ←/h back  •  →/l advance  •  r reset  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 6, .6)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
