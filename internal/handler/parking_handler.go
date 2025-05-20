package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syahriarreza/go-parking-lot-system/internal/service"
)

type ParkRequest struct {
	VehicleType   string `json:"vehicleType"`
	VehicleNumber string `json:"vehicleNumber"`
}

type UnparkRequest struct {
	SpotID        string `json:"spotId"`
	VehicleNumber string `json:"vehicleNumber"`
}

type SearchResponse struct {
	SpotID string `json:"spotId,omitempty"`
	Found  bool   `json:"found"`
}

type AvailableResponse struct {
	Spots []string `json:"spots"`
}

func Park(c echo.Context) error {
	var req ParkRequest
	if err := c.Bind(&req); err != nil || req.VehicleType == "" || req.VehicleNumber == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	spotID, err := service.Park(req.VehicleType, req.VehicleNumber)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"spotId": spotID})
}

func Unpark(c echo.Context) error {
	var req UnparkRequest
	if err := c.Bind(&req); err != nil || req.SpotID == "" || req.VehicleNumber == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	if err := service.Unpark(req.SpotID, req.VehicleNumber); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "vehicle unparked"})
}

func AvailableSpot(c echo.Context) error {
	vehicleType := c.QueryParam("vehicleType")
	if vehicleType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vehicleType is required"})
	}

	spots := service.AvailableSpot(vehicleType)
	return c.JSON(http.StatusOK, AvailableResponse{Spots: spots})
}

func SearchVehicle(c echo.Context) error {
	vehicleNumber := c.Param("vehicleNumber")
	if vehicleNumber == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vehicleNumber is required"})
	}

	spotID, found := service.SearchVehicle(vehicleNumber)
	return c.JSON(http.StatusOK, SearchResponse{SpotID: spotID, Found: found})
}
