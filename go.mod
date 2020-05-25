module robot

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.4
	github.com/gorilla/websocket v1.4.1
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
