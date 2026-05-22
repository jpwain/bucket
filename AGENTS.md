# AGENTS.md

Guidance for agentic development work in this repository.

## Scope

- This file applies to the whole repo unless a deeper `AGENTS.md` overrides it.
- Keep changes aligned with the existing Go app, release workflow, and Homebrew tap setup.

## Working Rules

- Prefer small, focused commits that group related changes.
- Use conventional commit messages https://www.conventionalcommits.org/en/v1.0.0/ -- with no body or footer unless the user confirms the content of those first.
- Do not revert user changes unless explicitly asked.
- Use `apply_patch` for manual file edits.
- Prefer `rg` for searches and `sed` for reads.
- Keep README content user-facing; keep release/process notes in `plans/`.

## Release Workflow

- Follow `plans/release-checklist.md` for release preparation.
- The release checklist covers tagging, GitHub Release assets, checksum updates, formula updates, and install verification.
- The installable Homebrew formula lives at `Formula/bucket.rb`.

## Project Layout

- `src/cmd/bucket/main.go` is the CLI entrypoint.
- `plans/` is for internal planning and release guidance.
- `Formula/bucket.rb` is the installable Homebrew formula for the repo tap.

## Verification

- Run `cd src && go test ./...` before committing meaningful Go changes.
- For release-related changes, verify the checklist steps still match the actual workflow.
