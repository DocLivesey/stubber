package main

import (
	"fmt"
	"os"

	// "github.com/charmbracelet/bubbles/list"

	"github.com/DocLivesey/stubber/data"
	"github.com/DocLivesey/stubber/viewstub"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list viewstub.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	stubs := data.Populate()
	var items []viewstub.Item
	for _, s := range stubs {
		items = append(items, s)
	}
	// items := []viewstub.Item{
	// 	Stub{jar: "jar1.jar", path: "/home/stubfolder1", state: true, pid: "10002"},
	// 	Stub{jar: "jar2.jar", path: "/home/stubfolder2", state: false, pid: "-"},
	// }

	m := model{list: viewstub.New(items, viewstub.NewDefaultDelegate(), 0, 0)}
	// m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
