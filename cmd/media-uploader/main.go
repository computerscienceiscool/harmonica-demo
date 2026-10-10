package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"math"
	"os"
	"strings"
	"time"
)

type tick time.Time
type file struct {
	name                       string
	target, progress, velocity float64
}
type model struct {
	selected int
	files    []file
	spring   harmonica.Spring
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
			if m.selected < len(m.files)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "enter":
			m.files[m.selected].target = 1
		case "r":
			for i := range m.files {
				m.files[i].target = 0
				m.files[i].progress = 0
			}
		}
	case tick:
		for i := range m.files {
			m.files[i].progress, m.files[i].velocity = m.spring.Update(m.files[i].progress, m.files[i].velocity, m.files[i].target)
		}
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	s := "\n  MEDIA UPLOADER\n\n"
	for i, x := range m.files {
		p := int(math.Round(x.progress * 20))
		mark := " "
		if i == m.selected {
			mark = "›"
		}
		s += fmt.Sprintf("  %s %-16s [%s%s] %3.0f%%\n", mark, x.name, strings.Repeat("█", p), strings.Repeat("░", 20-p), x.progress*100)
	}
	return s + "\n  ↑/↓ select  •  Enter upload  •  r clear  •  q quit\n"
}
func main() {
	files := []file{{name: "interview.mov"}, {name: "cover-art.png"}, {name: "notes.pdf"}}
	if err := tea.NewProgram(model{files: files, spring: harmonica.NewSpring(harmonica.FPS(60), 5, .65)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
