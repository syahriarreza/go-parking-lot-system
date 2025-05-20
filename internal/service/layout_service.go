package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/syahriarreza/go-parking-lot-system/internal/utils"
)

var (
	layoutData = make(map[int][][]string)
	layoutLock = sync.RWMutex{}
)

func SaveLayoutForFloor(floor int, layout [][]string) error {
	if floor < 1 {
		return errors.New("floor number must be >= 1")
	}

	layoutLock.Lock()
	defer layoutLock.Unlock()

	if len(layoutData) > 0 {
		return errors.New("layout data already exists, please reset data first")
	}

	layoutData[floor] = layout
	return nil
}

func GetLayoutForFloor(floor int) ([][]string, error) {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	if layout, ok := layoutData[floor]; ok {
		return layout, nil
	}
	return nil, fmt.Errorf("layout for floor %d not found", floor)
}

func GetAllFloors() []int {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	var floors []int
	for floor := range layoutData {
		floors = append(floors, floor)
	}
	return floors
}

func InjectSpotMap(sm map[string]string) {
	spotMap = sm
}

func GetLayoutMapFormatted(floor int) ([]string, error) {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	layout, ok := layoutData[floor]
	if !ok {
		return nil, fmt.Errorf("layout for floor %d not found", floor)
	}

	var result []string
	for r, row := range layout {
		var line []string
		for c, spot := range row {
			spotID := fmt.Sprintf("%d-%03d-%03d", floor, r+1, c+1)
			occupied := false
			if spotMap != nil {
				if _, exists := spotMap[spotID]; exists {
					occupied = true
				}
			}
			cell := utils.FormatSpot(spot, occupied)
			line = append(line, cell)
		}
		result = append(result, strings.Join(line, "  "))
	}
	return result, nil
}

func ResetAllLayouts() {
	layoutLock.Lock()
	defer layoutLock.Unlock()
	layoutData = make(map[int][][]string)
}
