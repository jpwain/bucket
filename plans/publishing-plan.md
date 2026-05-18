# Bucket Publishing Plan

## Goal

Ship `bucket` as a single Go binary that is easy to build, release, and install with Homebrew. Keep the process low overhead: no runtime dependencies, no extra services, and no release machinery that is more complicated than the app itself.

## Core Decisions

- Keep the app as a standard Go module with one executable at `cmd/bucket`.
- Build a single release binary where practical, with CGO disabled if the app continues to work without it.
- Publish tagged GitHub releases with prebuilt binaries and use those tags as the source of truth for Homebrew.
- Keep the release pipeline in GitHub Actions rather than adding a separate release tool unless the manual steps become painful.
- Update the Homebrew formula `sha256` per release.
- Keep release-only build flags explicit so debug builds stay available when needed.

## Human Actions

1. Create a release versioning convention.
   - Use semantic version tags such as `v1.1.0`.
   - Treat the Git tag as the release identity and the only version source Homebrew should reference.
   - Decide whether prereleases should be published or skipped.

2. Choose the Homebrew distribution strategy.
   - Keep the formula in the repo root `Formula/bucket.rb` so the repository can be tapped directly with the explicit GitHub URL.
   - The formula should install a prebuilt release tarball or zip, not build from source.

3. Prepare release notes.
   - Summarize visible user-facing changes.
   - Call out any terminal behavior, keybinding changes, or compatibility notes.

4. Publish a release.
   - Create and push the tag from a clean working tree.
   - Let GitHub Actions build the release artifacts for the supported platforms.
   - Publish a GitHub Release for that tag.
   - Update the Homebrew formula with the new version and checksum in the tap repo.
   - The release workflow lives at `.github/workflows/release.yml`.
   - Release asset names should follow `bucket_<version>_<goos>_<goarch>.tar.gz`.

5. Test the published install path.
   - Install with Homebrew from the tap.
   - Run `bucket --version`.
   - Run the binary against the demo files and confirm it matches the local build.

## Coding Agent Actions

1. Keep the build entrypoint stable.
   - Ensure `cmd/bucket/main.go` remains the single executable entrypoint.
   - Keep package layout idiomatic and predictable.

2. Add release-friendly build support.
   - Make sure `go build ./cmd/bucket` succeeds cleanly.
   - Add a documented release build command that includes `CGO_ENABLED=0`, `-trimpath`, and `-ldflags="-s -w"` if the app continues to work without CGO.
   - Keep a separate non-stripped build path available for debugging.

3. Add version injection support if useful.
   - Expose a build-time version string if the binary should report its version.
   - Keep this minimal and optional.
   - Default the version to `dev` so local builds behave sensibly before release injection is wired in.

4. Prepare Homebrew formula support.
   - Keep a formula template or documentation for `bucket` at `homebrew/Formula/bucket.rb` in this repo, and keep the installable tap formula at `Formula/bucket.rb`.
   - Keep the formula small and standard.
   - Prefer the formula to install a prebuilt release artifact.
   - Document where the tarball or zip is uploaded and how the URL is formed from the tag.
   - Keep the formula pinned to the release artifact checksum instead of building from source.
   - Keep the formula version in plain semver form, matching the release binary's `--version` output.
   - Use the tap repo URL `https://github.com/jpwain/bucket.git`.

5. Make checksum updates straightforward.
   - Document where the release `sha256` lives.
   - If possible, generate the checksum as part of release packaging or as a GitHub Actions step.
   - Keep the checksum update process one manual edit or one script run at most.

6. Provide a release checklist in the repo.
   - Document the exact commands used to build, package, and verify a release.
   - Include the Git tag, GitHub Release, checksum, Homebrew update, and install verification steps.
   - Include the exact command to verify `bucket --version` from the installed formula.

7. Add a GitHub Actions release workflow.
   - Build from a tag push only.
   - Run tests before packaging.
   - Build the release artifacts with the documented release flags.
   - Package one archive per target OS and architecture.
   - Upload artifacts to the GitHub Release.
   - Emit the checksum needed by the tap formula update.

## Suggested Build Shape

```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bucket ./cmd/bucket
```

This is the preferred release shape because it keeps the binary simple, portable, and easier to publish.

For local debugging, use a plain `go build ./cmd/bucket` so symbols and file paths stay available.

## Suggested Release Flow

1. Update code.
2. Run tests locally.
3. Tag a version with `vX.Y.Z`.
4. Push the tag to GitHub.
5. Let the release workflow build the binaries, package them, and create the GitHub Release.
6. Update the Homebrew formula with the new version, archive URL, and `sha256`.
7. Test `brew install` from the tap and confirm `bucket --version`.

## Suggested Homebrew Formula Shape

- Formula name: `bucket`
- Tap repo: `jpwain/bucket`
- Tap URL: `https://github.com/jpwain/bucket.git`
- Install method: prebuilt tarball or zip download from the GitHub Release
- Verification: `bucket --help` or a small demo invocation
- Keep the URL and checksum derived from the tag so the update path is mechanical.
- If multiple platforms are supported, document the matching asset name pattern for each.

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
