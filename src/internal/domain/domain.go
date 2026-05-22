package domain

type Side int

const (
	Left Side = iota
	Right
)

type Bucket struct {
	Path                string
	Name                string
	Lines               []string
	SavedBaselineText   string
	SavedBaselineCount  int
	Newline             string
	BaselineHadTrailing bool
	Cursor              int
	Scroll              int
}

type InteractionState struct {
	Focus       Side
	LeftCursor  int
	RightCursor int
}

type MoveEntry struct {
	Source      Side
	Destination Side
	SourceIndex int
	DestIndex   int
	Line        string
	Before      InteractionState
	After       InteractionState
	BeforeLeft  []string
	BeforeRight []string
	AfterLeft   []string
	AfterRight  []string
}

type State struct {
	Left         Bucket
	Right        Bucket
	Focus        Side
	Wrap         bool
	Undo         []MoveEntry
	Redo         []MoveEntry
	HistoryLimit int
}

func ClampCursor(lines []string, cur int) int {
	if len(lines) == 0 {
		return 0
	}
	if cur < 0 {
		return 0
	}
	if cur >= len(lines) {
		return len(lines) - 1
	}
	return cur
}

func (s *State) switchFocus() {
	if s.Focus == Left {
		s.Focus = Right
	} else {
		s.Focus = Left
	}
}

func (s *State) MoveCursor(side Side, delta int) {
	b := s.bucket(side)
	b.Cursor = ClampCursor(b.Lines, b.Cursor+delta)
}

func (s *State) MoveSelected(source, destination Side, below bool) bool {
	src := s.bucket(source)
	dst := s.bucket(destination)
	if len(src.Lines) == 0 {
		return false
	}

	before := s.snapshot()
	beforeLeft := clone(s.Left.Lines)
	beforeRight := clone(s.Right.Lines)
	srcIdx := ClampCursor(src.Lines, src.Cursor)
	line := src.Lines[srcIdx]
	src.Lines = append(src.Lines[:srcIdx], src.Lines[srcIdx+1:]...)
	src.Cursor = ClampCursor(src.Lines, srcIdx)

	dstInsert := ClampCursor(dst.Lines, dst.Cursor)
	if len(dst.Lines) == 0 {
		dstInsert = 0
	} else if below {
		dstInsert++
		if dstInsert > len(dst.Lines) {
			dstInsert = len(dst.Lines)
		}
	}
	dst.Lines = insertLine(dst.Lines, dstInsert, line)
	dst.Cursor = dstInsert

	s.Left.Cursor = ClampCursor(s.Left.Lines, s.Left.Cursor)
	s.Right.Cursor = ClampCursor(s.Right.Lines, s.Right.Cursor)

	after := s.snapshot()
	entry := MoveEntry{
		Source:      source,
		Destination: destination,
		SourceIndex: srcIdx,
		DestIndex:   dstInsert,
		Line:        line,
		Before:      before,
		After:       after,
		BeforeLeft:  beforeLeft,
		BeforeRight: beforeRight,
		AfterLeft:   clone(s.Left.Lines),
		AfterRight:  clone(s.Right.Lines),
	}
	s.pushUndo(entry)
	s.Redo = nil
	return true
}

func (s *State) UndoMove() bool {
	if len(s.Undo) == 0 {
		return false
	}
	entry := s.Undo[len(s.Undo)-1]
	s.Undo = s.Undo[:len(s.Undo)-1]
	s.Left.Lines = clone(entry.BeforeLeft)
	s.Right.Lines = clone(entry.BeforeRight)
	s.applySnapshot(entry.Before)
	s.Redo = append(s.Redo, entry)
	return true
}

func (s *State) RedoMove() bool {
	if len(s.Redo) == 0 {
		return false
	}
	entry := s.Redo[len(s.Redo)-1]
	s.Redo = s.Redo[:len(s.Redo)-1]
	s.Left.Lines = clone(entry.AfterLeft)
	s.Right.Lines = clone(entry.AfterRight)
	s.applySnapshot(entry.After)
	s.pushUndo(entry)
	return true
}

func (s *State) LineDelta(side Side) int {
	b := s.bucket(side)
	return len(b.Lines) - b.SavedBaselineCount
}

func (s *State) bucket(side Side) *Bucket {
	if side == Left {
		return &s.Left
	}
	return &s.Right
}

func (s *State) snapshot() InteractionState {
	return InteractionState{
		Focus:       s.Focus,
		LeftCursor:  s.Left.Cursor,
		RightCursor: s.Right.Cursor,
	}
}

func (s *State) applySnapshot(i InteractionState) {
	s.Focus = i.Focus
	s.Left.Cursor = ClampCursor(s.Left.Lines, i.LeftCursor)
	s.Right.Cursor = ClampCursor(s.Right.Lines, i.RightCursor)
}

func (s *State) pushUndo(entry MoveEntry) {
	s.Undo = append(s.Undo, entry)
	if s.HistoryLimit <= 0 {
		s.HistoryLimit = 1000
	}
	if len(s.Undo) > s.HistoryLimit {
		s.Undo = s.Undo[len(s.Undo)-s.HistoryLimit:]
	}
}

func insertLine(lines []string, at int, line string) []string {
	if at < 0 {
		at = 0
	}
	if at > len(lines) {
		at = len(lines)
	}
	lines = append(lines, "")
	copy(lines[at+1:], lines[at:])
	lines[at] = line
	return lines
}

func clone(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	return out
}
