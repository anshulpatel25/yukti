package battery

import (
	"fmt"
	"log"
)

// Service handles battery charging management business logic
type Service struct {
	repo   Repository
	logger *log.Logger
}

// NewService creates a new battery management service
func NewService(repo Repository, logger *log.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// ManageCharging reads battery capacity and adjusts charging state accordingly
func (s *Service) ManageCharging() error {
	// Read current battery capacity
	capacity, err := s.repo.GetCapacity()
	if err != nil {
		return fmt.Errorf("failed to get battery capacity: %w", err)
	}

	// Read current charging state
	currentState, err := s.repo.GetChargingEnabled()
	if err != nil {
		return fmt.Errorf("failed to get current charging state: %w", err)
	}

	// Determine desired charging state based on capacity and current state
	desiredState := DetermineChargingState(capacity, currentState)

	// Only update if state needs to change
	if desiredState != currentState {
		if err := s.repo.SetChargingEnabled(desiredState); err != nil {
			return fmt.Errorf("failed to set charging state: %w", err)
		}
		stateStr := "enabled"
		if desiredState == ChargingDisabled {
			stateStr = "disabled"
		}
		s.logger.Printf("Battery capacity: %d%%, Charging: %s (changed)", capacity, stateStr)
	} else {
		stateStr := "enabled"
		if currentState == ChargingDisabled {
			stateStr = "disabled"
		}
		s.logger.Printf("Battery capacity: %d%%, Charging: %s (unchanged)", capacity, stateStr)
	}

	return nil
}
