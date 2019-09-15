# consul-client

[example](https://github.com/sillyhatxu/consul-client/example)

[server-example](https://github.com/sillyhatxu/consul-client/server_test.go)


### Register (health check grpc)

```go
package main

import (
	"fmt"
	"github.com/sillyhatxu/consul-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	hv1 "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

const (
	address     = "localhost:8500"
	clusterName = "test-grpc"
	clusterHost = "host.docker.internal"
	//clusterHost = "localhost"
	clusterPort = 8802
)

var consulServer = consul.NewConsulServer(address, clusterName, clusterHost, clusterPort, consul.CheckType(consul.HealthCheckGRPC))

func main() {
	err := consulServer.Register()
	if err != nil {
		panic(err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", clusterPort))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus(consul.DefaultHealthCheckGRPCServerName, hv1.HealthCheckResponse_SERVING)
	hv1.RegisterHealthServer(server, healthServer)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

### Register (health check http)

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sillyhatxu/consul-client"
	"github.com/sirupsen/logrus"
)

const (
	address     = "localhost:8500"
	clusterName = "test-http"
	clusterHost = "host.docker.internal"
	//clusterHost = "localhost"
	clusterPort = 8801
)

var consulServer = consul.NewConsulServer(address, clusterName, clusterHost, clusterPort, consul.CheckType(consul.HealthCheckHttp))

func main() {
	err := consulServer.Register()
	if err != nil {
		panic(err)
	}
	router := SetupRouter()
	logrus.Infof("init http")
	err = router.Run(fmt.Sprintf(":%d", clusterPort))
	if err != nil {
		panic(err)
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "UP", "message": "OK"}) })
	return router
}

```