package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/culdo/bbs-restful-api/middleware"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AdminTestSuite struct {
	BaseTestSuite
}
func (s *AdminTestSuite) SetupSuite() {
	SetupBaseTest(s)
}
func TestAdmin(t *testing.T) {
    suite.Run(t, new(AdminTestSuite))
}

func (s *AdminTestSuite) setupTestRouter(){
	router := gin.Default()
	admin := router.Group("/admin")

	admin.GET("/posts", middleware.DoHidePost(false), FetchPosts)
	admin.GET("/posts/search", SearchAllPost)
	admin.PATCH("/posts/:id", UpdatePost)
	admin.PATCH("/users/:id", UpdateUser)
	s.testRouter = router
}

func fakeUser(id uint) (string, string){
	user := model.User{}
	user.ID = id
	user.Username = "test_username" + fmt.Sprint(user.ID)
	password := "test_hashedpass" + fmt.Sprint(user.ID)
	var err error
	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		panic(err.Error())
	}
	err = model.Save(&user)
	if err != nil {
		panic(err.Error())
	}
	return user.Username, password
}

func (s *AdminTestSuite) adminRunTask(taskType string, subId uint, updateReq map[string]interface{}) (map[string] interface{}, error) {
	
	if taskType=="post" {
		fakePost(subId) 
	}else if taskType=="user" {
		fakeUser(subId)
	}

	postResp := httptest.NewRecorder()

	var jsonStr, _ = json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/admin/"+taskType+"s/"+fmt.Sprint(subId), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	s.testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] interface{}
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	resp["code"] = postResp.Code
	return resp, nil
}

func (s *AdminTestSuite)TestAdminPost() {
	resp, err := s.adminRunTask("post", 101, map[string]interface{}{"hidden":true})
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "Post is updated!", resp["message"])
	assert.Equal(s.T(), "101", resp["pid"])
	resp, err = s.adminRunTask("post", 101, map[string]interface{}{"hidden":false})
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "Post is updated!", resp["message"])
	assert.Equal(s.T(), "101", resp["pid"])
}

func (s *AdminTestSuite)TestAdminUser() {
	resp, err := s.adminRunTask("user", 101, map[string]interface{}{"active":true})
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "User is updated!", resp["message"])
	assert.Equal(s.T(), "101", resp["uid"])
	resp, err = s.adminRunTask("user", 101, map[string]interface{}{"active":false})
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "User is updated!", resp["message"])
	assert.Equal(s.T(), "101", resp["uid"])
}