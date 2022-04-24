package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/culdo/bbs-restful-api/auth"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/middleware"
	"github.com/culdo/bbs-restful-api/migration"
	"github.com/culdo/bbs-restful-api/model"
)

type BaseTestSuite struct {
	suite.Suite
	testRouter  *gin.Engine
	loginResp map[string] interface{}
}

type BaseTester interface {
	setupDB()
	setupBaseRouter()
	setupTestRouter()
}

func (s *BaseTestSuite) setupDB() {
	db := model.Init()
	migration.Migrate(db)
	// clean(db)
	model.CreateAdmin()
}

func (s *BaseTestSuite) setupBaseRouter() {
	router := gin.Default()

	router.Use(auth.Session("bbssession"))
	router.GET("/", Index)
	router.POST("/login", auth.Login)
	router.GET("/logout", auth.Logout)
	s.testRouter = router
}

func SetupBaseTest(bt BaseTester) {
	bt.setupDB()
	bt.setupBaseRouter()
	bt.setupTestRouter()
}

type CommonTestSuite struct {
	BaseTestSuite
}	

func (s *CommonTestSuite) SetupSuite() {
	SetupBaseTest(s)
}	

func TestCommon(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}	

func (s *CommonTestSuite) setupTestRouter() {
	router := s.testRouter
	router.GET("/posts", middleware.DoHidePost(true), FetchPosts)
	router.Use(auth.AuthRequired("user")) 
	{
		router.POST("/posts", CreatePost)
		router.POST("/posts/:id/comments", CreateComment)
	}
}

func clean(db *gorm.DB) {
	stmt := `DO $$
		BEGIN
		DELETE FROM users;
		DELETE FROM posts;
		ALTER SEQUENCE users_id_seq RESTART WITH 1;
		ALTER SEQUENCE posts_id_seq RESTART WITH 1;
		END
        $$;`
	db.Exec(stmt)
}

func (s *BaseTestSuite)login(username string, password string) error {
	w := httptest.NewRecorder()
	var userReq model.UserRequest 
	userReq.Username = username
	userReq.Password = password
	var jsonStr, _ = json.Marshal(userReq)

	if _, err := model.FindUserByName(userReq.Username);err!=nil{
		return err
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	s.testRouter.ServeHTTP(w, req)
	if err := json.NewDecoder(w.Body).Decode(&s.loginResp); err!=nil {
		log.Print(err.Error())
	}
	s.loginResp["cookie"] = w.Header().Get("Set-Cookie")
	s.loginResp["code"] = w.Code
	
	return nil
}

func (s *CommonTestSuite)TestLogin() {
	err := s.login("test_name", "test_pass")

	if err == gorm.ErrRecordNotFound {
		assert.Equal(s.T(), http.StatusUnauthorized, s.loginResp["code"])
		assert.Equal(s.T(), "record not found", s.loginResp["message"])
	} else {
		assert.Equal(s.T(), http.StatusOK, s.loginResp["code"])
		assert.GreaterOrEqual(s.T(), len(s.loginResp["cookie"].(string)), 30)
	}

}

func (s *CommonTestSuite)TestAdminLogin() {
	err := s.login("admin", config.AdminPasswd)
	
	if err == gorm.ErrRecordNotFound {
		assert.Equal(s.T(), http.StatusUnauthorized, s.loginResp["code"])
		assert.Equal(s.T(), "record not found", s.loginResp["message"])
	} else {
		assert.Equal(s.T(), http.StatusOK, s.loginResp["code"])
		assert.GreaterOrEqual(s.T(), len(s.loginResp["cookie"].(string)), 30)
	}

}

func (s *CommonTestSuite)createPost(title string, content string) (map[string] interface{}, error) {
	postResp := httptest.NewRecorder()

	var postReq model.PostRequest 
	postReq.Title = title
	postReq.Content = content
	var jsonStr, _ = json.Marshal(postReq)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", s.loginResp["cookie"].(string))
	s.testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	resp["code"] = postResp.Code
	return resp, nil
}
func (s *CommonTestSuite)TestCreatePost() {
	err := s.login("test_name", "test_pass")
	if err != nil {
		log.Print(err.Error())
	}
	title := "test_title"
	content := "test_content"
	resp, err := s.createPost(title, content)
	if err != nil {
		log.Print(err.Error())
	}

	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), title, resp["Post"].(map[string]interface{})["title"])
	assert.Equal(s.T(), content, resp["Post"].(map[string]interface{})["content"])
}

func (s *CommonTestSuite)TestUserFetchPosts() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	q := req.URL.Query()
	q.Add("page", "1")
	s.testRouter.ServeHTTP(w, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}

	posts, err := model.FetchPosts(true, 30, 0)
	if err!=nil {
		log.Print(err.Error())
	}
	if len(posts)>0{
		assert.Equal(s.T(), 200, w.Code)
		assert.Equal(s.T(), len(resp["data"].([]interface{})), len(posts))
	} else {
		assert.Equal(s.T(), 404, w.Code)
		assert.Equal(s.T(), "No Posts found", resp["message"])
	}

}
