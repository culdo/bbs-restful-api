package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/culdo/bbs-restful-api/migration"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/culdo/bbs-restful-api/router"
)

func TestMain(m *testing.M) {
	db := model.Init()
	migration.Migrate(db)
	model.CreateAdmin()
	
	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = "test_name"
	userReq.Password = "test_pass"
	var jsonStr, _ = json.Marshal(userReq)

	_, err := model.FindUserByName(userReq.Username)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var m map[string]interface{}
	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusCreated, w.Code)
		if err := json.NewDecoder(w.Body).Decode(&m); err!=nil {
			log.Print(err.Error())
		}
		assert.Equal(t, "User created successfully", m["message"])
	} else {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		if err := json.NewDecoder(w.Body).Decode(&m); err!=nil {
			log.Print(err.Error())
		}
		assert.Equal(t, "User already exists", m["error"])
	}
}

func login(router *gin.Engine, username string, password string) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = username
	userReq.Password = password
	var jsonStr, _ = json.Marshal(userReq)

	if _, err := model.FindUserByName(userReq.Username);err!=nil{
		return nil, err
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w, nil
}

func TestLogin(t *testing.T) {
	router := router.SetupRouter()
	w, err := login(router, "test_name", "test_pass")
	
	
	var m map[string] string
	if err := json.NewDecoder(w.Body).Decode(&m); err!=nil {
		log.Print(err.Error())
	}
	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "record not found", m["message"])
	} else {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.GreaterOrEqual(t, len(m["expire"]), 30)
		assert.Equal(t, 145, len(m["token"]))
	}

}

func TestCreatePost(t *testing.T) {
	router := router.SetupRouter()
	loginResp, err := login(router, "test_name", "test_pass")
	if err != nil {
		log.Print(err.Error())
	}
	var m map[string] interface{}
	if err := json.NewDecoder(loginResp.Body).Decode(&m); err!=nil {
		log.Print(err.Error())
	}
	token := m["token"].(string)
	postResp := httptest.NewRecorder()

	var postReq model.PostRequest 
	postReq.Title = "test_title"
	postReq.Content = "test_content"
	var jsonStr, _ = json.Marshal(postReq)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(postResp, req)
	
	if err := json.NewDecoder(postResp.Body).Decode(&m); err!=nil {
		log.Print(err.Error())
	}

	assert.Equal(t, http.StatusCreated, postResp.Code)
	assert.Equal(t, postReq.Title, m["Post"].(map[string]interface{})["title"])
	assert.Equal(t, postReq.Content, m["Post"].(map[string]interface{})["content"])

}

func TestFetchPosts(t *testing.T) {
	router := router.SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	q := req.URL.Query()
	q.Add("page", "1")
	router.ServeHTTP(w, req)
	
	var m map[string] interface{}
	if err := json.NewDecoder(w.Body).Decode(&m); err!=nil {
		log.Print(err.Error())
	}

	posts, err := model.FetchPosts(true, 30, 0)
	if err!=nil {
		log.Print(err.Error())
	}
	if len(posts)>0{
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, len(m["data"].([]interface{})), len(posts))
	} else {
		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "No Posts found", m["message"])
	}

}
