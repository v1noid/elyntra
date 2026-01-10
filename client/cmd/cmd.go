package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	keys string
	ti   textinput.Model
}

func initModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Hie there"
	ti.Width = 20
	return model{
		ti: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if val, ok := msg.(tea.KeyMsg); ok {
		key := val.String()
		m.keys += key
		os.WriteFile("key.log", []byte(m.keys), 0644)
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "enter" {
			os.WriteFile("output.log", []byte(m.ti.Value()), 0644)
			return m, tea.Quit
		}

	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("Hi ther, what's is your name \n %v", m.ti.View())
}

func main() {
	t := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
