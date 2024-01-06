package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SanjaySinghRajpoot/realNotification/controller"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome Real notification"}`
	r := SetUpRouter()
	r.GET("/", HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNotification(t *testing.T) {

	r := SetUpRouter()

	r.POST("/user/notification", controller.Notification)

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
