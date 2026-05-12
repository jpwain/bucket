# Bucket Publishing Plan

## Goal

Ship `bucket` as a single Go binary that is easy to build, release, and install with Homebrew. Keep the process low overhead: no runtime dependencies, no extra services, and no release machinery that is more complicated than the app itself.

## Core Decisions

- Keep the app as a standard Go module with one executable at `cmd/bucket`.
- Build a single statically linked binary where practical.
- Publish tagged releases and use those tags as the source of truth for Homebrew.
- Avoid adding release tooling unless it clearly reduces manual work.
- Update `sha256` in the Homebrew formula per release if that remains simple.

## Human Actions

1. Create a release versioning convention.
   - Use semantic version tags such as `v0.1.0`.
   - Treat the Git tag as the release identity.

2. Decide the release artifact shape.
   - Prefer GitHub release assets containing OS-specific tarballs or zip files.
   - Keep the artifact contents minimal: binary plus README/license if needed.

3. Choose the Homebrew distribution strategy.
   - Prefer a personal tap or organization tap for formulas.
   - Keep the formula in a separate repo if that makes releases easier to manage.

4. Prepare release notes.
   - Summarize visible user-facing changes.
   - Call out any terminal behavior or keybinding changes.

5. Publish a release.
   - Tag the source tree.
   - Build release artifacts.
   - Upload artifacts to GitHub Releases.
   - Update the Homebrew formula with the new version and checksum.

6. Test the published install path.
   - Install with Homebrew from the tap.
   - Run the binary with the demo files.
   - Confirm the versioned release behaves like the local build.

## Coding Agent Actions

1. Keep the build entrypoint stable.
   - Ensure `cmd/bucket/main.go` remains the single executable entrypoint.
   - Keep package layout idiomatic and predictable.

2. Add release-friendly build support.
   - Make sure `go build ./cmd/bucket` succeeds cleanly.
   - Prefer `CGO_ENABLED=0` if the app continues to work without CGO.
   - Use `-trimpath` and `-ldflags="-s -w"` if they do not interfere with debugging or version reporting.

3. Add version injection support if useful.
   - Expose a build-time version string if the binary should report its version.
   - Keep this minimal and optional.

4. Prepare Homebrew formula support.
   - Add a formula template or documentation for `bucket`.
   - Keep the formula small and standard.
   - Prefer the formula to install a prebuilt release artifact if available.
   - If source builds are simpler, keep the formula using `system "go", "build", ...` and document the tradeoff.

5. Make checksum updates straightforward.
   - Document where the release `sha256` lives.
   - If possible, generate the checksum as part of release packaging.
   - Keep the checksum update process one manual edit or one script run at most.

6. Provide a release checklist in the repo.
   - Document the exact commands used to build, package, and verify a release.
   - Include the Homebrew update step and install verification.

## Suggested Build Shape

```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bucket ./cmd/bucket
```

This is the preferred shape because it keeps the binary simple, portable, and easy to publish.

## Suggested Release Flow

1. Update code.
2. Run tests locally.
3. Tag a version.
4. Build release binaries for the target platforms.
5. Upload release assets.
6. Update the Homebrew formula with the new version and `sha256`.
7. Test `brew install` from the tap.

## Suggested Homebrew Formula Shape

- Formula name: `bucket`
- Tap: `your-org/bucket` or similar
- Install method: either prebuilt tarball download or source build from the release tag
- Verification: `bucket --help` or a small demo invocation

## Non-Goals

- Automatic update servers.
- In-app self-updating.
- Complex release orchestration.
- Cross-platform installers beyond Homebrew and straightforward binary downloads.

## Acceptance Criteria

- A tagged release can be built into a single `bucket` binary.
- The binary can be installed with Homebrew from a tap.
- The Homebrew formula is small and maintainable.
- Release checksum updates are straightforward.
- The release process remains simple enough to repeat manually without specialized infrastructure.
