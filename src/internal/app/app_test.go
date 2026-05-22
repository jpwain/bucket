package app

import (
	"testing"

	"bucket/internal/domain"

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
