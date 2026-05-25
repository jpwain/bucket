# Bucket Release Checklist

Use this when shipping the latest tagged release or when refreshing the tap formula against the latest published release.

1. Update the code and run tests.
   ```bash
   cd src && go test ./...
   ```

2. For a new release, create the next release tag from a clean working tree.
   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

3. Wait for the GitHub Actions release workflow to finish.
   - Workflow: `.github/workflows/release.yml`
   - Expected assets: `bucket_X.Y.Z_<goos>_<goarch>.tar.gz`

4. For the release tag you are targeting, fetch the asset digests from GitHub Release metadata.
   ```bash
   gh release view vX.Y.Z --repo jpwain/bucket --json assets --jq '.assets[] | {name: .name, digest: .digest}'
   ```

5. Update the Homebrew formula.
   - Edit `Formula/bucket.rb` in this repo.
   - Keep the version pinned to the release tag being referenced.
   - Set the release URL for each supported OS and architecture.
   - Replace the `sha256` values with the published release digests for each tarball.

6. Commit and push the formula update.

7. Verify the published install path.
   ```bash
   brew tap jpwain/bucket https://github.com/jpwain/bucket.git
   brew update
   brew install jpwain/bucket/bucket
   bucket --version
   ```

8. Run a real usage check with the sample files you trust.
   ```bash
   bucket left.txt right.txt
   ```

9. If the version output or binary behavior looks wrong, fix the tap formula or release artifact before announcing the release.
