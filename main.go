package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anshulpatel25/yukti/battery"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[Battery Manager] ", log.LstdFlags)

	// Initialize battery repository with file system implementation
	capacityPath := "/sys/class/power_supply/battery/capacity"
	chargingEnabledPath := "/sys/class/power_supply/battery/charging_enabled"
	repo := battery.NewFileSystemRepository(capacityPath, chargingEnabledPath)

	// Initialize battery service
	service := battery.NewService(repo, logger)

	// Create ticker for 60-second intervals
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Println("Battery charging manager started")

	// Run immediately on start
	if err := service.ManageCharging(); err != nil {
		logger.Printf("Error managing charging: %v", err)
	}

	// Main loop
	for {
		select {
		case <-ticker.C:
			if err := service.ManageCharging(); err != nil {
				logger.Printf("Error managing charging: %v", err)
			}
		case <-sigChan:
			logger.Println("Shutting down gracefully...")
			return
		}
	}
}
