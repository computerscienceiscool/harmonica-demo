package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"os"
	"time"
)

type tick time.Time
type model struct {
	selected       int
	flags          []bool
	glow, velocity float64
	spring         harmonica.Spring
}

var cohorts = []string{"staff", "canary", "25%", "75%", "everybody"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			if m.selected > 0 {
				m.selected--
			}
		case "right", "l":
			if m.selected < len(m.flags)-1 {
				m.selected++
			}
		case " ":
			m.flags[m.selected] = !m.flags[m.selected]
			m.glow = 1
		case "a":
			for i := range m.flags {
				m.flags[i] = true
			}
			m.glow = 1
		}
	case tick:
		m.glow, m.velocity = m.spring.Update(m.glow, m.velocity, 0)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  FEATURE FLAG CONTROL ROOM\n  flag: extra-crispy-search • rollback lever armed\n\n"
	for i, c := range cohorts {
		cursor := " "
		if i == m.selected {
			cursor = "›"
		}
		state := "OFF"
		lamp := "○"
		if m.flags[i] {
			state = " ON"
			lamp = "●"
		}
		s += fmt.Sprintf("  %s [%s]  %-10s  %s\n", cursor, lamp, c, state)
	}
	return s + "\n  ←/→ choose cohort  •  Space toggle flag  •  a enable all  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{flags: make([]bool, 5), spring: harmonica.NewSpring(harmonica.FPS(60), 8, .6)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
