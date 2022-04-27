package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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
}

type BaseTester interface {
	setupTestRouter()
}

func setupDB() {
	config.DatabaseUrl = "postgresql:///bbstest"
	db := model.Init()
	migration.Migrate(db)
	clean(db)
}

func SetupBaseTest(bt BaseTester) {
	setupDB()
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
	router := gin.Default()

	router.Use(auth.Session("bbssession"))
	router.GET("/", Index)
	router.POST("/login", auth.Login)
	router.GET("/posts", middleware.DoHidePost(true), FetchPosts)

	router.POST("/posts", CreatePost)
	router.POST("/posts/:id/comments", CreateComment)
	
	s.testRouter = router
}

func clean(db *gorm.DB) {
	stmt := `DO $$
		BEGIN
		DELETE FROM users;
		DELETE FROM comments;
		DELETE FROM posts;
		ALTER SEQUENCE users_id_seq RESTART WITH 1;
		ALTER SEQUENCE comments_id_seq RESTART WITH 1;
		ALTER SEQUENCE posts_id_seq RESTART WITH 1;
		END
        $$;`
	db.Exec(stmt)
}

func (s *BaseTestSuite)login(username string, password string) map[string] interface{} {
	w := httptest.NewRecorder()
	var userReq auth.LoginRequest 
	userReq.Username = username
	userReq.Password = password
	var jsonStr, _ = json.Marshal(userReq)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	s.testRouter.ServeHTTP(w, req)
	log.Print(w.Body.String())
	var loginResp map[string] interface{}
	if err := json.NewDecoder(w.Body).Decode(&loginResp); err!=nil {
		log.Print(err.Error())
	}
	
	loginResp["cookie"] = w.Header().Get("Set-Cookie")
	loginResp["code"] = w.Code
	
	return loginResp
}

func (s *CommonTestSuite)TestLogin() {
	model.Register("test_name", "test_pass")
	loginResp := s.login("test_name", "test_pass")
	assert.Equal(s.T(), http.StatusOK, loginResp["code"])
	assert.GreaterOrEqual(s.T(), len(loginResp["cookie"].(string)), 30)

	loginResp = s.login("user_noexist", "passnoexist")
	assert.Equal(s.T(), http.StatusUnauthorized, loginResp["code"])
	assert.Equal(s.T(), "Login failed", loginResp["error"])
}

func (s *CommonTestSuite)TestAdminLogin() {
	model.CreateAdmin()
	loginResp := s.login("admin", config.AdminPasswd)
	assert.Equal(s.T(), http.StatusOK, loginResp["code"])
	assert.GreaterOrEqual(s.T(), len(loginResp["cookie"].(string)), 30)
}

func (s *CommonTestSuite)createPost(title string, content string) (map[string] interface{}, error) {
	loginResp := s.login(fakeUser(10))
	
	var postReq PostRequest 
	postReq.Title = title
	postReq.Content = content
	
	var jsonStr, _ = json.Marshal(postReq)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", loginResp["cookie"].(string))
	
	postResp := httptest.NewRecorder()
	s.testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	resp["code"] = postResp.Code
	return resp, nil
}

func (s *CommonTestSuite)createComment(pid string, content string) (map[string] interface{}, error) {
	loginResp := s.login(fakeUser(11))

	var commentReq CommentRequest 
	commentReq.Content = content
	
	var jsonStr, _ = json.Marshal(commentReq)
	req, _ := http.NewRequest("POST", "/posts/"+pid+"/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", loginResp["cookie"].(string))
	
	postResp := httptest.NewRecorder()
	s.testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	resp["code"] = postResp.Code
	return resp, nil
}
func fakePost(id uint) {
	post := model.Post{}
	post.ID = id
	post.Title = "test_title" + fmt.Sprint(post.ID)
	post.Content = "test_content" + fmt.Sprint(post.ID)
	err := model.Save(&post)
	if err != nil {
		log.Print(err.Error())
	}
}
func (s *CommonTestSuite)TestCreatePost() {
	title := "test_title"
	content := "test_content"
	resp, err := s.createPost(title, content)
	if err != nil {
		log.Print(err.Error())
	}

	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), title, resp["post"].(map[string]interface{})["title"])
	assert.Equal(s.T(), content, resp["post"].(map[string]interface{})["content"])
}

func (s *CommonTestSuite)TestCreateComment() {
	fakePost(11)
	content := "test_comment_content"
	resp, err := s.createComment("11",content)
	log.Print(resp["post"])
	assert.Equal(s.T(), err, nil)
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "Comment created successfully!", resp["message"])
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
