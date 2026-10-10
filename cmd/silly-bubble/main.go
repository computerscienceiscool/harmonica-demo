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
	selected    int
	ingredients []bool
	fill, v     float64
	spring      harmonica.Spring
}

var bits = []string{"ice", "jasmine tea", "boba pearls", "brown sugar", "seal"}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return tick(t) })
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.selected < len(bits)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case " ":
			m.ingredients[m.selected] = !m.ingredients[m.selected]
		}
	case tick:
		count := 0
		for _, on := range m.ingredients {
			if on {
				count++
			}
		}
		m.fill, m.v = m.spring.Update(m.fill, m.v, float64(count)/5)
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	rows := int(m.fill*5 + .5)
	cup := ""
	for i := 0; i < 5; i++ {
		line := "     "
		if i >= 5-rows {
			line = "~~~~~"
		}
		cup += fmt.Sprintf("  │%s│\n", line)
	}
	s := "\n  BUBBLE TEA BUILDER\n\n"
	for i, b := range bits {
		cursor := " "
		if i == m.selected {
			cursor = "›"
		}
		check := "○"
		if m.ingredients[i] {
			check = "●"
		}
		s += fmt.Sprintf("  %s %s %s\n", cursor, check, b)
	}
	return s + "\n" + cup + "  j/k choose  •  Space add/remove ingredient  •  q quit\n"
}
func main() {
	if err := tea.NewProgram(model{ingredients: make([]bool, 5), spring: harmonica.NewSpring(harmonica.FPS(60), 7, .55)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
