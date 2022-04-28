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
```bash
curl -i -X GET -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/posts
```
### 註冊
```bash
curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/register
```
### 登入
```bash
cookie=$(curl -i -X POST -d '{"username":"test_login","password":"test_password"}' -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/login | grep -Po '(?<=Set-Cookie: ).+(?=\r)')
```
### 發布留言
```bash
curl -i -X POST -d '{"title":"test_title2","content":"test_content2"}' --cookie "$cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/posts
```
### 回覆留言(id = 2)
```bash
curl -i -X POST -d '{"content":"test_comment"}' --cookie "$cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/posts/2/comments
```

## 管理員
### 登入
```bash
admin_cookie=$(curl -i -X POST -d '{"username":"admin","password":'"\"$ADMIN_PASSWD\""'}' -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/login | grep -Po '(?<=Set-Cookie: ).+(?=\r)')
```
### 停權使用者(id = 2)
```bash
curl -i -X PATCH -d '{"active": false}' --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/users/2
```
### 解封使用者(id = 2)
```bash
curl -i -X PATCH -d '{"active": true}' --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/users/2
```
### 隱藏留言(id = 2)
```bash
curl -i -X PATCH -d '{"hidden": true}' --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/posts/2
```
### 顯示留言(id = 2)
```bash
curl -i -X PATCH -d '{"hidden": false}' --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/posts/2
```
### 取得所有留言(包含隱藏留言)
```bash
curl -i -X GET --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/posts
```
### 搜尋留言內容，回傳留言
```bash
curl -i -X GET --cookie "$admin_cookie" -H 'Content-Type: application/json' bbs-restful-api.herokuapp.com/admin/posts/search?keyword="123"
```

# To-do-list
- [x] Use bcrypt on password accessing
- [x] Deploy on cloud service
- [x] Testing
