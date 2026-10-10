package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

type tick time.Time
type model struct {
	focus, owner int
	x, v         float64
	spring       harmonica.Spring
}

var roles = []string{"ON-CALL", "INCIDENT LEAD", "COMMS"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right":
			m.focus = (m.focus + 1) % len(roles)
		case "left":
			m.focus = (m.focus + len(roles) - 1) % len(roles)
		case "enter":
			m.owner = m.focus
		}
	case tick:
		m.x, m.v = m.spring.Update(m.x, m.v, float64(m.owner)/2)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  INCIDENT HANDOFF BOARD\n  SRE-418 • checkout latency • severity: loud\n\n"
	for i, r := range roles {
		mark := " "
		if i == m.focus {
			mark = "›"
		}
		baton := ""
		if i == m.owner {
			baton = "  [ BATON ]"
		}
		s += fmt.Sprintf("  %s ┌─────────────────┐\n    │ %-15s │%s\n    └─────────────────┘\n", mark, r, baton)
	}
	return s + "\n  Tab choose recipient  •  Enter hand off baton  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 7, .55)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
