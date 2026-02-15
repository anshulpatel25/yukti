package battery

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FileSystemRepository implements Repository using the file system
type FileSystemRepository struct {
	capacityPath        string
	chargingEnabledPath string
}

// NewFileSystemRepository creates a new file system based repository
func NewFileSystemRepository(capacityPath, chargingEnabledPath string) *FileSystemRepository {
	return &FileSystemRepository{
		capacityPath:        capacityPath,
		chargingEnabledPath: chargingEnabledPath,
	}
}

// GetCapacity reads the battery capacity from the file system
func (r *FileSystemRepository) GetCapacity() (int, error) {
	data, err := os.ReadFile(r.capacityPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read capacity file: %w", err)
	}

	capacityStr := strings.TrimSpace(string(data))
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse capacity value '%s': %w", capacityStr, err)
	}

	return capacity, nil
}

// GetChargingEnabled reads the current charging state from the file system
func (r *FileSystemRepository) GetChargingEnabled() (ChargingState, error) {
	data, err := os.ReadFile(r.chargingEnabledPath)
	if err != nil {
		return ChargingDisabled, fmt.Errorf("failed to read charging enabled file: %w", err)
	}

	stateStr := strings.TrimSpace(string(data))
	state, err := strconv.Atoi(stateStr)
	if err != nil {
		return ChargingDisabled, fmt.Errorf("failed to parse charging state value '%s': %w", stateStr, err)
	}

	return ChargingState(state), nil
}

// SetChargingEnabled writes the charging state to the file system
func (r *FileSystemRepository) SetChargingEnabled(state ChargingState) error {
	data := fmt.Sprintf("%d\n", state)
	err := os.WriteFile(r.chargingEnabledPath, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write charging state: %w", err)
	}

	return nil
}
