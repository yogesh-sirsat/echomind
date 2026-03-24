package ui

import (
	"fmt"
	"io"

	"echomind/internal/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(MainColor).Bold(true)
)

type item struct {
	entry config.HistoryEntry
}

func (i item) Title() string       { return i.entry.FileName }
func (i item) Description() string { return fmt.Sprintf("%s | %s", i.entry.Timestamp.Format("2006-01-02 15:04"), i.entry.FilePath) }
func (i item) FilterValue() string { return i.entry.FileName }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s\n   %s", index+1, i.Title(), i.Description())

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("> "+str))
	} else {
		fmt.Fprint(w, itemStyle.Render(str))
	}
}

type HistoryModel struct {
	list     list.Model
	err      error
	quitting bool
}

func InitialHistoryModel(entries []config.HistoryEntry) HistoryModel {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = item{entry: e}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📊 Recording History"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle

	return HistoryModel{list: l}
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				_ = openFile(i.entry.FilePath)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Padding(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HistoryModel) View() string {
	if m.quitting {
		return ""
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}
