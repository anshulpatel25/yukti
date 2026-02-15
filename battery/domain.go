package battery

// ChargingState represents whether charging should be enabled or disabled
type ChargingState int

const (
	// ChargingDisabled represents charging disabled state (0)
	ChargingDisabled ChargingState = 0
	// ChargingEnabled represents charging enabled state (1)
	ChargingEnabled ChargingState = 1
)

const (
	// LowerThreshold is the minimum capacity to disable charging
	LowerThreshold = 40
	// UpperThreshold is the maximum capacity to re-enable charging
	UpperThreshold = 70
)

// DetermineChargingState calculates the desired charging state based on battery capacity
// and current state, implementing hysteresis to prevent rapid switching.
// - Enable charging when capacity drops to 40% or below
// - Disable charging when capacity reaches 70% or above
// - Maintain current state when capacity is between 40% and 70% (hysteresis)
func DetermineChargingState(capacity int, currentState ChargingState) ChargingState {
	// If capacity drops to 40% or below, enable charging
	if capacity <= LowerThreshold {
		return ChargingEnabled
	}

	// If capacity reaches 70% or above, disable charging
	if capacity >= UpperThreshold {
		return ChargingDisabled
	}

	// Between 40% and 70%, maintain current state (hysteresis)
	return currentState
}
