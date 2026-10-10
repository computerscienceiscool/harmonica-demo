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
	lane   int
	x, v   float64
	spring harmonica.Spring
}

var lanes = []string{"TICKET", "PREP", "SIZZLE", "PLATE", "RUNNER"}

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
			if m.lane < len(lanes)-1 {
				m.lane++
			}
		case "h", "left":
			if m.lane > 0 {
				m.lane--
			}
		}
	case tick:
		m.x, m.v = m.spring.Update(m.x, m.v, float64(m.lane)/4)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  KITCHEN LINE • TABLE 12\n  one mushroom burger, fries, no judgement\n\n"
	for i, l := range lanes {
		ticket := "                 "
		if i == m.lane {
			ticket = " [ #12 BURGER ] "
		}
		s += fmt.Sprintf("  %-8s │%s│\n", l, ticket)
	}
	return s + "\n  ←/→ move ticket across stations  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 6, .58)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
