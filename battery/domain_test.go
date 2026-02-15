package battery

import "testing"

func TestDetermineChargingState(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		currentState ChargingState
		expected     ChargingState
	}{
		// Test cases at lower threshold (40%)
		{
			name:         "capacity at 40% - should enable charging",
			capacity:     40,
			currentState: ChargingDisabled,
			expected:     ChargingEnabled,
		},
		{
			name:         "capacity below 40% - should enable charging",
			capacity:     35,
			currentState: ChargingDisabled,
			expected:     ChargingEnabled,
		},
		{
			name:         "capacity very low - should enable charging",
			capacity:     20,
			currentState: ChargingDisabled,
			expected:     ChargingEnabled,
		},
		// Test cases at upper threshold (70%)
		{
			name:         "capacity at 70% - should disable charging",
			capacity:     70,
			currentState: ChargingEnabled,
			expected:     ChargingDisabled,
		},
		{
			name:         "capacity above 70% - should disable charging",
			capacity:     85,
			currentState: ChargingEnabled,
			expected:     ChargingDisabled,
		},
		{
			name:         "capacity at 100% - should disable charging",
			capacity:     100,
			currentState: ChargingEnabled,
			expected:     ChargingDisabled,
		},
		// Hysteresis test cases (between 40% and 70%)
		{
			name:         "capacity at 50% while charging - maintain charging enabled",
			capacity:     50,
			currentState: ChargingEnabled,
			expected:     ChargingEnabled,
		},
		{
			name:         "capacity at 50% while not charging - maintain charging disabled",
			capacity:     50,
			currentState: ChargingDisabled,
			expected:     ChargingDisabled,
		},
		{
			name:         "capacity at 41% while charging - maintain charging enabled",
			capacity:     41,
			currentState: ChargingEnabled,
			expected:     ChargingEnabled,
		},
		{
			name:         "capacity at 69% while not charging - maintain charging disabled",
			capacity:     69,
			currentState: ChargingDisabled,
			expected:     ChargingDisabled,
		},
		{
			name:         "capacity at 60% while charging - maintain charging enabled",
			capacity:     60,
			currentState: ChargingEnabled,
			expected:     ChargingEnabled,
		},
		{
			name:         "capacity at 60% while not charging - maintain charging disabled",
			capacity:     60,
			currentState: ChargingDisabled,
			expected:     ChargingDisabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineChargingState(tt.capacity, tt.currentState)
			if result != tt.expected {
				t.Errorf("DetermineChargingState(capacity=%d, currentState=%v) = %v, expected %v",
					tt.capacity, tt.currentState, result, tt.expected)
			}
		})
	}
}
