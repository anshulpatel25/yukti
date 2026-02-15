package battery

// Repository defines the interface for battery operations
type Repository interface {
	// GetCapacity reads the current battery capacity percentage
	GetCapacity() (int, error)
	// GetChargingEnabled reads the current charging state
	GetChargingEnabled() (ChargingState, error)
	// SetChargingEnabled sets the charging state
	SetChargingEnabled(state ChargingState) error
}
