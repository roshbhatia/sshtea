package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/roshbhatia/sshtea/config"
	"github.com/roshbhatia/sshtea/host"
)

type model struct {
	hosts     list.Model
	state     string
	textInput textinput.Model
	err       error
}

const (
	stateList          = "list"
	stateAdd           = "add"
	stateEdit          = "edit"
	stateConfirmDelete = "confirm_delete"
	stateHelp          = "help"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func InitialModel() model {
	hosts := config.LoadHosts()

	l := list.New(hosts, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Hosts"

	ti := textinput.New()
	ti.Placeholder = "Enter host details"
	ti.Focus()

	return model{
		hosts:     l,
		state:     stateList,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateList:
			return m.handleListState(msg)
		case stateAdd:
			return m.handleAddState(msg)
		case stateEdit:
			return m.handleEditState(msg)
		case stateConfirmDelete:
			return m.handleConfirmDeleteState(msg)
		case stateHelp:
			return m.handleHelpState(msg)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.hosts.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.hosts, cmd = m.hosts.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateList:
		return docStyle.Render(m.hosts.View() + "\n\nPress (a) to add, (e) to edit, (d) to delete, (h) for help")
	case stateAdd:
		return fmt.Sprintf(
			"Add new host (format: name hostname user port)\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to cancel)",
		) + "\n"
	case stateEdit:
		return fmt.Sprintf(
			"Edit host (format: name hostname user port)\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to cancel)",
		) + "\n"
	case stateConfirmDelete:
		return "Are you sure you want to delete this host? (y/n)"
	case stateHelp:
		return `
SSH Host Manager Help:

a: Add a new host
e: Edit the selected host
d: Delete the selected host
h: Show this help menu
q: Quit the application

Press any key to return to the list view.
`
	default:
		return "Invalid state"
	}
}

func (m model) handleListState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "a":
		m.state = stateAdd
		m.textInput.SetValue("")
		m.textInput.Focus()
	case "e":
		if i, ok := m.hosts.SelectedItem().(host.Host); ok {
			m.state = stateEdit
			m.textInput.SetValue(fmt.Sprintf("%s %s %s %s", i.Name, i.Hostname, i.User, i.Port))
			m.textInput.Focus()
		}
	case "d":
		m.state = stateConfirmDelete
	case "h":
		m.state = stateHelp
	}
	return m, nil
}

func (m model) handleAddState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		parts := strings.Split(m.textInput.Value(), " ")
		if len(parts) == 4 {
			newHost := host.Host{Name: parts[0], Hostname: parts[1], User: parts[2], Port: parts[3]}
			m.hosts.InsertItem(len(m.hosts.Items())-1, newHost)
			err := config.SaveHosts(m.hosts.Items())
			if err != nil {
				m.err = err
			}
			m.state = stateList
		}
	case "esc":
		m.state = stateList
	}
	return m, nil
}

func (m model) handleEditState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		parts := strings.Split(m.textInput.Value(), " ")
		if len(parts) == 4 {
			editedHost := host.Host{Name: parts[0], Hostname: parts[1], User: parts[2], Port: parts[3]}
			index := m.hosts.Index()
			m.hosts.RemoveItem(index)
			m.hosts.InsertItem(index, editedHost)
			err := config.SaveHosts(m.hosts.Items())
			if err != nil {
				m.err = err
			}
			m.state = stateList
		}
	case "esc":
		m.state = stateList
	}
	return m, nil
}

func (m model) handleConfirmDeleteState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if i, ok := m.hosts.SelectedItem().(host.Host); ok {
			newHosts := []list.Item{}
			for _, h := range m.hosts.Items() {
				if h.(host.Host).Name != i.Name {
					newHosts = append(newHosts, h)
				}
			}
			m.hosts.SetItems(newHosts)
			err := config.SaveHosts(newHosts)
			if err != nil {
				m.err = err
			}
		}
		m.state = stateList
	case "n", "esc":
		m.state = stateList
	}
	return m, nil
}

func (m model) handleHelpState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "h":
		m.state = stateList
	}
	return m, nil
}
