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
		r.GET("/do", DoSomethine)
		r.Run(":8081")
		fmt.Println("start api server success.")
	}()

	t.Log("waiting 2 second for server startup")
	time.Sleep(2 * time.Second)
}

func TestSession(t *testing.T) {
	StartTestServer(t)
	getPreSession()

	u := url.Values{}
	u.Set("user_email", "test@mail.com")
	testUrl := apiBaseUrl + "/do?" + u.Encode()
	fmt.Printf("test url is %s\n", testUrl)

	resp, err := http.Get(testUrl)
	defer resp.Body.Close()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp.StatusCode)
}

func getPreSession() {
	resp, err := http.Get(apiBaseUrl + "/pre")
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Fatalf("pre session request error")
	}

	fmt.Printf("cookie of pre session request is %v\n", resp.Cookies())
}
