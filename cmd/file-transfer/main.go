package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const width = 38

var accent = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DD3FC"))

type tick time.Time
type model struct {
	progress, velocity, target float64
	spring                     harmonica.Spring
	paused                     bool
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
		case " ":
			m.paused = !m.paused
		case "r":
			m.target = 0
			m.paused = false
		case "enter":
			m.target = 1
			m.paused = false
		}
	case tick:
		if !m.paused {
			m.progress, m.velocity = m.spring.Update(m.progress, m.velocity, m.target)
		}
		return m, m.Init()
	}
	return m, nil
}
func (m model) View() string {
	filled := int(math.Round(m.progress * width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := accent.Render(strings.Repeat("█", filled)) + strings.Repeat("░", width-filled)
	state := "ready"
	if m.target == 1 {
		state = "copying"
	}
	if m.paused {
		state = "paused"
	}
	if m.progress > .995 {
		state = "checksum verified"
	}
	return fmt.Sprintf("\n  FILE TRANSFER MONITOR\n  nightly-backup.tar.zst → archive/cas\n\n  %s  %3.0f%%\n  %s\n\n  Enter start  •  Space pause/resume  •  r reset  •  q quit\n", bar, m.progress*100, state)
}
func main() {
	if err := tea.NewProgram(model{spring: harmonica.NewSpring(harmonica.FPS(60), 5, .7)}).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
