# Publishing to pkg.go.dev

This document explains how to publish igpsport_sync to pkg.go.dev.

## Prerequisites

- A GitHub account
- Go installed locally
- The repository pushed to GitHub

## Steps to Publish

### 1. Create a Release on GitHub

```bash
# Create and push a git tag
git tag v0.1.0
git push origin v0.1.0
```

Then go to GitHub and create a Release from the tag.

### 2. Automatic Indexing

Once you push a tag that follows semantic versioning (v0.1.0, v1.0.0, etc.):

1. pkg.go.dev automatically indexes your package within a few minutes
2. Visit `https://pkg.go.dev/github.com/NenoSann/igpsport_sync` to verify

### 3. Verify Package Documentation

Check that:
- README.md is properly rendered
- All exported functions and types are documented
- Code examples are clear and correct
- Package badge appears in README

### 4. Add Badge to README (Optional)

```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/NenoSann/igpsport_sync.svg)](https://pkg.go.dev/github.com/NenoSann/igpsport_sync)
```

## Requirements for pkg.go.dev

✅ Repository on GitHub
✅ Semantic versioning (v0.1.0, v1.0.0)
✅ LICENSE file (MIT)
✅ Proper module path in go.mod
✅ README.md with package description
✅ Documented exported types and functions
✅ Valid Go code that compiles

## Semantic Versioning

- **v0.1.0**: Initial release / API not stable
- **v0.2.0**: New features / minor changes
- **v1.0.0**: Stable API / suitable for production use
- **v1.0.1**: Bug fixes

## After Publishing

### Update Version in go.mod

For users to get the latest version:

```bash
go get -u github.com/NenoSann/igpsport_sync
```

### Monitor Package Health

Check pkg.go.dev for:
- Documentation quality score
- License information
- Dependency details
- Version history

### Troubleshooting

If your package doesn't appear on pkg.go.dev:

1. Check that go.mod is correctly formatted
2. Ensure the repository is public on GitHub
3. Wait a few minutes for automatic indexing
4. Check the status at: https://pkg.go.dev/github.com/NenoSann/igpsport_sync

If still not visible, visit https://pkg.go.dev/search and manually request indexing.

## Version Upgrade Process

To publish a new version:

```bash
# Update version in code if needed
# Update CHANGELOG.md
# Commit changes
git add .
git commit -m "Release v0.2.0"

# Create and push tag
git tag v0.2.0
git push origin v0.2.0

# pkg.go.dev will automatically index the new version
```

## More Information

- [Go Modules Documentation](https://golang.org/doc/modules)
- [pkg.go.dev Guide](https://pkg.go.dev/about)
- [Semantic Versioning](https://semver.org/)
