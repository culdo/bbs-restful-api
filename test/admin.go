package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

)

func hidePost(loginResp map[string] string, pid int) (map[string] string, *httptest.ResponseRecorder, error) {
	token := loginResp["token"]
	postResp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/admin/posts/"+fmt.Sprint(pid)+"/hide", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] string
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	return resp, postResp, nil
}


func unhidePost(loginResp map[string] string, pid int) (map[string] string, *httptest.ResponseRecorder, error) {
	token := loginResp["token"]
	postResp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/admin/posts/"+fmt.Sprint(pid)+"/hide", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] string
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	return resp, postResp, nil
}


func banUser(loginResp map[string] string, uid int) (map[string] string, *httptest.ResponseRecorder, error) {
	token := loginResp["token"]
	postResp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/admin/users/"+fmt.Sprint(uid)+"/ban", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(postResp, req)
	
	var resp map[string] string
	if err := json.NewDecoder(postResp.Body).Decode(&resp); err!=nil {
		log.Print(err.Error())
	}
	return resp, postResp, nil
}