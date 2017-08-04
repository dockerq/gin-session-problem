package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/clear", clear)
	r.GET("/pre", preSession)
	r.GET("/do", DoSomethine)
	r.Run(":8000")
}

func preSession(c *gin.Context) {
	userAccessToken := "test.access.token"
	session := sessions.Default(c)
	session.Set("test@mail.com", userAccessToken)
	session.Save()
	fmt.Printf("[preSession] user access token :%s has been saved to session\n", userAccessToken)
	c.JSON(http.StatusOK, nil)
}

func DoSomethine(c *gin.Context) {
	userEmail := c.Query("user_email")
	if userEmail == "" {
		panic("can not get user email")
	}
	session := sessions.Default(c)
	userAccessToken := session.Get(userEmail)
	fmt.Printf("[DoSomethine] user access token is %s\n", userAccessToken)
	c.JSON(http.StatusOK, nil)
}

func clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "clear session",
	})
}
