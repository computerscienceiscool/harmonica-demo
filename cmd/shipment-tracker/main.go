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
	step               int
	position, velocity float64
	spring             harmonica.Spring
}

var stops = []string{"Order packed", "Carrier collected", "Out for delivery", "Delivered"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "l", "enter":
			if m.step < len(stops)-1 {
				m.step++
			}
		case "left", "h":
			if m.step > 0 {
				m.step--
			}
		case "r":
			m.step = 0
		}
	case tick:
		m.position, m.velocity = m.spring.Update(m.position, m.velocity, float64(m.step)/3)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  SHIPMENT TRACKER\n  Order #4821 • Portland, OR\n\n"
	for i, stop := range stops {
		dot := "○"
		if float64(i)/3 < m.position+.03 {
			dot = "●"
		}
		s += fmt.Sprintf("  %s %s\n", dot, stop)
	}
	return s + "\n  ←/→ move parcel  •  r restart route  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 5, .6)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
