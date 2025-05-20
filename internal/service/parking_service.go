package service

import (
	"errors"
	"fmt"
	"sync"
)

type VehicleRecord struct {
	VehicleType   string
	VehicleNumber string
	SpotID        string
}

var (
	parkingState = make(map[string]*VehicleRecord) // vehicleNumber -> record
	spotMap      = make(map[string]string)         // spotID -> vehicleNumber
	vehicleLog   = make(map[string]string)         // vehicleNumber -> last spotID
	lock         sync.RWMutex
)

// Park tries to assign a free spot for the given vehicle type
func Park(vehicleType, vehicleNumber string) (string, error) {
	lock.Lock()
	defer lock.Unlock()

	if _, exists := parkingState[vehicleNumber]; exists {
		return "", errors.New("vehicle already parked")
	}

	for floor, layout := range layoutData {
		for r, row := range layout {
			for c, spot := range row {
				if isSpotMatch(spot, vehicleType) {
					spotID := fmt.Sprintf("%d-%03d-%03d", floor, r+1, c+1)
					if _, occupied := spotMap[spotID]; !occupied {
						spotMap[spotID] = vehicleNumber
						parkingState[vehicleNumber] = &VehicleRecord{
							VehicleType:   vehicleType,
							VehicleNumber: vehicleNumber,
							SpotID:        spotID,
						}
						return spotID, nil
					}
				}
			}
		}
	}

	return "", errors.New("no available spot")
}

// Unpark removes a vehicle from a spot
func Unpark(spotID, vehicleNumber string) error {
	lock.Lock()
	defer lock.Unlock()

	record, ok := parkingState[vehicleNumber]
	if !ok || record.SpotID != spotID {
		return errors.New("vehicle not found at given spot")
	}

	delete(parkingState, vehicleNumber)
	delete(spotMap, spotID)
	vehicleLog[vehicleNumber] = spotID
	return nil
}

// AvailableSpot lists all free spots for a vehicle type
func AvailableSpot(vehicleType string) []string {
	lock.RLock()
	defer lock.RUnlock()

	var result []string
	for floor, layout := range layoutData {
		for r, row := range layout {
			for c, spot := range row {
				if isSpotMatch(spot, vehicleType) {
					spotID := fmt.Sprintf("%d-%03d-%03d", floor, r+1, c+1)
					if _, occupied := spotMap[spotID]; !occupied {
						result = append(result, spotID)
					}
				}
			}
		}
	}
	return result
}

// SearchVehicle returns the last spot a vehicle was seen at
func SearchVehicle(vehicleNumber string) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()

	if rec, ok := parkingState[vehicleNumber]; ok {
		return rec.SpotID, true
	}
	if spotID, ok := vehicleLog[vehicleNumber]; ok {
		return spotID, true
	}
	return "", false
}

func isSpotMatch(spot, vehicleType string) bool {
	return len(spot) > 0 && spot[0:1] == vehicleType
}

// GetSpotMap exposes the internal spotMap to other packages (e.g. for layout rendering)
func GetSpotMap() map[string]string {
	return spotMap
}

// ResetParkingState clears all parked vehicles and logs
func ResetParkingState() {
	lock.Lock()
	defer lock.Unlock()

	parkingState = make(map[string]*VehicleRecord)
	spotMap = make(map[string]string)
	vehicleLog = make(map[string]string)
}
