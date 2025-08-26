package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NimbleMarkets/ntcharts/sparkline"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rabidaudio/led-eq/eq"
)

// arbitrary, just make the scale bigger cheaply
const scaleFactor = 8

type TerminalDisplay struct {
	eq  *eq.EQ
	msg chan tea.Msg
	sl  sparkline.Model
	avg float64
}

type render struct{ data []float64 }
type done struct{}

var _ tea.Model = done{}

func (done) Init() tea.Cmd {
	return nil
}

func (done) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return done{}, nil
}

func (done) View() string {
	return "\n"
}

func (td *TerminalDisplay) awaitNext() tea.Cmd {
	return tea.Every(1*time.Second/60.0, func(t time.Time) tea.Msg {
		return <-td.msg
	})
}

func (td *TerminalDisplay) Init() tea.Cmd {
	return td.awaitNext()
}

func (td *TerminalDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return done{}, tea.Quit
		}
	case done:
		return done{}, tea.Quit
	case render:
		td.avg = 0
		for _, v := range msg.data {
			td.avg += v
		}
		td.avg /= float64(len(msg.data))
		td.sl.PushAll(msg.data)
		td.sl.Draw()
		return td, td.awaitNext()
	}
	return td, nil
}

func (td *TerminalDisplay) View() string {
	return fmt.Sprintf("avg: %f\n", td.avg) + td.sl.View()
}

var _ Display = (*TerminalDisplay)(nil)
var _ tea.Model = (*TerminalDisplay)(nil)

func NewTerminalDisplay(eq *eq.EQ) *TerminalDisplay {
	sl := sparkline.New(eq.OutBins.Len(), scaleFactor)
	td := TerminalDisplay{eq: eq, sl: sl}
	td.msg = make(chan tea.Msg)
	return &td
}

func (td *TerminalDisplay) Render(values []float64) error {
	v := make([]float64, len(values))
	for i := range values {
		v[i] = values[i] * scaleFactor
	}
	td.msg <- render{data: v}
	return nil
}

func (td *TerminalDisplay) Done() {
	td.msg <- done{}
}

func (td *TerminalDisplay) Run() {
	if _, err := tea.NewProgram(td, tea.WithFPS(90)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
