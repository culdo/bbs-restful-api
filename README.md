# BBS RESTful API
Try it on [bbs-restful-api.herokuapp.com](https://bbs-restful-api.herokuapp.com/)
# 功能
* `/posts`: 瀏覽目前的留言(不包含隱藏留言，預設30筆)
* `/register`: 註冊帳號
* `/login`: 登入帳號
## 使用者
* `POST /posts`: 創建留言
* `POST /posts/:id/comments`: 回覆特定留言
## 管理員
* `/admin/posts`: 瀏覽目前的留言(包含隱藏留言，預設30筆)
* `/admin/posts/search`: 搜尋留言內容
* `POST /admin/posts/:id/hide`: 隱藏特定留言
* `POST /admin/posts/:id/unhide`: 顯示特定留言
* `POST /admin/users/:id/ban`:停權特定使用者
* `POST /admin/users/:id/activate`:解封特定使用者
# Developing
```bash
export DATABASE_URL="your_postgredb_url"
export SESSION_KEY="session_store_secret_you_want"
export ADMIN_PASSWD="bbs_admin_password_you_want"
go run main.go
```
# Testing
`go test -count=1 ./...`
# Testing with GitHub Action
Check out `.github/workflows/main.yaml` to see how it works
# Delopy on Heroku
Set `DATABASE_URL`, `SESSION_KEY` and `ADMIN_PASSWD` as your app's `Config Vars`
# cURL examples
## 使用者
### 取得所有留言(不包含隱藏留言)
`curl -i -X GET -H 'Content-Type: application/json' 127.0.0.1:8080/posts`
### 註冊
`curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' 127.0.0.1:8080/register`
### 登入
`curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' 127.0.0.1:8080/login`
### 發布留言
`curl -i -X POST -d '{"ID":1,"title":"test_title2","content":"test_content2"}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/posts`
### 回覆留言(id = 2)
`curl -i -X POST -d '{"content":"test_comment"}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/post/2/comments`

## 管理員
### 登入
`curl -i -X POST -d '{"username":"admin","password":"your_admin_pass"}' -H 'Content-Type: application/json' 127.0.0.1:8080/login`
### 停權使用者(id = 1)
`curl -i -X POST -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/user/1/ban`
### 解封使用者(id = 1)
`curl -i -X POST -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/user/1/activate`
### 隱藏留言(id = 2)
`curl -i -X POST -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2/hide`
### 顯示留言(id = 2)
`curl -i -X POST -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2/unhide`
### 取得所有留言(包含隱藏留言)
`curl -i -X GET -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts`
### 搜尋留言內容，回傳留言
`curl -i -X GET -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/search?keyword="123"`

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