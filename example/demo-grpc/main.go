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
	//hv1.RegisterHealthServer(server, &HealthImpl{})
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//type HealthImpl struct{}
//
//func (h *HealthImpl) Watch(*hv1.HealthCheckRequest, hv1.Health_WatchServer) error {
//	panic("implement me")
//}
//
// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
//func (h *HealthImpl) Check(ctx context.Context, req *hv1.HealthCheckRequest) (*hv1.HealthCheckResponse, error) {
//	logrus.Info("health check")
//	return &hv1.HealthCheckResponse{
//		Status: hv1.HealthCheckResponse_SERVING,
//	}, nil
//}
