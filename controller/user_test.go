package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/culdo/bbs-restful-api/migration"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setup() *gin.Engine {
	db := model.Init()
	migration.Migrate(db)
	
	router := gin.Default()
	router.POST("/register", RegisterEndpoint)
	return router
}

func TestRegister(t *testing.T) {
	testRouter := setup()

	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = "test_name"
	userReq.Password = "test_pass"
	var jsonStr, _ = json.Marshal(userReq)

	_, err := model.FindUserByName(userReq.Username)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	log.Print(w.Body.String())
	
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "User created successfully", resp["message"])
	} else {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "User already exists", resp["error"])
	}
}