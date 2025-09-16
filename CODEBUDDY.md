# HasteGUI / HasteCLI Codebase Guide

## Project Overview
This is a Go-based duplicate file scanner with both CLI and GUI interfaces. The project uses Fyne for the GUI and implements a shared core scanning engine. It's designed to detect duplicate files using various scanning modes (basic, video, text, image) with configurable processing strategies.

## Development Commands

### Building
```bash
# Build CLI version
go build -o hastecli.exe .\cmd\hastecli

# Build GUI version (requires C toolchain on Windows)
go build -o hastegui.exe .\cmd\hastegui

# Build GUI with software rendering (no C toolchain needed)
go build -tags nogl -o hastegui.exe .\cmd\hastegui
```

### Running
```bash
# CLI example
./hastecli.exe --paths "D:\\,E:\\docs" --exclude "*.tmp;node_modules" --mode basic --concurrency 4

# GUI
./hastegui.exe
```

### Testing
```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/core
```

## Architecture

### Core Components

**Dual Interface Design**: The project implements both CLI (`cmd/hastecli`) and GUI (`cmd/hastegui`) entry points that share the same core scanning engine.

**Scanner Engine Interface** (`internal/core/engine.go`): Defines `ScannerEngine` interface that all scanning implementations must follow. Current implementation is `SimpleScanner` with plans for more optimized versions.

**Configuration System** (`internal/core/model.go`): 
- `ScanConfig` struct centralizes all scan parameters
- Supports multiple modes: basic, video, text, image
- Configurable filters (size, patterns, hash algorithms)
- Progress callback system for UI updates

**Processing Pipeline**:
1. Path walking with exclusion patterns
2. File hashing (SHA1/SHA256/MD5)
3. Duplicate grouping by hash
4. Strategy-based processing (delete, move, rename)

### GUI Architecture (`internal/gui/`)

**State Management** (`state.go`): Centralized application state with mutex protection for concurrent access.

**Page System** (`pages*.go`): Modular page design matching the requirements:
- Scan configuration page
- Strategy configuration page  
- Execution monitoring page
- Settings page

**Internationalization** (`i18n.go`): Built-in support for Chinese and English localization.

### Key Features

**Multi-Mode Scanning**: Different algorithms for different file types (basic hash, video similarity, text comparison, image analysis).

**Strategy System** (`policy.go`, `executor.go`): Configurable processing rules for handling duplicates with preview and undo capabilities.

**Media Processing** (`media*.go`): Specialized handling for video files with FFmpeg integration and caching system.

**Preset Management** (`presets.go`): Save/load common configuration templates.

## Development Notes

**Windows Focus**: GUI build requires Visual Studio Build Tools or can use software rendering with `-tags nogl`.

**FFmpeg Integration**: Video processing requires FFmpeg binary, path configurable via `HASTE_FFMPEG_PATH` environment variable.

**Concurrency**: Scanner supports configurable concurrency levels for performance tuning.

**Progress Reporting**: All scanning operations support progress callbacks for real-time UI updates.

**Error Handling**: Graceful handling of file permissions, locks, and I/O errors throughout the pipeline.