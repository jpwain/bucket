package app

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"bucket/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestRenderPanelsMatchesTerminalWidth(t *testing.T) {
	m := model{
		state: domain.State{
			Left:  domain.Bucket{Name: "left.txt", Lines: []string{"a"}, Cursor: 0},
			Right: domain.Bucket{Name: "right.txt", Lines: []string{"b"}, Cursor: 0},
			Focus: domain.Left,
		},
		height: 24,
	}

	for _, width := range []int{60, 61, 80, 81, 120, 121} {
		m.width = width
		got := lipgloss.Width(m.renderPanels())
		if got != width {
			t.Fatalf("width %d rendered as %d", width, got)
		}
	}
}

func TestScrollAdvancesWhenCursorMovesPastViewport(t *testing.T) {
	m := model{
		state: domain.State{
			Left: domain.Bucket{
				Name:   "left.txt",
				Lines:  makeLines("left", 12),
				Cursor: 0,
			},
			Right: domain.Bucket{
				Name:   "right.txt",
				Lines:  []string{"only"},
				Cursor: 0,
			},
			Focus: domain.Left,
		},
		width:  80,
		height: 12,
	}

	for i := 0; i < 9; i++ {
		next, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = next.(model)
	}

	if got := m.state.Left.Scroll; got != 2 {
		t.Fatalf("scroll = %d, want 2", got)
	}

	rendered := stripANSI(m.renderBucket(m.state.Left, 38, 10, true))
	if !strings.Contains(rendered, "left-10") {
		t.Fatalf("expected scrolled view to include left-10, got:\n%s", rendered)
	}
	if strings.Contains(rendered, "left-01") {
		t.Fatalf("expected top line to be scrolled out, got:\n%s", rendered)
	}
}

func TestWrapKeepsBucketHeightFixed(t *testing.T) {
	m := model{
		state: domain.State{
			Left: domain.Bucket{
				Name:   "left.txt",
				Lines:  []string{strings.Repeat("wrap me ", 12)},
				Cursor: 0,
			},
			Right: domain.Bucket{
				Name:   "right.txt",
				Lines:  []string{"short"},
				Cursor: 0,
			},
			Focus: domain.Left,
			Wrap:  true,
		},
	}

	rendered := m.renderBucket(m.state.Left, 38, 8, true)
	if got := lipgloss.Height(rendered); got != 8 {
		t.Fatalf("wrapped height = %d, want 8", got)
	}
	if !strings.Contains(stripANSI(rendered), "wrap me") {
		t.Fatalf("expected wrapped text to be present, got:\n%s", rendered)
	}
}

func makeLines(prefix string, n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("%s-%02d", prefix, i+1)
	}
	return lines
}

func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[[:alpha:]]`)
	return re.ReplaceAllString(s, "")
}
