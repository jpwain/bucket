package app

import (
	"fmt"
	"strings"

	"bucket/internal/domain"
	fileio "bucket/internal/io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/colors"
)

type dialogMode int

const (
	dialogNone dialogMode = iota
	dialogHelp
	dialogSave
	dialogQuit
)

type model struct {
	state    domain.State
	width    int
	height   int
	dialog   dialogMode
	status   string
	statusOk bool
}

func Run(left, right domain.Bucket) error {
	m := model{
		state: domain.State{
			Left:         left,
			Right:        right,
			Focus:        domain.Left,
			Wrap:         false,
			HistoryLimit: 1000,
		},
	}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if m.dialog != dialogNone {
			return m.handleDialog(msg)
		}
		switch msg.String() {
		case "tab":
			if m.state.Focus == domain.Left {
				m.state.Focus = domain.Right
			} else {
				m.state.Focus = domain.Left
			}
		case "up":
			m.state.MoveCursor(m.state.Focus, -1)
		case "down":
			m.state.MoveCursor(m.state.Focus, 1)
		case "shift+up":
			m.state.MoveCursor(other(m.state.Focus), -1)
		case "shift+down":
			m.state.MoveCursor(other(m.state.Focus), 1)
		case "left":
			m.move(domain.Right, domain.Left, false)
		case "right":
			m.move(domain.Left, domain.Right, false)
		case "shift+left":
			m.move(domain.Right, domain.Left, true)
		case "shift+right":
			m.move(domain.Left, domain.Right, true)
		case "z":
			if !m.state.UndoMove() {
				m.status = "nothing to undo"
				m.statusOk = false
			} else {
				m.status = "undo"
				m.statusOk = true
			}
		case "Z":
			if !m.state.RedoMove() {
				m.status = "nothing to redo"
				m.statusOk = false
			} else {
				m.status = "redo"
				m.statusOk = true
			}
		case "w":
			m.state.Wrap = !m.state.Wrap
		case "s":
			m.dialog = dialogSave
		case "q":
			m.dialog = dialogQuit
		case "?":
			m.dialog = dialogHelp
		}
	}
	return m, nil
}

func (m *model) handleDialog(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.dialog {
	case dialogHelp:
		if k.String() == "esc" || k.String() == "q" || k.String() == "enter" {
			m.dialog = dialogNone
		}
	case dialogSave:
		switch k.String() {
		case "y", "enter":
			if err := fileio.SaveDirty(&m.state.Left); err != nil {
				m.status = fmt.Sprintf("save left failed: %v", err)
				m.statusOk = false
			} else if err := fileio.SaveDirty(&m.state.Right); err != nil {
				m.status = fmt.Sprintf("save right failed: %v", err)
				m.statusOk = false
			} else {
				m.status = "saved"
				m.statusOk = true
			}
			m.dialog = dialogNone
		case "n", "esc":
			m.dialog = dialogNone
		}
	case dialogQuit:
		switch k.String() {
		case "y":
			return m, tea.Quit
		case "s":
			if err := fileio.SaveDirty(&m.state.Left); err != nil {
				m.status = fmt.Sprintf("save left failed: %v", err)
				m.statusOk = false
				m.dialog = dialogNone
				return m, nil
			}
			if err := fileio.SaveDirty(&m.state.Right); err != nil {
				m.status = fmt.Sprintf("save right failed: %v", err)
				m.statusOk = false
				m.dialog = dialogNone
				return m, nil
			}
			return m, tea.Quit
		case "n":
			return m, tea.Quit
		case "esc", "c":
			m.dialog = dialogNone
		}
	}
	return m, nil
}

func (m *model) move(source, dest domain.Side, below bool) {
	if !m.state.MoveSelected(source, dest, below) {
		m.status = "no line to move"
		m.statusOk = false
		return
	}
	m.status = "moved line"
	m.statusOk = true
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	top := m.renderStatus()
	panels := m.renderPanels()
	bottom := m.renderHints()
	body := lipgloss.JoinVertical(lipgloss.Left, top, panels, bottom)
	if m.dialog != dialogNone {
		body += "\n" + m.renderDialog()
	}
	return body
}

func (m model) renderStatus() string {
	fileStyle := lipgloss.NewStyle().Foreground(colors.Normal)
	dotStyle := lipgloss.NewStyle().Foreground(colors.Gray)
	cmdStyle := lipgloss.NewStyle().Foreground(colors.Gray)
	leftDelta := m.state.LineDelta(domain.Left)
	rightDelta := m.state.LineDelta(domain.Right)
	status := fileStyle.Render(m.state.Left.Name) + " " + renderDelta(leftDelta) +
		" " + dotStyle.Render("·") + " " +
		fileStyle.Render(m.state.Right.Name) + " " + renderDelta(rightDelta)
	if m.status != "" {
		status += " " + cmdStyle.Render(":: "+m.status)
	}
	return status
}

func (m model) renderPanels() string {
	// Width() is measured before the border is drawn, but padding is already
	// accounted for inside the rendered block. Reserve only the outer borders.
	leftW := (m.width - 4) / 2
	rightW := m.width - 4 - leftW
	if leftW < 1 {
		leftW = 1
	}
	if rightW < 1 {
		rightW = 1
	}
	h := m.height - 4
	if h < 8 {
		h = 8
	}
	left := m.renderBucket(m.state.Left, leftW, h, m.state.Focus == domain.Left)
	right := m.renderBucket(m.state.Right, rightW, h, m.state.Focus == domain.Right)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderBucket(b domain.Bucket, width, height int, focused bool) string {
	border := lipgloss.NormalBorder()
	s := lipgloss.NewStyle().
		Border(border).
		Width(width).
		Height(height).
		Padding(0, 1)
	if focused {
		s = s.Bold(true).BorderForeground(colors.GrayBright)
	} else {
		s = s.BorderForeground(colors.GrayBrightDim)
	}
	lines := make([]string, 0, len(b.Lines)+1)
	lineNoStyle := lipgloss.NewStyle().Foreground(colors.GrayMid)
	rowStyle := lipgloss.NewStyle().Foreground(colors.NormalDim)
	if focused {
		rowStyle = lipgloss.NewStyle().Foreground(colors.Normal)
	}
	selectedStyle := lipgloss.NewStyle().
		Background(colors.GrayDark).
		Foreground(colors.NormalDim)
	if focused {
		selectedStyle = lipgloss.NewStyle().
			Background(colors.IndigoSubtle).
			Foreground(colors.WhiteBright)
	}
	if len(b.Lines) == 0 {
		lines = append(lines, "  [empty]")
	}
	for i, line := range b.Lines {
		prefix := "  " + fmt.Sprintf("%4d ", i+1)
		if i == b.Cursor {
			prefix = "> " + fmt.Sprintf("%4d ", i+1)
		}
		plainLine := line
		if !m.state.Wrap {
			maxText := width - 2 - len(prefix)
			if maxText < 0 {
				maxText = 0
			}
			if len(plainLine) > maxText {
				plainLine = plainLine[:maxText]
			}
		}
		row := lineNoStyle.Render(prefix) + rowStyle.Render(plainLine)
		if i == b.Cursor {
			row = selectedStyle.Render(prefix + plainLine)
		}
		lines = append(lines, row)
	}
	return s.Render(strings.Join(lines, "\n"))
}

func (m model) renderHints() string {
	keyStyle := lipgloss.NewStyle().Foreground(colors.Normal)
	descStyle := lipgloss.NewStyle().Foreground(colors.Gray)
	sepStyle := lipgloss.NewStyle().Foreground(colors.Gray)

	parts := []string{
		keyStyle.Render("tab") + " " + descStyle.Render("switch"),
		keyStyle.Render("↕") + " " + descStyle.Render("select"),
		keyStyle.Render("↔") + " " + descStyle.Render("move"),
		keyStyle.Render("z") + " " + descStyle.Render("undo"),
		keyStyle.Render("Z") + " " + descStyle.Render("redo"),
		keyStyle.Render("w") + " " + descStyle.Render("wrap"),
		keyStyle.Render("s") + " " + descStyle.Render("save"),
		keyStyle.Render("q") + " " + descStyle.Render("quit"),
		keyStyle.Render("?") + " " + descStyle.Render("help"),
	}
	return strings.Join(parts, sepStyle.Render("  "))
}

func (m model) renderDialog() string {
	switch m.dialog {
	case dialogHelp:
		return "Help: tab switch, up/down move cursor, shift+up/down move other cursor, left/right move, shift+left/right move below, z undo, Z redo, w wrap, s save, q quit, esc close."
	case dialogSave:
		return "Save changes? [y]es / [n]o"
	case dialogQuit:
		if fileio.IsDirty(m.state.Left) || fileio.IsDirty(m.state.Right) {
			return "Unsaved changes. [s] save+quit / [n] discard+quit / [c]ancel"
		}
		return "Quit? [y]es / [c]ancel"
	}
	return ""
}

func other(s domain.Side) domain.Side {
	if s == domain.Left {
		return domain.Right
	}
	return domain.Left
}

func formatDelta(v int) string {
	if v > 0 {
		return fmt.Sprintf("+%d", v)
	}
	return fmt.Sprintf("%d", v)
}

func renderDelta(v int) string {
	parenStyle := lipgloss.NewStyle().Foreground(colors.GrayMid)
	valueStyle := lipgloss.NewStyle().Foreground(colors.GrayMid)
	if v > 0 {
		valueStyle = lipgloss.NewStyle().Foreground(colors.Green)
	} else if v < 0 {
		valueStyle = lipgloss.NewStyle().Foreground(colors.Red)
	}
	return parenStyle.Render("(") + valueStyle.Render(formatDelta(v)) + parenStyle.Render(")")
}
