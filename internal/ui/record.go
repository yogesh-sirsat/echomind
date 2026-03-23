package ui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"echomind/internal/audio"
	"echomind/internal/config"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SessionState int

const (
	StateStarting SessionState = iota
	StateRecording
	StateSaving
	StateFinished
)

type RecordModel struct {
	state       SessionState
	config      config.Config
	progress    progress.Model
	waveform    []float64
	startTime   time.Time
	fileName    string
	filePath    string
	tickCount   int
	blink       bool
	recorder    *audio.Recorder
	err         error
	
	fileNameInput textinput.Model
	formats       []string
	formatIndex   int
	dirInput      textinput.Model
	focusIndex    int
}

func InitialRecordModel(cfg config.Config) RecordModel {
	p := progress.New(progress.WithDefaultGradient())
	
	sampleRate := uint32(44100)
	switch cfg.DefaultQuality {
	case "low":
		sampleRate = 22050
	case "medium":
		sampleRate = 44100
	case "high":
		sampleRate = 48000
	}
	rec, _ := audio.NewRecorder(sampleRate)

	fni := textinput.New()
	fni.Placeholder = "Enter file name..."
	fni.SetValue(fmt.Sprintf("recording_%s", time.Now().Format("2006-01-02_15-04-05")))
	fni.Focus()

	formats := []string{"wav", "mp3", "flac"}
	formatIndex := 0
	for i, f := range formats {
		if f == cfg.DefaultFormat {
			formatIndex = i
			break
		}
	}

	di := textinput.New()
	di.Placeholder = "C:\\Users\\..."
	di.SetValue(cfg.DefaultDirectory)

	return RecordModel{
		state:    StateStarting,
		config:   cfg,
		progress: p,
		waveform: make([]float64, 20),
		recorder: rec,
		fileNameInput: fni,
		formats:       formats,
		formatIndex:   formatIndex,
		dirInput:      di,
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m RecordModel) Init() tea.Cmd {
	return tick()
}

func (m RecordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.recorder != nil {
				m.recorder.Stop()
				m.recorder.Close()
			}
			return m, tea.Quit
		case "o":
			if m.state == StateFinished {
				_ = openFile(m.filePath)
				return m, nil
			}
		case "r":
			if m.state == StateFinished {
				// Reset for a new recording
				newModel := InitialRecordModel(m.config)
				return newModel, newModel.Init()
			}
		case "enter":
			if m.state == StateRecording {
				m.recorder.Stop()
				m.state = StateSaving
				return m, textinput.Blink
			} else if m.state == StateSaving {
				if m.focusIndex == 2 {
					m.fileName = m.fileNameInput.Value() + "." + m.formats[m.formatIndex]
					m.filePath = filepath.Join(m.dirInput.Value(), m.fileName)
					
					_, err := m.recorder.Save(m.filePath, m.formats[m.formatIndex])
					if err != nil {
						m.err = err
						return m, nil
					}

					// Save to history
					_ = config.AddToHistory(config.HistoryEntry{
						Timestamp: time.Now(),
						FileName:  m.fileName,
						FilePath:  m.filePath,
						Format:    m.formats[m.formatIndex],
					})

					m.state = StateFinished
					return m, nil
				}
				m.focusIndex++
				return m.updateFocus()
			}
		case "tab", "down":
			if m.state == StateSaving {
				m.focusIndex = (m.focusIndex + 1) % 3
				return m.updateFocus()
			}
		case "up":
			if m.state == StateSaving {
				m.focusIndex = (m.focusIndex - 1 + 3) % 3
				return m.updateFocus()
			}
		case "left", "right":
			if m.state == StateSaving && m.focusIndex == 1 {
				if msg.String() == "left" {
					m.formatIndex = (m.formatIndex - 1 + len(m.formats)) % len(m.formats)
				} else {
					m.formatIndex = (m.formatIndex + 1) % len(m.formats)
				}
				return m, nil
			}
		}

	case tickMsg:
		if m.state == StateStarting {
			m.tickCount++
			cmd := m.progress.SetPercent(float64(m.tickCount) / 20.0)
			if m.tickCount >= 20 {
				m.state = StateRecording
				m.startTime = time.Now()
				if m.recorder != nil {
					err := m.recorder.Start()
					if err != nil {
						m.err = err
					}
				}
				return m, tick()
			}
			return m, tea.Batch(cmd, tick())
		}

		if m.state == StateRecording {
			m.tickCount++
			m.blink = !m.blink
			
			// Use real amplitude
			amp := 0.0
			if m.recorder != nil {
				amp = float64(m.recorder.GetAmplitude())
			}
			
			// Shift waveform
			for i := 0; i < len(m.waveform)-1; i++ {
				m.waveform[i] = m.waveform[i+1]
			}
			m.waveform[len(m.waveform)-1] = amp * 5.0 // Scale for visibility
			
			return m, tick()
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	if m.state == StateSaving {
		var cmd tea.Cmd
		if m.focusIndex == 0 {
			m.fileNameInput, cmd = m.fileNameInput.Update(msg)
		} else if m.focusIndex == 2 {
			m.dirInput, cmd = m.dirInput.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

func (m RecordModel) updateFocus() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		cmd = m.fileNameInput.Focus()
		m.dirInput.Blur()
	case 1:
		m.fileNameInput.Blur()
		m.dirInput.Blur()
	case 2:
		m.fileNameInput.Blur()
		cmd = m.dirInput.Focus()
	}
	return m, cmd
}

func (m RecordModel) View() string {
	var s strings.Builder

	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	switch m.state {
	case StateStarting:
		s.WriteString(TitleStyle.Render("🚀 Initializing EchoMind..."))
		s.WriteString("\n\n")
		s.WriteString(m.progress.View())
		s.WriteString("\n\n")
		s.WriteString(StatusStyle.Render("Preparing your recording session..."))

	case StateRecording:
		s.WriteString(TitleStyle.Render("🎙️  Recording Session"))
		s.WriteString("\n\n")

		indicator := lipgloss.NewStyle().Foreground(ErrorColor).Render("● REC")
		if !m.blink {
			indicator = lipgloss.NewStyle().Foreground(MutedColor).Render("● REC")
		}

		duration := time.Since(m.startTime).Round(time.Second)
		timerStr := lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true).Render(duration.String())

		s.WriteString(fmt.Sprintf("%s   %s\n\n", indicator, timerStr))

		s.WriteString(PromptStyle.Render("Waveform:"))
		s.WriteString("\n")
		for _, v := range m.waveform {
			barLen := int(v * 30)
			if barLen > 30 {
				barLen = 30
			}
			if barLen < 1 && v > 0.01 {
				barLen = 1
			}
			s.WriteString(lipgloss.NewStyle().Foreground(MainColor).Render(strings.Repeat("█", barLen)) + "\n")
		}
		s.WriteString("\n")
		s.WriteString(StatusStyle.Render("Press ENTER to stop recording"))

	case StateSaving:
		s.WriteString(TitleStyle.Render("💾 Save Recording"))
		s.WriteString("\n\n")

		// File Name
		s.WriteString(m.renderField("File Name:", m.fileNameInput.View(), m.focusIndex == 0))
		s.WriteString("\n")

		// Format
		var formatView strings.Builder
		for i, f := range m.formats {
			style := lipgloss.NewStyle().Padding(0, 1)
			if i == m.formatIndex {
				style = style.Foreground(lipgloss.Color("#000000")).Background(MainColor).Bold(true)
			} else {
				style = style.Foreground(MutedColor)
			}
			formatView.WriteString(style.Render(strings.ToUpper(f)) + " ")
		}
		s.WriteString(m.renderField("File Format (Left/Right arrows):", formatView.String(), m.focusIndex == 1))
		s.WriteString("\n")

		// Directory
		s.WriteString(m.renderField("Save Directory:", m.dirInput.View(), m.focusIndex == 2))
		s.WriteString("\n")

		s.WriteString(StatusStyle.Render("(arrows to navigate, enter to confirm)"))

	case StateFinished:
		s.WriteString(TitleStyle.Render("✅ Recording Saved!"))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("%s %s\n", PromptStyle.Render("File:"), m.fileName))
		s.WriteString(fmt.Sprintf("%s %s\n", PromptStyle.Render("Location:"), filepath.Dir(m.filePath)))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(SecondaryColor).Render("Your voice has been captured in the digital void."))
		s.WriteString("\n\n")
		s.WriteString(StatusStyle.Render("Press 'o' to open, 'r' to record again, 'q' to exit"))
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s.String())
}

func openFile(path string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", "", path}
	case "darwin":
		cmd = "open"
		args = []string{path}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{path}
	}

	return exec.Command(cmd, args...).Start()
}

func (m RecordModel) renderField(label string, value string, focused bool) string {
	labelStyle := PromptStyle
	if focused {
		labelStyle = labelStyle.Copy().Foreground(lipgloss.Color("#FFF")).Background(MainColor).Padding(0, 1)
	}
	return fmt.Sprintf("%s\n%s\n", labelStyle.Render(label), value)
}
