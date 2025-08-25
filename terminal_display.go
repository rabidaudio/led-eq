package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NimbleMarkets/ntcharts/sparkline"
	tea "github.com/charmbracelet/bubbletea"
)

type TerminalDisplay struct {
	eq      *EQ
	compute func(res []float64)

	sl sparkline.Model
}

type render struct{ data []float64 }

func (td *TerminalDisplay) every() tea.Cmd {
	return tea.Every(td.eq.Timestep, func(t time.Time) tea.Msg {
		v := make([]float64, td.eq.NumBins)
		td.compute(v)
		for i := range v {
			v[i] = v[i] * 10
		}
		return render{data: v}
	})
}

func (td *TerminalDisplay) Init() tea.Cmd {
	return td.every()
}

func (td *TerminalDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return td, tea.Quit
		}
	case render:
		td.sl.PushAll(msg.data)
		td.sl.Draw()
		return td, td.every()
	}
	return td, nil
}

func (td *TerminalDisplay) View() string {
	return td.sl.View()
}

var _ tea.Model = (*TerminalDisplay)(nil)

func NewTerminalDisplay(eq *EQ, compute func(res []float64)) *TerminalDisplay {
	sl := sparkline.New(eq.NumBins, 10)
	td := TerminalDisplay{eq: eq, sl: sl, compute: compute}
	return &td
}

// func (td *TerminalDisplay) Render(values []float64) error {
// 	v := make([]float64, len(values))
// 	for i := range values {
// 		v[i] = values[i] * 10
// 	}
// 	td.c <- render{data: v}
// 	return nil
// }

func (td *TerminalDisplay) Run() {
	if _, err := tea.NewProgram(td).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
