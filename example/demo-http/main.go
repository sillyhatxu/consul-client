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
