package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/culdo/bbs-restful-api/model"
	"github.com/culdo/bbs-restful-api/router"
)


func TestRegister(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	var user model.UserRequest 
	user.Username = "test_name"
	user.Password = "test_pass"
	var jsonStr, _ = json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/vi/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

}

func TestLogin(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

}

func TestCreatePost(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

}

func TestFetchAllPost(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

}
