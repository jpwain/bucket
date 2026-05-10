# Bucket

## Goal

Build a terminal tool for reviewing two text files side by side and moving individual lines between them. The app should be simple, predictable, and optimized for line triage: load two files, move lines in memory, save changes intentionally.

The implementation should prioritize correct data behavior, clear keyboard interaction, and idiomatic Bubble Tea structure.

## Product Behavior

### Core Workflow

- The user launches the app with two file paths.
- The app loads both files as text and displays them as two side-by-side buckets.
- Each bucket has a cursor indicating the selected line.
- The user can move selected lines from one bucket to the other.
- Edits are held in memory until the user saves.
- The user can undo and redo line moves.
- The user can quit, with confirmation when there are unsaved changes.

### CLI Behavior

- The command requires exactly two file arguments.
- Each path is resolved to an absolute path.
- The two paths must refer to different files.
- Both files must be readable as text.
- Invalid arguments or load failures print an error and exit non-zero.

### Screen Layout

- A top status line shows both file names and their current line-count deltas.
- The main area shows two equal-width buckets: left file and right file.
- Each bucket shows line numbers, line text, and the selected line.
- A bottom hint line lists the primary keyboard commands.
- Save, quit, and help flows appear as modal dialogs or equivalent blocking prompts.

The layout requirements are structural only.

### Presentation

- Use Charm's Lip Gloss for terminal layout and presentation.
- Prefer Lip Gloss's declarative style model for boxes, padding, borders, alignment, width, height, and text emphasis.
- Keep styling idiomatic to Bubble Tea and Lip Gloss rather than defining a large custom visual system.
- Use terminal-aware styling that works in light themes, dark themes, limited-color terminals, and plain ANSI environments.
- Use visual treatment only to clarify function:
  - focused bucket vs unfocused bucket,
  - selected line,
  - muted metadata such as line numbers and hints,
  - success, error, and neutral messages,
  - active dialog option.
- Avoid carrying over fixed application-specific colors, backgrounds, or theme decisions unless they are introduced deliberately for this app.
- Favor simple, readable defaults over detailed visual customization.

### Status Information

- Each file shows a line-count delta relative to its current saved baseline.
- Positive deltas are formatted with a leading `+`.
- Zero deltas are shown as `0`.
- Transient messages may be shown for save results, errors, undo/redo outcomes, and no-op moves.
- The bottom hint line should communicate:
  - `tab switch`
  - `up/down select`
  - `left/right move`
  - `z undo`
  - `Z redo`
  - `w wrap`
  - `s save`
  - `q quit`
  - `? help`

## Keyboard Interaction

### Navigation

- `tab`: switch focused bucket.
- `up`: move the focused bucket cursor up.
- `down`: move the focused bucket cursor down.
- `shift+up`: move the unfocused bucket cursor up.
- `shift+down`: move the unfocused bucket cursor down.

Cursor movement clamps to valid line positions. Empty buckets keep their cursor at position `0`.

### Moving Lines

- `right`: move the selected line from the left bucket to the right bucket, inserted above the right cursor.
- `left`: move the selected line from the right bucket to the left bucket, inserted above the left cursor.
- `shift+right`: move the selected line from the left bucket to the right bucket, inserted below the right cursor.
- `shift+left`: move the selected line from the right bucket to the left bucket, inserted below the left cursor.

Move direction is determined by the key pressed, not by which bucket is focused.

When a move succeeds:

- The source line is removed.
- The destination receives the moved line at the requested insertion point.
- Source and destination cursors are clamped to valid positions.
- The destination cursor moves to the inserted line.
- The move is recorded in undo history.
- Any redo history is cleared.

Moving from an empty source bucket is a no-op and should not create history.

### Undo And Redo

- `z`: undo the last line move.
- `Z`: redo the last undone line move.

Undo restores the previous bucket contents, focus, and cursor positions for that move. Redo reapplies the move and restores the post-move interaction state.

The history stack should have a fixed limit, defaulting to 1000 moves.

### Wrap

- `w`: toggle line wrapping.

When wrapping is off, long lines are truncated to the visible bucket width. When wrapping is on, long lines wrap within their bucket. In both modes, cursor visibility should be preserved while navigating.

### Save

- `s`: open save confirmation.
- Confirming save writes only files whose serialized content differs from their saved baseline.
- Canceling save returns to the main view without writing.
- Save errors are shown to the user and do not update baselines.

After a successful save:

- The saved text becomes the new baseline for dirty detection.
- Line-count deltas are recalculated from the new baseline.
- The app remains open unless the save was part of a quit flow.

### Quit

- `q`: open quit confirmation.
- If there are unsaved changes, the user can save and quit, discard and quit, or cancel.
- If there are no unsaved changes, quitting should be direct or require only minimal confirmation.
- `ctrl+c`: exit immediately without saving.

### Help

- `?`: open help.
- Help lists the available commands and their effects.
- `esc` closes help and returns to the main view.

## Data Model

### Bucket

Each bucket stores:

- Absolute file path.
- Display name.
- Current lines.
- Saved baseline text.
- Saved baseline line count.
- Newline sequence used for serialization.
- Whether the baseline ended with a trailing newline.
- Current cursor position.

### App State

The app state stores:

- Left bucket.
- Right bucket.
- Focused side.
- Wrap mode.
- Undo and redo history.
- Current dialog or prompt mode.
- Transient message state.
- Terminal dimensions and scroll offsets.

### Move History Entry

Each successful move records:

- Source side.
- Destination side.
- Source index.
- Destination index.
- Moved line text.
- Interaction state before the move.
- Interaction state after the move.

Interaction state includes focused side and both cursor positions.

## File Semantics

### Loading

- Files are loaded as text.
- Lines are represented internally without newline terminators.
- Empty files produce an empty line slice.
- The app detects the newline sequence used for future serialization.
- The app records whether the file ended with a trailing newline.

### Serialization

- Lines are joined using the bucket's recorded newline sequence.
- The original trailing-newline behavior is preserved.
- Serialization is the source of truth for dirty detection.

### Dirty Detection

- A bucket is dirty when its serialized current text differs from its saved baseline text.
- Dirty detection should not depend only on line counts.

### Saving

- Save only dirty buckets.
- Write through a temporary file in the same directory.
- Flush and close the temporary file before rename.
- Rename the temporary file over the target path.
- After a successful write, update that bucket's saved baseline text and saved baseline line count.

## Implementation Shape

### Suggested Packages

- `cmd/buckets/main.go`: command entrypoint.
- `internal/cli`: argument parsing and validation.
- `internal/io`: file loading, newline handling, serialization helpers, atomic writes.
- `internal/domain`: buckets, moves, cursor logic, dirty detection, undo/redo.
- `internal/app`: Bubble Tea model, update loop, rendering, dialogs, scrolling.

### Library Choices

- Use Bubble Tea for the application loop.
- Use Lip Gloss for layout and presentation, including borders, padding, alignment, sizing, and lightweight text emphasis.
- Use Bubbles components only when they reduce complexity, such as viewport handling for long files.

Keep domain logic independent from Bubble Tea so movement, history, serialization, and save behavior can be tested without terminal rendering.

## Development Plan

### Phase 1: Foundation

- Initialize the Go module and command entrypoint.
- Parse and validate exactly two file arguments.
- Load files and build initial bucket state.
- Start a Bubble Tea program with a minimal main view.

### Phase 2: Domain Logic

- Implement line splitting and serialization.
- Implement cursor movement and clamping.
- Implement directional line moves.
- Implement dirty detection and line-count deltas.
- Add unit tests for the domain behavior.

### Phase 3: History

- Record successful moves.
- Implement undo and redo.
- Restore focus and cursor snapshots.
- Enforce the history limit.
- Clear redo history after a new move.

### Phase 4: Terminal App

- Render the top status line, two buckets, and bottom command hints.
- Implement keyboard handling.
- Implement wrapping and scroll behavior.
- Ensure rendered output fits within the terminal at common sizes.

### Phase 5: Dialogs

- Implement save confirmation.
- Implement quit confirmation.
- Implement help.
- Ensure `esc` and `ctrl+c` behavior is consistent.

### Phase 6: Persistence

- Save only changed files.
- Write atomically through temporary files.
- Update baselines after successful saves.
- Report save results and errors clearly.

### Phase 7: Hardening

- Test small terminal sizes, long files, long lines, empty files, and files without trailing newlines.
- Verify behavior when save partially fails.
- Verify symlink or hard-link handling for duplicate file detection.
- Manually test the demo workflow end to end.

## Testing Plan

### Unit Tests

- CLI argument validation.
- Same-file detection.
- Text splitting and serialization.
- Newline and trailing-newline preservation.
- Cursor clamping.
- Directional move behavior.
- Empty-source move no-op behavior.
- Undo and redo.
- Redo clearing after a new move.
- History limit enforcement.
- Dirty detection.
- Save baseline updates.

### Integration Tests

- File-system save tests using temporary directories.
- Atomic write behavior.
- Key-flow tests for move, undo, redo, save, help, and quit.
- Render sizing tests to ensure the view does not overflow terminal dimensions.

### Manual Checklist

- Launch with two demo files.
- Move lines left to right and right to left.
- Move above and below destination cursors.
- Navigate both focused and unfocused buckets.
- Toggle wrapping with long lines.
- Undo and redo several moves.
- Save changes and confirm files changed on disk.
- Quit with unsaved changes and test save, discard, and cancel.
- Open and close help.
- Test empty files and files without trailing newlines.

## Non-Goals

- Multi-select editing.
- Search or filtering.
- Mouse support.
- Plugin support.
- In-place text editing.
- Automatic saving.

## Acceptance Criteria

- The app can be run with two files and used entirely from the keyboard.
- Movement semantics are deterministic and independent of focus.
- Undo and redo preserve file contents and interaction state.
- Dirty detection is based on serialized text.
- Saving writes only changed files and updates baselines correctly.
- The terminal view remains structurally readable across common terminal sizes.
- Core behavior is covered by unit tests.
- The implementation is self-contained and follows the behavior described in this plan.
