package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/syahriarreza/go-parking-lot-system/internal/handler"
)

func setupServer() *echo.Echo {
	e := echo.New()
	v1 := e.Group("/api/v1")
	v1.POST("/park", handler.Park)
	v1.POST("/unpark", handler.Unpark)
	v1.GET("/available", handler.AvailableSpot)
	v1.GET("/search/:vehicleNumber", handler.SearchVehicle)
	v1.POST("/layout/floor", handler.SetLayoutPerFloor)
	v1.POST("/layout/upload", handler.UploadExcelLayout)
	v1.GET("/layout/map", handler.GetLayoutMap)
	v1.POST("/reset", handler.ResetSystem)
	return e
}

func TestFullSystemFlow(t *testing.T) {
	e := setupServer()

	t.Run("Reset system", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reset", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Contains(t, resp["message"], "success")
	})

	t.Run("Upload layout via Excel", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		folderPath := "testdata"
		filePath := filepath.Join(folderPath, "layout.xlsx")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			assert.FailNow(t, fmt.Sprintf("layout.xlsx not found in: %s", filepath.Join(wd, folderPath)))
		}

		part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
		file, _ := os.Open(filePath)
		defer file.Close()
		io.Copy(part, file)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/layout/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "uploaded successfully")
	})

	var spotID string

	t.Run("Park vehicle", func(t *testing.T) {
		body := map[string]string{"vehicleType": "M", "vehicleNumber": "M-001"}
		data, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/park", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		spotID = resp["spotId"]
		if assert.NotEmpty(t, spotID) == false {
			fmt.Println(">>> Response:", resp)
		}
	})

	t.Run("Check available spots", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/available?vehicleType=M", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		type availableStruct struct {
			Spots []string `json:"spots"`
		}
		var resp availableStruct
		json.Unmarshal(rec.Body.Bytes(), &resp)

		if assert.Greater(t, len(resp.Spots), 0) == false {
			fmt.Println(">>> Response:", resp)
		}
	})

	t.Run("Get layout map", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/layout/map?floor=1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["floor"])
		assert.NotEmpty(t, resp["map"])
	})

	t.Run("Unpark vehicle", func(t *testing.T) {
		body := map[string]string{"spotId": spotID, "vehicleNumber": "M-001"}
		data, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/unpark", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
