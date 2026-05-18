# Bucket Release Checklist

Use this when shipping a tagged release.

1. Update the code and run tests.
   ```bash
   go test ./...
   ```

2. Bump the version in the release notes or changelog if you keep one.

3. Create the release tag from a clean working tree.
   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

4. Wait for the GitHub Actions release workflow to finish.
   - Workflow: `.github/workflows/release.yml`
   - Expected assets: `bucket_X.Y.Z_<goos>_<goarch>.tar.gz`

5. Download the macOS release asset you need for the tap update.
   ```bash
   gh release download vX.Y.Z --repo jpwain/bucket --pattern 'bucket_X.Y.Z_darwin_amd64.tar.gz'
   shasum -a 256 bucket_X.Y.Z_darwin_amd64.tar.gz
   ```

6. Update the Homebrew formula.
   - Edit `Formula/bucket.rb` in this repo.
   - Replace the version with the new release version.
   - Replace the release URL with the new asset URL.
   - Replace the `sha256` value with the checksum from the downloaded tarball.

7. Commit and push the formula update.

8. Verify the published install path.
   ```bash
   brew tap jpwain/bucket https://github.com/jpwain/bucket.git
   brew update
   brew install jpwain/bucket/bucket
   bucket --version
   ```

9. Run a real usage check with the sample files you trust.
   ```bash
   bucket left.txt right.txt
   ```

10. If the version output or binary behavior looks wrong, fix the tap formula or release artifact before announcing the release.
