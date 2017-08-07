GO_BUILD_FLAGS=

# init:
# 	go get github.com/gin-gonic/gin
# 	go get github.com/gin-contrib/sessions
all:
	go build $(GO_BUILD_FLAGS) 
test:
	go test --cover -test.v
