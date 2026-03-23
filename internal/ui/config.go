package ui

import (
	"fmt"
	"strings"

	"echomind/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfigModel struct {
	config      config.Config
	formatInput textinput.Model
	dirInput    textinput.Model
	qualities   []string
	qualityIdx  int
	focusIndex  int
	submitted   bool
	err         error
}

func InitialConfigModel(cfg config.Config) ConfigModel {
	fi := textinput.New()
	fi.Placeholder = "wav, mp3, flac"
	fi.SetValue(cfg.DefaultFormat)
	fi.Focus()

	di := textinput.New()
	di.Placeholder = "C:\\Users\\..."
	di.SetValue(cfg.DefaultDirectory)

	qualities := []string{"low", "medium", "high"}
	qualityIdx := 1 // medium
	for i, q := range qualities {
		if q == cfg.DefaultQuality {
			qualityIdx = i
			break
		}
	}

	return ConfigModel{
		config:      cfg,
		formatInput: fi,
		dirInput:    di,
		qualities:   qualities,
		qualityIdx:  qualityIdx,
		focusIndex:  0,
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			return m.updateFocus()

		case "left", "right":
			if m.focusIndex == 1 {
				if msg.String() == "left" {
					m.qualityIdx = (m.qualityIdx - 1 + len(m.qualities)) % len(m.qualities)
				} else {
					m.qualityIdx = (m.qualityIdx + 1) % len(m.qualities)
				}
				return m, nil
			}

		case "enter":
			if m.focusIndex == 2 {
				m.config.DefaultFormat = strings.ToLower(m.formatInput.Value())
				m.config.DefaultQuality = m.qualities[m.qualityIdx]
				m.config.DefaultDirectory = m.dirInput.Value()
				err := config.Save(m.config)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.submitted = true
				return m, tea.Quit
			}
			m.focusIndex++
			return m.updateFocus()
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.formatInput, cmd = m.formatInput.Update(msg)
	} else if m.focusIndex == 2 {
		m.dirInput, cmd = m.dirInput.Update(msg)
	}
	return m, cmd
}

func (m ConfigModel) updateFocus() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		cmd = m.formatInput.Focus()
		m.dirInput.Blur()
	case 1:
		m.formatInput.Blur()
		m.dirInput.Blur()
	case 2:
		m.formatInput.Blur()
		cmd = m.dirInput.Focus()
	}
	return m, cmd
}

func (m ConfigModel) View() string {
	if m.submitted {
		return lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Padding(1).
			Render("✨ Configuration saved successfully!")
	}

	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var s strings.Builder
	s.WriteString(TitleStyle.Render("⚙️  EchoMind Settings"))
	s.WriteString("\n\n")

	s.WriteString(m.renderField("Default Format:", m.formatInput.View(), m.focusIndex == 0))
	s.WriteString("\n")

	var qualView strings.Builder
	for i, q := range m.qualities {
		style := lipgloss.NewStyle().Padding(0, 1)
		if i == m.qualityIdx {
			style = style.Foreground(lipgloss.Color("#000000")).Background(MainColor).Bold(true)
		} else {
			style = style.Foreground(MutedColor)
		}
		qualView.WriteString(style.Render(strings.ToUpper(q)) + " ")
	}
	s.WriteString(m.renderField("Audio Quality (Left/Right):", qualView.String(), m.focusIndex == 1))
	s.WriteString("\n")

	s.WriteString(m.renderField("Default Directory:", m.dirInput.View(), m.focusIndex == 2))
	s.WriteString("\n")

	s.WriteString(StatusStyle.Render("(arrows/tab to navigate, enter to save, q to quit)"))

	return s.String()
}

func (m ConfigModel) renderField(label string, value string, focused bool) string {
	labelStyle := PromptStyle
	if focused {
		labelStyle = labelStyle.Copy().Foreground(lipgloss.Color("#FFF")).Background(MainColor).Padding(0, 1)
	}
	return fmt.Sprintf("%s\n%s\n", labelStyle.Render(label), value)
}
