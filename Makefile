init:
	go get github.com/gin-gonic/gin
	go get github.com/gin-contrib/sessions
test:
	go test --cover -test.v
