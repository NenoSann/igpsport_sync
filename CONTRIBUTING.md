# Contributing to igpsport_sync

Thank you for your interest in contributing to igpsport_sync! We welcome contributions from everyone.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps which reproduce the problem**
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed after following the steps**
- **Explain which behavior you expected to see instead and why**
- **Include screenshots and animated GIFs if possible**
- **Include your Go version** (`go version`)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Use a clear and descriptive title**
- **Provide a step-by-step description of the suggested enhancement**
- **Provide specific examples to demonstrate the steps**
- **Describe the current behavior and expected behavior**
- **Explain why this enhancement would be useful**

### Pull Requests

- Fill in the required template
- Follow the Go code style guidelines
- End all files with a newline
- Include appropriate test cases
- Update documentation as needed
- Ensure all tests pass: `go test ./...`

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/igpsport_sync.git`
3. Create a branch for your feature: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `go test -v ./...`
6. Run linter: `go fmt ./...` and `go vet ./...`
7. Commit your changes: `git commit -am 'Add some feature'`
8. Push to the branch: `git push origin feature/your-feature-name`
9. Create a Pull Request

## Code Style

We follow the standard Go code style:

- Use `gofmt` to format code
- Use `govet` to check for common errors
- Use meaningful variable and function names
- Write comments for exported functions
- Keep functions small and focused

## Testing

Please write tests for your changes:

```go
func TestYourFeature(t *testing.T) {
    // Your test code
}
```

Run tests with:
```bash
go test -v ./...
go test -race ./...  # Race condition detection
```

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## License

By contributing to igpsport_sync, you agree that your contributions will be licensed under its MIT License.

## Questions?

Feel free to open an issue or start a discussion if you have any questions!
