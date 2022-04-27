# BBS RESTful API
Try it on [bbs-restful-api.herokuapp.com](https://bbs-restful-api.herokuapp.com/)
# Dependencies
* gin
* gorm
* gin-sessions
# 功能
* `GET /posts`: 瀏覽目前的留言(不包含隱藏留言，預設30筆)
* `GET /register`: 註冊帳號
* `GET /login`: 登入帳號
## 使用者
* `POST /posts`: 創建留言
* `POST /posts/:id/comments`: 回覆特定留言
## 管理員
* `GET /admin/posts`: 瀏覽目前的留言(包含隱藏留言，預設30筆)
* `GET /admin/posts/search`: 搜尋留言內容
* `PATCH /admin/posts/:id`: 隱藏/顯示特定留言
* `PATCH /admin/users/:id`:停權/解封特定使用者
# Developing
```bash
# You need to create a db user if you have't one
sudo -u postgres createuser -s <username>

export DATABASE_URL="your_postgredb_url"
export SESSION_KEY="session_store_secret_you_want"
export ADMIN_PASSWD="bbs_admin_password_you_want"
go run main.go
```
# Testing
```bash
# You need to create a test db if you have't one
sudo -u postgres createdb bbstest
go test -count=1 ./...
```
# Testing with GitHub Action
Check out `.github/workflows/main.yaml` to see how it works
# Depoly on Heroku
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
`curl -i -X PATCH -d '{"active": false}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/users/1`
### 解封使用者(id = 1)
`curl -i -X PATCH -d '{"active": true}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/users/1`
### 隱藏留言(id = 2)
`curl -i -X PATCH -d '{"hidden": true}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2`
### 顯示留言(id = 2)
`curl -i -X PATCH -d '{"hidden": false}' -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/2`
### 取得所有留言(包含隱藏留言)
`curl -i -X GET -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts`
### 搜尋留言內容，回傳留言
`curl -i -X GET -H "Cookie: <login-returned-set-cookies>" -H 'Content-Type: application/json' 127.0.0.1:8080/admin/posts/search?keyword="123"`

# To-do-list
- [x] Use bcrypt on password accessing
- [x] Deploy on cloud service
- [x] Testing
