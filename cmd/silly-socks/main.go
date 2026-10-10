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
	cursor, found int
	pulse, v      float64
	spring        harmonica.Spring
}

var rooms = []string{"hamper", "dryer", "couch", "closet", "under bed", "laundry pile", "shoe rack", "cat fort", "victory"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			if m.cursor%3 > 0 {
				m.cursor--
			}
		case "l", "right":
			if m.cursor%3 < 2 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor >= 3 {
				m.cursor -= 3
			}
		case "j", "down":
			if m.cursor < 6 {
				m.cursor += 3
			}
		case "enter":
			if m.cursor == 8 {
				m.found = 1
			}
			m.pulse = 1
		}
	case tick:
		m.pulse, m.v = m.spring.Update(m.pulse, m.v, 0)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  LOST SOCK RECOVERY MISSION\n  matching confidence: "
	if m.found == 1 {
		s += "100%  ♥ PAIR REUNITED ♥\n\n"
	} else {
		s += "12%  scanning...\n\n"
	}
	for i, r := range rooms {
		mark := " "
		if i == m.cursor {
			mark = "›"
		}
		s += fmt.Sprintf("  %s %-14s", mark, r)
		if i%3 == 2 {
			s += "\n"
		}
	}
	return s + "\n  h/j/k/l explore map  •  Enter search tile  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 9, .5)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
