package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/culdo/bbs-restful-api/auth"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/middleware"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	router := s.testRouter
	admin := router.Group("/admin")
	admin.Use(auth.AuthRequired("admin"))
	{
		admin.GET("/posts", middleware.DoHidePost(false), FetchPosts)
		admin.GET("/posts/search", SearchAllPost)
		admin.POST("/posts/:id/hide", HidePost)
		admin.POST("/posts/:id/unhide", UnhidePost)
		admin.POST("/users/:id/ban", BanUser)
		admin.POST("/users/:id/activate", ActivateUser)
	}
}

func (s *AdminTestSuite) adminRunTask(name string, subId int) (map[string] interface{}, error) {
	postResp := httptest.NewRecorder()
	taskName := strings.Split(name, " ")[0]
	taskType := strings.Split(name, " ")[1]

	req, _ := http.NewRequest("POST", "/admin/"+taskType+"s/"+fmt.Sprint(subId)+"/"+taskName, nil)
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

func (s *AdminTestSuite)TestHidePost() {
	err := s.login("admin", config.AdminPasswd)
	// resp, _, err := login("test_name", "test_pass")
	if err!=nil {
		log.Print(err.Error())
	}
	testPost := &model.Post{}
	testPost.Title = "test_title"
	testPost.Content = "test_content"
	model.Save(testPost)
	pid := 1
	resp, err := s.adminRunTask("hide post", pid)
	if err!=nil {
		log.Print(err.Error())
	}
	assert.Equal(s.T(), http.StatusCreated, resp["code"])
	assert.Equal(s.T(), "Post is hidden!", resp["message"])
	assert.Equal(s.T(), fmt.Sprint(pid), resp["pid"])
}

