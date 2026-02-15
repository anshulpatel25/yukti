# AGENTS.md - AI Agent Guide for Yukti

This document provides essential context and guidelines for AI agents working with the Yukti Battery Management System repository.

## Repository Overview

**Yukti** is a clean, idiomatic Go application for managing battery charging based on capacity thresholds with hysteresis logic to optimize battery health and longevity. It runs on Linux systems with battery charging control interfaces (e.g., `/sys/class/power_supply/battery/`).

### Key Features
- Battery capacity monitoring (60-second intervals)
- Hysteresis charging logic (40%-70% thresholds)
- Clean architecture with separated concerns
- Comprehensive unit test coverage
- Graceful shutdown support

## Project Structure

```
.
├── main.go                 # Application entry point with dependency injection
├── go.mod                  # Go module definition
├── README.md               # User-facing documentation
├── LICENSE                 # Apache 2.0 license
└── battery/                # Core battery management package
    ├── domain.go          # Business logic and charging rules
    ├── domain_test.go     # Unit tests for domain logic
    ├── repository.go      # Repository interface
    ├── filesystem.go      # File system implementation
    └── service.go         # Application service layer
```

## Architecture

The application follows **Clean Architecture** principles with clear layer separation:

### 1. Domain Layer (`battery/domain.go`)
- **Purpose**: Core business logic for battery charging decisions
- **Key Components**:
  - `ChargingState` enum (Enabled=1, Disabled=0)
  - `LowerThreshold` = 40% (enable charging at or below)
  - `UpperThreshold` = 70% (disable charging at or above)
  - `DetermineChargingState()` function implementing hysteresis
- **No Dependencies**: Pure business logic with no external dependencies

### 2. Repository Layer (`battery/repository.go`, `battery/filesystem.go`)
- **Purpose**: Abstract data access to battery system files
- **Interface**: `Repository` defines contract for battery operations
- **Implementation**: `FileSystemRepository` reads/writes to sysfs paths
- **Operations**:
  - `GetCapacity()`: Read battery capacity percentage
  - `GetChargingEnabled()`: Read current charging state
  - `SetChargingEnabled()`: Write charging state

### 3. Service Layer (`battery/service.go`)
- **Purpose**: Orchestrate business logic and coordinate between layers
- **Key Method**: `ManageCharging()` - main workflow for charging management
- **Responsibilities**:
  - Read battery capacity
  - Get current charging state
  - Determine desired state using domain logic
  - Update charging state if needed
  - Log state changes

### 4. Main (`main.go`)
- **Purpose**: Application bootstrapping and runtime management
- **Responsibilities**:
  - Initialize logger and dependencies
  - Set up 60-second ticker
  - Handle graceful shutdown (SIGINT, SIGTERM)
  - Run main event loop

## Domain Knowledge

### Hysteresis Logic
The charging logic implements **hysteresis** to prevent rapid charge cycling:

| Battery Level | Action | Reason |
|--------------|--------|---------|
| ≤ 40% | Enable charging | Prevent battery depletion |
| 41% - 69% | Maintain current state | Hysteresis zone |
| ≥ 70% | Disable charging | Optimal charge level for health |

**Why Hysteresis?**
- Prevents constant switching near a single threshold
- Reduces wear on battery and charging circuitry
- Improves battery longevity

### System Integration
- **Platform**: Linux systems with sysfs battery interface
- **Privileges**: Requires root/sudo to write to charging control
- **Paths** (configurable in `main.go`):
  - Capacity: `/sys/class/power_supply/battery/capacity`
  - Charging control: `/sys/class/power_supply/battery/charging_enabled`

## Build, Test, and Deployment

### Building
```bash
# Native build
go build -o yukti

# ARM64 cross-compilation (Raspberry Pi, Android, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o yukti-arm64

# Optimized ARM64 build
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o yukti-arm64
```

### Testing
```bash
# Run tests
go test ./battery

# With coverage
go test -cover ./battery

# Verbose output
go test -v ./battery
```

### Running
```bash
# Requires root privileges
sudo ./yukti
```

### Deployment
The application is typically deployed as a **systemd service** for automatic startup and restart. See README.md for systemd configuration.

## Coding Conventions

### Go Style
- Follow idiomatic Go conventions
- Use standard Go formatting (`go fmt`)
- Go version: 1.21 or higher
- No external dependencies (standard library only)

### Architecture Patterns
- **Clean Architecture**: Maintain layer separation
- **Dependency Injection**: Pass dependencies via constructors
- **Interface-based Design**: Use interfaces for abstractions (Repository)
- **Error Handling**: Wrap errors with context using `fmt.Errorf` with `%w`

### Naming Conventions
- **Constants**: PascalCase (e.g., `LowerThreshold`, `ChargingEnabled`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Packages**: lowercase single word (e.g., `battery`)

### Testing Practices
- **Test Files**: `*_test.go` in same package
- **Test Functions**: `Test*` prefix
- **Coverage**: Aim for comprehensive coverage
- **Table-Driven Tests**: Use for multiple scenarios (see `domain_test.go`)

## Common Tasks for Agents

### Adding New Features
1. **Domain Logic First**: Start with `battery/domain.go` for business rules
2. **Write Tests**: Add tests in `battery/domain_test.go` first (TDD)
3. **Update Service**: Modify `battery/service.go` if orchestration changes needed
4. **Update Repository**: Add new interface methods if data access changes
5. **Update Main**: Modify `main.go` only for configuration or bootstrapping changes

### Modifying Thresholds
- **Location**: `battery/domain.go` constants
- **Considerations**: Update tests and documentation
- **Testing**: Verify hysteresis behavior with new values

### Changing Monitoring Interval
- **Location**: `main.go` ticker initialization
- **Current**: 60 seconds
- **Considerations**: More frequent = more CPU, less frequent = slower response

### Adding Configuration
- Current paths are hardcoded in `main.go`
- Consider environment variables or config file
- Maintain backward compatibility
- Update documentation

### Extending Repository
1. Add method to `Repository` interface
2. Implement in `FileSystemRepository`
3. Update `Service` to use new method
4. Add tests for new functionality

## File Modification Guidelines

### High-Risk Changes
- `battery/domain.go`: Core business logic, verify hysteresis behavior
- `battery/service.go`: Main workflow, ensure error handling preserved
- Test files: Maintain existing coverage

### Medium-Risk Changes
- `main.go`: Application bootstrap, verify graceful shutdown
- `battery/filesystem.go`: File I/O, handle errors properly

### Low-Risk Changes
- Documentation files (README.md, this file)
- Comments and logging messages

## Dependencies

- **Standard Library Only**: No external dependencies
- **Go Version**: 1.21 or higher (specified in `go.mod`)
- **Platform**: Linux with sysfs battery interface

## Error Handling

- Always wrap errors with context using `fmt.Errorf` with `%w` verb
- Log errors but don't crash the application (graceful degradation)
- Repository layer returns errors for I/O failures
- Service layer logs errors and continues monitoring

## Security Considerations

- **Privilege Escalation**: Application requires root to write to sysfs
- **Path Injection**: Paths are hardcoded, not user-supplied
- **Input Validation**: Battery capacity is read from system files (trusted)
- **Resource Safety**: Graceful shutdown prevents resource leaks

## Testing Strategy

### Unit Tests (`battery/domain_test.go`)
- Test hysteresis logic with all threshold combinations
- Use table-driven tests for comprehensive coverage
- Test edge cases (0%, 100%, boundary values)

### Integration Testing
- Manual testing on target hardware recommended
- Verify sysfs paths exist on target platform
- Test graceful shutdown behavior

## Performance Considerations

- **Monitoring Interval**: 60 seconds balances responsiveness and efficiency
- **I/O Operations**: Minimal (two reads, conditional write per cycle)
- **Memory Usage**: Low (~5-10 MB typical)
- **CPU Usage**: Negligible (mostly idle waiting for ticker)

## Troubleshooting for Agents

### Common Issues

**Permission Errors**
- Symptom: "permission denied" writing to charging_enabled
- Solution: Application must run with root privileges (`sudo`)

**File Not Found**
- Symptom: Error reading capacity or charging_enabled
- Solution: Verify sysfs paths exist on target system
- Note: Paths may vary by device (e.g., `/sys/class/power_supply/BAT0/`)

**State Not Changing**
- Symptom: Charging state doesn't match expected behavior
- Solution: Check if device supports charging control via sysfs
- Note: Some devices don't expose charging_enabled interface

## Best Practices for Agents

1. **Read Documentation First**: Review README.md and code before making changes
2. **Test Before Committing**: Run `go test ./battery` after code changes
3. **Maintain Test Coverage**: Add tests for new functionality
4. **Preserve Architecture**: Keep layer separation intact
5. **Follow Go Idioms**: Use standard Go conventions and formatting
6. **Error Context**: Always wrap errors with descriptive context
7. **Minimal Changes**: Make smallest changes necessary to accomplish goals
8. **Update Documentation**: Keep README.md and this file in sync with code

## Quick Reference

### Key Functions
- `DetermineChargingState(capacity, currentState)`: Core hysteresis logic
- `service.ManageCharging()`: Main workflow orchestration
- `repo.GetCapacity()`: Read battery percentage
- `repo.SetChargingEnabled(state)`: Control charging

### Key Constants
- `LowerThreshold = 40`: Enable charging threshold
- `UpperThreshold = 70`: Disable charging threshold
- `ChargingEnabled = 1`: Charging on state
- `ChargingDisabled = 0`: Charging off state

### Key Paths
- `/sys/class/power_supply/battery/capacity`: Battery percentage
- `/sys/class/power_supply/battery/charging_enabled`: Charging control

## Version Information

- **Go Version**: 1.21+
- **Module**: github.com/anshulpatel25/yukti
- **License**: Apache 2.0

## Additional Resources

- **README.md**: User-facing documentation and usage examples
- **Domain Tests**: `battery/domain_test.go` shows expected behavior
- **Go Documentation**: Run `go doc ./battery` for package documentation

---

*This document is maintained for AI agents. When making code changes, update this file if architecture, conventions, or important patterns change.*
