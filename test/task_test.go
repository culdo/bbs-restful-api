package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/migration"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/culdo/bbs-restful-api/router"
)
var testRouter *gin.Engine

func clean(db *gorm.DB) {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM posts")
	db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE posts_id_seq RESTART WITH 1")
}

func TestMain(m *testing.M) {
	db := model.Init()
	migration.Migrate(db)
	// clean(db)
	model.CreateAdmin()
	testRouter = router.SetupRouter()
	
	m.Run()
	
	os.Exit(0)
}

func TestRegister(t *testing.T) {

	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = "test_name"
	userReq.Password = "test_pass"
	var jsonStr, _ = json.Marshal(userReq)

	_, err := model.FindUserByName(userReq.Username)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)
	
	var resp map[string]interface{}
	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusCreated, w.Code)
		if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
			log.Print(err.Error())
		}
		assert.Equal(t, "User created successfully", resp["message"])
	} else {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
			log.Print(err.Error())
		}
		assert.Equal(t, "User already exists", resp["error"])
	}
}

func login(username string, password string) (map[string] string, *httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = username
	userReq.Password = password
	var jsonStr, _ = json.Marshal(userReq)

	if _, err := model.FindUserByName(userReq.Username);err!=nil{
		return nil, nil, err
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)
	var resp map[string] string
	if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	return resp, w, nil
}

func TestLogin(t *testing.T) {
	resp, w, err := login("test_name", "test_pass")

	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "record not found", resp["message"])
	} else {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.GreaterOrEqual(t, len(resp["expire"]), 30)
		assert.GreaterOrEqual(t, len(resp["token"]), 140)
	}

}

func TestAdminLogin(t *testing.T) {
	resp, w, err := login("admin", config.AdminPasswd)
	
	if err == gorm.ErrRecordNotFound {
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "record not found", resp["message"])
	} else {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.GreaterOrEqual(t, len(resp["expire"]), 30)
		assert.GreaterOrEqual(t, len(resp["token"]), 140)
	}

}

func createPost(loginResp map[string] string, title string, content string) (map[string] interface{}, *httptest.ResponseRecorder, error) {
	token := loginResp["token"]
	postResp := httptest.NewRecorder()

	var postReq model.PostRequest 
	postReq.Title = title
	postReq.Content = content
	var jsonStr, _ = json.Marshal(postReq)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	return resp, postResp, nil
}
func TestCreatePost(t *testing.T) {
	loginResp, _, err := login("test_name", "test_pass")
	if err != nil {
		log.Print(err.Error())
	}
	title := "test_title"
	content := "test_content"
	resp, postResp, err := createPost(loginResp, title, content)
	if err != nil {
		log.Print(err.Error())
	}

	assert.Equal(t, http.StatusCreated, postResp.Code)
	assert.Equal(t, title, resp["Post"].(map[string]interface{})["title"])
	assert.Equal(t, content, resp["Post"].(map[string]interface{})["content"])
}

func TestHidePost(t *testing.T) {
	resp, _, err := login("admin", config.AdminPasswd)
	// resp, _, err := login("test_name", "test_pass")
	if err!=nil {
		log.Print(err.Error())
	}
	pid := 1
	resp, w, err := hidePost(resp, pid)
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Post is hidden!", resp["message"])
	assert.Equal(t, fmt.Sprint(pid), resp["pid"])
}

func TestUserFetchPosts(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	q := req.URL.Query()
	q.Add("page", "1")
	testRouter.ServeHTTP(w, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}

	posts, err := model.FetchPosts(true, 30, 0)
	if err!=nil {
		log.Print(err.Error())
	}
	if len(posts)>0{
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, len(resp["data"].([]interface{})), len(posts))
	} else {
		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "No Posts found", resp["message"])
	}

}

func TestAdminFetchPosts(t *testing.T) {
	loginResp, _, err := login("admin", config.AdminPasswd)
	if err!=nil {
		log.Print(err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/posts", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp["token"])
	q := req.URL.Query()
	q.Add("page", "1")
	testRouter.ServeHTTP(w, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}

	posts, err := model.FetchPosts(false, 30, 0)
	if err!=nil {
		log.Print(err.Error())
	}
	if len(posts)>0{
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, len(resp["data"].([]interface{})), len(posts))
	} else {
		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "No Posts found", resp["message"])
	}

}
