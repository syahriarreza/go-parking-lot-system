package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syahriarreza/go-parking-lot-system/internal/service"
	"github.com/syahriarreza/go-parking-lot-system/internal/utils"
)

type FloorLayoutRequest struct {
	Floor  int        `json:"floor"`
	Layout [][]string `json:"layout"`
}

type UploadResponse struct {
	Message string `json:"message"`
}

func SetLayoutPerFloor(c echo.Context) error {
	var req FloorLayoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid layout input"})
	}

	if err := service.SaveLayoutForFloor(req.Floor, req.Layout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UploadResponse{Message: "Layout for floor saved successfully"})
}

func UploadExcelLayout(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to open uploaded file"})
	}
	defer src.Close()

	layouts, err := utils.ParseExcelLayout(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	for floor, layout := range layouts {
		if err := service.SaveLayoutForFloor(floor, layout); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, UploadResponse{Message: "Excel layout uploaded successfully"})
}

func GetLayoutMap(c echo.Context) error {
	floorStr := c.QueryParam("floor")
	if floorStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "floor parameter is required"})
	}
	var floor int
	if _, err := fmt.Sscanf(floorStr, "%d", &floor); err != nil || floor < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid floor number"})
	}

	rows, err := service.GetLayoutMapFormatted(floor)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"floor": floor,
		"map":   rows,
	})
}

func ResetSystem(c echo.Context) error {
	service.ResetAllLayouts()
	service.ResetParkingState()
	return c.JSON(http.StatusOK, map[string]string{"message": "System reset successful"})
}
