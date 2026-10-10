package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"os"
	"strings"
	"time"
)

type tick time.Time
type model struct {
	snacks, helmet, fuel, launch bool
	thrust, v                    float64
	spring                       harmonica.Spring
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
		case "s":
			m.snacks = !m.snacks
		case "h":
			m.helmet = !m.helmet
		case "f":
			m.fuel = !m.fuel
		case "enter":
			if m.snacks && m.helmet && m.fuel {
				m.launch = true
			}
		}
	case tick:
		target := 0.0
		if m.launch {
			target = 1
		}
		m.thrust, m.v = m.spring.Update(m.thrust, m.v, target)
		return m, m.Init()
	}
	return m, nil
}
func sw(on bool) string {
	if on {
		return "[ ON ]"
	}
	return "[OFF ]"
}
func (m model) View() string {
	bars := int(m.thrust*12 + .5)
	s := "\n  SPACE RACCOON COCKPIT\n  pilot: Professor Trashpanda\n\n" + fmt.Sprintf("  snacks  %s\n  helmet  %s\n  fuel    %s\n\n", sw(m.snacks), sw(m.helmet), sw(m.fuel))
	s += fmt.Sprintf("  THRUST  |%s%s|  %3.0f%%\n", strings.Repeat("█", bars), strings.Repeat("·", 12-bars), m.thrust*100)
	if m.launch {
		s += "  ★ LIFTOFF! hold onto your tiny helmet ★\n"
	}
	return s + "\n  s snacks  •  h helmet  •  f fuel  •  Enter launch  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 5, .55)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
