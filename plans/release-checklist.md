# Bucket Release Checklist

Use this when shipping a tagged release.

1. Update the code and run tests.
   ```bash
   go test ./...
   ```

2. Bump the version in the release notes or changelog if you keep one.

3. Create the release tag from a clean working tree.
   ```bash
   git tag v1.1.0
   git push origin v1.1.0
   ```

4. Wait for the GitHub Actions release workflow to finish.
   - Workflow: `.github/workflows/release.yml`
   - Expected assets: `bucket_1.1.0_<goos>_<goarch>.tar.gz`

5. Download the macOS release asset you need for the tap update.
   ```bash
   gh release download v1.1.0 --repo jpwain/bucket --pattern 'bucket_1.1.0_darwin_amd64.tar.gz'
   shasum -a 256 bucket_1.1.0_darwin_amd64.tar.gz
   ```

6. Update the Homebrew formula in the tap repo.
   - Replace the version with `1.1.0`.
   - Replace the release URL with the new asset URL.
   - Replace the `sha256` value with the checksum from the downloaded tarball.

7. Commit and push the tap update, or open a pull request in the tap repo.

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
