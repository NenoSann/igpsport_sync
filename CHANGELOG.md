# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-12-10

### Added
- Initial release
- Authentication with iGPSport account
- Activity list fetching with pagination
- Support for multiple file formats (FIT, GPX, TCX)
- Serial download functionality (`DownloadAllActivities`)
- Concurrent download functionality (`DownloadAllActivitiesWithConcurrency`)
- Callback-based processing for downloaded activities
- Comprehensive documentation and examples
- Unit tests for core functionality
- MIT License

### Features
- `GetActivityList`: Fetch activities with pagination
- `GetACtivityDownloadUrl`: Get download URL for a specific activity
- `DownloadFile`: Download a single file
- `DownloadAllActivities`: Sequential download of all activities
- `DownloadAllActivitiesWithConcurrency`: Parallel download with worker pool
