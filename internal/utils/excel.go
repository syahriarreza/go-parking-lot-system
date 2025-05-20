package utils

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

func ParseExcelLayout(file io.Reader) (map[int][][]string, error) {
	result := make(map[int][][]string)

	xl, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel: %w", err)
	}
	defer xl.Close()

	sheets := xl.GetSheetList()
	for _, sheet := range sheets {
		rows, err := xl.GetRows(sheet)
		if err != nil {
			return nil, fmt.Errorf("failed to read sheet %s: %w", sheet, err)
		}

		floorNum, err := parseFloorFromSheetName(sheet)
		if err != nil {
			return nil, fmt.Errorf("invalid sheet name %s: %w", sheet, err)
		}

		result[floorNum] = rows
	}

	return result, nil
}

func parseFloorFromSheetName(sheet string) (int, error) {
	var floor int
	_, err := fmt.Sscanf(sheet, "Floor%d", &floor)
	if err != nil {
		return 0, fmt.Errorf("sheet name must be like 'Floor1', 'Floor2', etc")
	}
	return floor, nil
}
