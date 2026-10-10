package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"os"
	"time"
)

type tick time.Time
type item struct{ name, status string }
type model struct {
	selected        int
	items           []item
	pulse, velocity float64
	spring          harmonica.Spring
}

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
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "enter":
			m.items[m.selected].status = "done"
			m.pulse = 1
		case "b":
			m.items[m.selected].status = "blocked"
			m.pulse = 1
		case "r":
			for i := range m.items {
				m.items[i].status = "queued"
			}
		}
	case tick:
		m.pulse, m.velocity = m.spring.Update(m.pulse, m.velocity, 0)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  AGENT WORK QUEUE\n\n"
	for i, x := range m.items {
		mark := " "
		if i == m.selected {
			mark = "›"
		}
		s += fmt.Sprintf("  %s %-27s %-8s\n", mark, x.name, x.status)
	}
	return s + "\n  ↑/↓ choose  •  Enter complete  •  b block  •  r reset  •  q quit\n"
}
func main() {
	items := []item{{"Inspect intake", "queued"}, {"Implement change", "queued"}, {"Run proof", "queued"}, {"Send handoff", "queued"}}
	if err := tea.NewProgram(model{items: items, spring: harmonica.NewSpring(harmonica.FPS(60), 7, .6)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
