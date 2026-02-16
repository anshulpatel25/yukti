# GitHub Workflows Documentation

This repository contains a GitHub Actions workflow for automated release management.

## Workflow

### Draft Release Workflow (`draft-release.yml`)

This workflow creates draft releases with ARM64 binaries and can be triggered manually.

#### Features

- **Automatic Semantic Versioning**: Analyzes commit messages to determine version bumps
  - `breaking:` or `major:` prefix → Major version bump (X.0.0)
  - `feat:` or `feature:` prefix → Minor version bump (0.X.0)
  - Any other commit → Patch version bump (0.0.X)
  - Initial release → Version 0.1.0
- **Release Notes Generation**: Automatically generates release notes from commit history
- **ARM64 Binary Build**: Cross-compiles optimized ARM64 binary for Linux
- **Archive with Permissions**: Creates tar.gz archive preserving executable permissions
- **Draft Release**: Creates a draft release on GitHub with the binary attached

#### Usage

This workflow is triggered manually via GitHub Actions UI:

1. Go to the **Actions** tab in your GitHub repository
2. Select **Draft Release** from the workflows list
3. Click **Run workflow**
4. (Optional) Specify a branch or leave as `main`
5. Click **Run workflow** button

The workflow will:
1. Build the ARM64 binary
2. Calculate the next version based on commits
3. Generate release notes
4. Create a draft release with the binary attached

#### Inputs

- `branch` (optional): The branch to create the release from. Default: `main`

#### Requirements

- The workflow must be triggered from the `main` branch
- Repository must have `contents: write` permissions

### Version Numbering

The workflow follows [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- Initial version: **0.1.0**

#### Commit Message Conventions

To control version bumping, use these prefixes in your commit messages:

- `breaking:` or `major:` - Increment major version (e.g., 1.0.0 → 2.0.0)
  ```
  breaking: Remove deprecated API endpoint
  ```

- `feat:` or `feature:` - Increment minor version (e.g., 1.0.0 → 1.1.0)
  ```
  feat: Add new battery monitoring feature
  ```

- Any other commit - Increment patch version (e.g., 1.0.0 → 1.0.1)
  ```
  fix: Resolve charging threshold bug
  docs: Update README
  ```

### Release Process

1. **Make your changes** and commit them with appropriate prefixes
2. **Push to main branch**
3. **Trigger the workflow** manually via GitHub Actions UI
4. **Review the draft release** in the Releases page
5. **Edit if needed** (title, description, etc.)
6. **Publish the release** when ready

### Binary Artifact

The workflow produces:
- **yukti-arm64-linux.tar.gz**: ARM64 Linux binary (1.7MB uncompressed, ~660KB compressed)
  - Compiled with: `GOOS=linux GOARCH=arm64`
  - Optimized with: `-ldflags="-s -w"` (stripped symbols)
  - Executable permissions preserved in archive

#### Installation from Release

Users can download and install the binary:

```bash
# Download the release artifact
wget https://github.com/anshulpatel25/yukti/releases/download/v0.1.0/yukti-arm64-linux.tar.gz

# Extract
tar -xzf yukti-arm64-linux.tar.gz

# Move to system path (optional)
sudo mv yukti /usr/local/bin/yukti

# Run (requires root privileges)
sudo yukti
```

## Technical Details

### Actions Used

All actions are from trusted sources:

- **actions/checkout@v6**: Official GitHub action for repository checkout
- **actions/setup-go@v6**: Official GitHub action for Go setup
- **softprops/action-gh-release@v2**: Popular action for creating releases (46k+ stars)

### Permissions

The workflows require:
- `contents: write` - To create releases and tags

### Security

- Uses `GITHUB_TOKEN` automatically provided by GitHub Actions
- No custom secrets required
- All actions pinned to specific versions

## Troubleshooting

### Workflow doesn't run
- Ensure you're on the `main` branch
- Check that the workflow files are in `.github/workflows/`
- Verify repository permissions allow workflow execution

### Version calculation is incorrect
- Check commit message prefixes
- Review git tags with `git tag -l`
- Ensure full git history is available (`fetch-depth: 0`)

### Binary build fails
- Verify Go version compatibility (requires Go 1.21+)
- Check for build errors in workflow logs
- Ensure all Go files compile successfully

### Release creation fails
- Verify `GITHUB_TOKEN` has `contents: write` permission
- Check for existing releases with the same tag
- Review action logs for specific error messages

## Future Enhancements

Potential improvements for the workflow:

1. **Multi-platform builds**: Add support for multiple architectures (amd64, arm, etc.)
2. **Changelog automation**: Use tools like `git-chglog` for more structured changelogs
3. **Release notes templates**: Add customizable templates for release notes
4. **Automated publishing**: Option to automatically publish releases after manual approval
5. **Build verification**: Add tests before creating release
6. **Checksums**: Generate SHA256 checksums for binary artifacts

## Examples

### Example Release Notes

```markdown
## What's Changed

### Commits since v0.1.0
- feat: Add configurable monitoring interval (abc1234)
- fix: Improve error handling in battery reader (def5678)
- docs: Update installation instructions (ghi9012)

## Binary Downloads

- **yukti-arm64-linux.tar.gz**: ARM64 Linux binary (optimized for Raspberry Pi, Android with root, AWS Graviton, etc.)

### Installation

\`\`\`bash
# Extract the archive
tar -xzf yukti-arm64-linux.tar.gz

# Move to system path (optional)
sudo mv yukti /usr/local/bin/yukti

# Run (requires root privileges)
sudo yukti
\`\`\`

⚠️ **Note**: This application requires a ROOTED Android device or root/sudo privileges on Linux systems.
```

## Contributing

When contributing changes that affect the release workflow:

1. Test workflow syntax with `yamllint` or online validators
2. Update this documentation if adding new features
3. Follow the commit message conventions for proper versioning
4. Test workflow changes in a fork before merging

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)
