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
	selected       int
	values, speeds []float64
	springs        []harmonica.Spring
}

var checks = []string{"go test ./...", "go vet ./...", "staticcheck ./...", "build release"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.selected < len(checks)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "enter":
			for i := range m.values {
				m.values[i] = 0
				m.speeds[i] = 0
			}
		case " ":
			m.values[m.selected] = 1
		}
	case tick:
		for i := range m.values {
			m.values[i], m.speeds[i] = m.springs[i].Update(m.values[i], m.speeds[i], m.values[i])
		}
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  TEST DASHBOARD\n\n"
	for i, check := range checks {
		marker := "○"
		if m.values[i] > .98 {
			marker = "✓"
		}
		if i == m.selected {
			marker = "›"
		}
		s += fmt.Sprintf("  %s %-22s %s\n", marker, check, map[bool]string{true: "passed", false: "waiting"}[m.values[i] > .98])
	}
	return s + "\n  ↑/↓ select  •  Space mark passed  •  Enter reset  •  q quit\n"
}
func main() {
	springs := make([]harmonica.Spring, len(checks))
	for i := range springs {
		springs[i] = harmonica.NewSpring(harmonica.FPS(60), 6, .65)
	}
	if err := tea.NewProgram(model{springs: springs, values: make([]float64, len(checks)), speeds: make([]float64, len(checks))}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
