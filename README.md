# gin session problem
## 环境
1. gin版本
```
{
    "checksumSHA1": "86tapazS8gfJ5JRCxVNTTDkUZwM=",
    "path": "github.com/gin-gonic/gin",
    "revision": "bbd4dfee5056087c640a75c6cc21567f4f47585d",
    "revisionTime": "2017-07-12T07:01:46Z"
}
```
2. golang 1.8.0
3. gin中设置session代码
```
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
```
4. preSession(这个函数目的是向启动的gin.Engine中添加一条session记录)
```
func preSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("test@mail.com", "test.access.token")
	session.Save()
	c.JSON(http.StatusOK, nil)
}
```
5. DoSomethine函数接受email参数，读取session中对应email的value
```
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
```
6. 测试方法，先调用getPreSession()设置一条session，然后访问`/do`查看是否拿到session。
```
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
```
7. 结果，注意**[DoSomethine] user access token is %!s(<nil>)**表明没有拿到session。
```
➜  session make test 
go test --cover -test.v
=== RUN   TestSession
cookie of pre session request is [mysession=MTUwMTgxNDYxNnxEdi1CQkFFQ180SUFBUkFCRUFBQU9QLUNBQUVHYzNSeWFXNW5EQThBRFhSbGMzUkFiV0ZwYkM1amIyMEdjM1J5YVc1bkRCTUFFWFJsYzNRdVlXTmpaWE56TG5SdmEyVnV8cc0AwUddF90uWUxOF8BONUo-tFGouRdSOfj62m9U8sE=; Path=/; Expires=Wed, 09 Jan 1771 04:16:32 GMT; Max-Age=1800000000000]
test url is http://0.0.0.0:8081/do?user_email=test%40mail.com
[DoSomethine] user access token is %!s(<nil>)
[GIN] 2017/08/04 - 10:43:36 | 200 |      99.642µs |       127.0.0.1 |  GET     /do
200
--- PASS: TestSession (2.00s)
        main_test.go:36: waiting 2 second for server startup
PASS
coverage: 48.0% of statements
ok      ginlab/session  2.010s
```

## 问题
为什么`/pre`设置session后，请求`/do`却拿不到session？是没session没有配置正确还是其它原因？

## 解决方法
**本质上是Cookie和Session的问题，因为不同的库实现机制不同，所以需要充分了解库文档和代码正确配置和使用**

在基于浏览器的测试情况下能测试通过是因为浏览器会在请求中加上本地的Cookie。而在我的`go test case`中并不是浏览器的使用场景，所以需要主动加上cookie。
```
req, err := http.NewRequest("GET", testUrl, nil)
if err != nil {
	logrus.Fatalf("generate req of url %s error: %v", testUrl, err)
}
req.AddCookie(cookie)
```

总结起来，还是自己关于[Cookie和Session](https://en.wikipedia.org/wiki/HTTP_cookie)的基础知识薄弱，没有深入了解[gin-contrib/session](https://github.com/gin-contrib/sessions)代码