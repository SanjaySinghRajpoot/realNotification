package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SanjaySinghRajpoot/realNotification/controller"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to Real notification"}`
	r := SetUpRouter()
	r.GET("/", utils.HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNotification(t *testing.T) {

	r := SetUpRouter()

	r.POST("/user/notification", controller.Notification)

	userId := make([]int, 0)

	userId = append(userId, 1)

	notification := models.NotificationPayload{
		Type:        "sms",
		Description: "testing kafka",
	}

	jsonValue, _ := json.Marshal(notification)

	reqFound, _ := http.NewRequest("POST", "/user/notification", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)
}
