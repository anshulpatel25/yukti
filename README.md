# Yukti - Battery Charging Manager

A clean, idiomatic Go application for managing battery charging based on capacity thresholds with hysteresis logic to optimize battery health and longevity.

## Features

- Monitors battery capacity every 60 seconds
- Implements hysteresis charging logic to prevent rapid cycling
- Enables charging when battery drops to 40% or below
- Disables charging when battery reaches 70% or above
- Maintains current state between thresholds (40-70%)
- Clean architecture with separated concerns
- Comprehensive error handling
- Graceful shutdown support
- Full unit test coverage

## ⚠️ Important: Android Device Requirements

**WARNING: Yukti requires a ROOTED Android device to function properly.**

This application needs write access to `/sys/class/power_supply/battery/charging_enabled`, which is only available with root privileges. On Android devices, this means:

- Your device **must be rooted** (using Magisk or similar)
- The application must be granted root/superuser permissions
- Standard (non-rooted) Android devices **will not work**

Without root access, the application cannot control the charging behavior and will fail when attempting to manage battery charging.

## Architecture

The application follows clean architecture principles:

```
├── main.go                 # Application entry point
└── battery/
    ├── domain.go          # Domain logic and business rules
    ├── domain_test.go     # Unit tests for domain logic
    ├── repository.go      # Repository interface
    ├── filesystem.go      # File system implementation
    └── service.go         # Application service layer
```

### Layers

- **Domain Layer** (`domain.go`): Contains business logic and rules for determining charging state
- **Repository Layer** (`repository.go`, `filesystem.go`): Abstracts data access operations
- **Service Layer** (`service.go`): Orchestrates business logic and coordinates between layers
- **Main** (`main.go`): Application entry point with dependency injection

## Usage

### Build for ARM64 (Cross-Compilation)

For ARM64 Linux devices (Android with root):

```bash
GOOS=linux GOARCH=arm64 go build -o yukti
```

### Run

```bash
sudo ./yukti
```

Note: Root privileges are required to write to `/sys/class/power_supply/battery/charging_enabled`.

### Test

```bash
go test ./battery
```

### Test with Coverage

```bash
go test -cover ./battery
```

## Configuration

The battery paths are currently hardcoded in `main.go`:
- Capacity: `/sys/class/power_supply/battery/capacity`
- Charging control: `/sys/class/power_supply/battery/charging_enabled`

These can be easily modified or made configurable through environment variables if needed.

## Charging Logic

The application implements **hysteresis** to prevent rapid charge cycling, which improves battery longevity:

### Behavior

| Battery Level | Action | Reason |
|--------------|--------|---------|
| ≤ 40% | **Enable charging** | Prevent battery depletion |
| 41% - 69% | **Maintain current state** | Hysteresis zone - no change |
| ≥ 70% | **Disable charging** | Optimal charge level for battery health |

### Example Charging Cycle

```
100% → Charging disabled (≥70%)
 ↓
 69% → Charging remains disabled (hysteresis)
 ↓
 40% → Charging enabled (≤40%)
 ↓
 41% → Charging remains enabled (hysteresis)
 ↓
 70% → Charging disabled (≥70%)
```

This prevents the battery from constantly switching between charging and not charging when hovering near a single threshold, reducing wear on both the battery and charging circuitry.

## Deployment

> **Note**: Before first time deployment, ensure that the battery state of charge is greater than 80%.

### Termux (Recommended for Android)

Termux should be granted with the `superuser` permissions using Magisk or similar.

**Installation**

- Install sudo package: `pkg install sudo`.
- Install termux services: `pkg install termux-services`.
- Download the latest ARM64 binary from the releases page.
- Create a directory for Yukti: `mkdir -p $HOME/yukti`.
- Place the downloaded `yukti` binary in `$HOME/yukti`.
- Grant execute permissions: `chmod +x $HOME/yukti/yukti`.

**Service**

- Create `/data/data/com.termux/files/usr/var/service/yukti/run` with the following content:

```sh
#!/data/data/com.termux/files/usr/bin/sh

exec sudo $HOME/yukti/yukti > $HOME/yukti/yukti.log 2>&1
```

- Grant execute permissions: `chmod +x /data/data/com.termux/files/usr/var/service/yukti/run`
- Enable the service: `sv-enable yukti`
- Start the service: `sv up yukti`
