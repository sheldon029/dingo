package handler

import (
	"github.com/dinever/dingo/app/model"
	"github.com/dinever/dingo/app/utils"
	"github.com/dinever/golf"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
)

const name = "Shawn Ding"
const email = "dingpeixuan911@gmail.com"
const password = "passwordfortest"

func InitTestApp() *golf.Application {
	app := golf.New()

	RegisterAdminURLHandlers(app)
	RegisterHomeHandler(app)
	utils.RegisterFuncMap(app)
	app.View.FuncMap["Setting"] = model.GetSettingValue
	app.View.FuncMap["Navigator"] = model.GetNavigators
	RegisterFunctions(app)

	app.Use(golf.RecoverMiddleware, golf.SessionMiddleware)
	app.View.SetTemplateLoader("base", "view")
	app.View.SetTemplateLoader("admin", filepath.Join("..", "..", "view", "admin"))
	app.View.SetTemplateLoader("theme", filepath.Join("..", "..", "view", "default"))
	app.SessionManager = golf.NewMemorySessionManager()
	app.Error(404, NotFoundHandler)
	return app
}

func makeTestHTTPRequest(body io.Reader, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	return req
}

func mockSignUpPostRequest(email, name, password, rePassword string) *golf.Context {
	w := httptest.NewRecorder()
	app := InitTestApp()
	form := url.Values{}
	form.Add("email", email)
	form.Add("name", name)
	form.Add("password", password)
	form.Add("re-password", rePassword)
	form.Add("remember-me", "on")
	req := makeTestHTTPRequest(strings.NewReader(form.Encode()), "POST", "/signup")
	req.PostForm = form
	return golf.NewContext(req, w, app)
}

func mockContext(form url.Values, method, path string) *golf.Context {
	w := httptest.NewRecorder()
	app := InitTestApp()
	r := makeTestHTTPRequest(strings.NewReader(form.Encode()), method, path)
	r.PostForm = form
	return golf.NewContext(r, w, app)
}

func mockLogInPostContext() *golf.Context {
	w := httptest.NewRecorder()
	app := InitTestApp()
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	form.Add("remember-me", "on")
	req := makeTestHTTPRequest(strings.NewReader(form.Encode()), "POST", "/login/")
	req.PostForm = form
	return golf.NewContext(req, w, app)
}

func retrieveSetCookieValue(rec *httptest.ResponseRecorder, key string) string {
	request := &http.Request{Header: http.Header{"Cookie": rec.HeaderMap["Set-Cookie"]}}
	cookie, _ := request.Cookie(key)
	return cookie.Value
}
