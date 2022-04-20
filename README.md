# BBS RESTful API
Try it on [bbs-restful-api.herokuapp.com](https://bbs-restful-api.herokuapp.com/)
# 功能
## 使用者
1. 可以瀏覽目前的留言(不包含隱藏留言)
2. 可以註冊帳號
3. 可以登入
4. 登入後可以留言
5. 登入後回覆特定留言
## 管理員
1. 可以瀏覽目前的留言(包含隱藏留言)
2. 可以搜尋留言內容
3. 可以隱藏/顯示留言
4. 可以將使用者停權/解封
# Develop
```bash
export DATABASE_URL="your_postgredb_url"
export JWT_KEY="JWT_secret_key_you_want"
export ADMIN_PASSWD="bbs_admin_password_you_want"
go run main.go
```
# Test
`go test ./test`
# Delopy on Heroku
Set `DATABASE_URL`, `JWT_KEY` and `ADMIN_PASSWD` as your app's `Config Vars`
# API
## 使用者
### 取得所有留言(不包含隱藏留言)
`curl -i -X GET -H 'Content-Type: application/json' 127.0.0.1:8080/posts`
### 註冊
`curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' 127.0.0.1:8080/register`
### 登入
`curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' 127.0.0.1:8080/login`
### 發布留言
`curl -i -X POST -d '{"ID":1,"title":"test_title2","content":"test_content2"}' -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/posts`
### 回覆留言(id = 2)
`curl -i -X POST -d '{"content":"test_comment"}' -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/post/2/comments`

## 管理員
### 登入
`curl -i -X POST -d '{"username":"admin","password":"your_admin_pass"}' -H 'Content-Type: application/json' 127.0.0.1:8080/login`
### 停權使用者(id = 1)
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/user/1/ban`
### 解封使用者(id = 1)
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/user/1/activate`
### 隱藏留言(id = 2)
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2/hide`
### 顯示留言(id = 2)
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2/unhide`
### 取得所有留言(包含隱藏留言)
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts`
### 搜尋留言內容，回傳留言
`curl -i -X GET -H "Authorization: Bearer $token" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/search?keyword="123"`

# CI, CD
使用Github Action部屬於Heroku
# To-do-list
- [x] Use bcrypt on password accessing
- [x] Deploy on cloud service
- [x] Testing

# Dependency
* gin
* gorm
* gin-jwt