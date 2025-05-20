package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/syahriarreza/go-parking-lot-system/internal/handler"
	"github.com/syahriarreza/go-parking-lot-system/internal/service"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Parking struct {
		Floors  int `mapstructure:"floors"`
		Rows    int `mapstructure:"rows"`
		Columns int `mapstructure:"columns"`
	} `mapstructure:"parking"`
}

var AppConfig Config

func initConfig() {
	viper.SetConfigFile("config/app.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

func main() {
	initConfig()

	e := echo.New()

	// inject spotMap into layout_service
	service.InjectSpotMap(service.GetSpotMap())

	v1 := e.Group("/api/v1")
	v1.POST("/park", handler.Park)
	v1.POST("/unpark", handler.Unpark)
	v1.GET("/available", handler.AvailableSpot)
	v1.GET("/search/:vehicleNumber", handler.SearchVehicle)
	v1.POST("/layout/floor", handler.SetLayoutPerFloor)
	v1.POST("/layout/upload", handler.UploadExcelLayout)
	v1.GET("/layout/map", handler.GetLayoutMap)
	v1.POST("/reset", handler.ResetSystem)

	addr := fmt.Sprintf(":%d", AppConfig.Server.Port)
	log.Printf("Server starting at %s...", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
