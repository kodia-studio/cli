# Changelog - Kodia CLI 🐨💻

All notable changes to the Kodia CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-04-23

### Added
- **Plugin Management**: New command `kodia plugin install` to fetch and register official plugins.
- **Service Generation**: Support for generating typed service layers via `kodia make service`.
- **Pre-flight Checks**: Improved system dependency checks before project initialization.

### Fixed
- Fixed directory permission issues on Windows during `kodia new`.
- Corrected templates for Svelte 5 runes compatibility.

## [1.1.0] - 2026-04-10

### Added
- **Docker Support**: Auto-generation of optimized Dockerfiles for both Backend and Frontend.
- **Env Validation**: Automatic validation of `.env.example` during startup.

### Improved
- Faster boilerplate extraction using concurrent processing.
- Enhanced color output for better readability in different terminal themes.

## [1.0.0] - 2026-03-25

### Added
- Initial stable release of Kodia CLI.
- Project scaffolding for Fullstack, Backend-only, and Frontend-only modes.
- Integrated `npm` and `go` dependency management.

---
© 2026 Kodia Studio. "Build like a user, code like a pro."
