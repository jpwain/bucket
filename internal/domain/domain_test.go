package domain

import "testing"

func TestMoveAndUndoRedo(t *testing.T) {
	s := State{
		Left:         Bucket{Lines: []string{"a", "b"}, Cursor: 0, SavedBaselineCount: 2},
		Right:        Bucket{Lines: []string{"x"}, Cursor: 0, SavedBaselineCount: 1},
		Focus:        Left,
		HistoryLimit: 1000,
	}
	ok := s.MoveSelected(Left, Right, false)
	if !ok {
		t.Fatal("expected move")
	}
	if len(s.Left.Lines) != 1 || s.Left.Lines[0] != "b" {
		t.Fatalf("left: %#v", s.Left.Lines)
	}
	if len(s.Right.Lines) != 2 || s.Right.Lines[0] != "a" {
		t.Fatalf("right: %#v", s.Right.Lines)
	}
	if !s.UndoMove() {
		t.Fatal("expected undo")
	}
	if len(s.Left.Lines) != 2 || s.Left.Lines[0] != "a" {
		t.Fatalf("undo left: %#v", s.Left.Lines)
	}
	if !s.RedoMove() {
		t.Fatal("expected redo")
	}
	if len(s.Right.Lines) != 2 || s.Right.Lines[0] != "a" {
		t.Fatalf("redo right: %#v", s.Right.Lines)
	}
}

func TestMoveNoopOnEmptySource(t *testing.T) {
	s := State{
		Left:         Bucket{Lines: []string{}, Cursor: 0},
		Right:        Bucket{Lines: []string{"x"}, Cursor: 0},
		HistoryLimit: 1000,
	}
	if s.MoveSelected(Left, Right, false) {
		t.Fatal("expected no-op")
	}
	if len(s.Undo) != 0 {
		t.Fatal("expected no history")
	}
}
