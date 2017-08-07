package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	apiBaseUrl = "http://0.0.0.0:8081"
)

func StartTestServer(t *testing.T) {
	go func() {
		r := gin.Default()

		store := sessions.NewCookieStore([]byte("secret"))
		store.Options(sessions.Options{
			MaxAge: int(30 * time.Minute), //30min
			Path:   "/",
		})

		r.Use(sessions.Sessions("mysession", store))
		r.GET("/pre", preSession)
		r.GET("/do", DoSomething)
		r.Run(":8081")
		fmt.Println("start api server success.")
	}()

	t.Log("waiting 2 second for server startup")
	time.Sleep(2 * time.Second)
}

func TestSession(t *testing.T) {
	StartTestServer(t)
	cookie := setAndGetCookie()

	u := url.Values{}
	u.Set("user_email", "test@mail.com")
	testUrl := apiBaseUrl + "/do?" + u.Encode()

	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		logrus.Fatalf("generate req of url %s error: %v", testUrl, err)
	}
	req.AddCookie(cookie)
	fmt.Printf("header of test session is %v\n", req.Header)
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		t.Error(err)
	}
	cookies := resp.Cookies()
	fmt.Printf("cookie of do session request is %v\n", cookies)
	fmt.Println(resp.StatusCode)
	fmt.Printf("resp header is %v\n", resp.Header)
}

func setAndGetCookie() *http.Cookie {
	resp, err := http.Get(apiBaseUrl + "/pre")
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Fatalf("pre session request error")
	}

	cookies := resp.Cookies()
	if len(cookies) < 1 {
		logrus.Fatalf("cookie length expected to be 1")
	}

	fmt.Printf("response headers of set cookie is %v\n", resp.Header)
	return cookies[0]
}
