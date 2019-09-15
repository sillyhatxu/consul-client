package consul

import "time"

const (
	DefaultHealthCheckGRPCServerName = "grpc.health.v1.Health"
)

const (
	defaultTimeout                        = 3 * time.Second
	defaultInterval                       = 11 * time.Second
	defaultDeregisterCriticalServiceAfter = 3 * time.Second
)

const (
	HealthCheckGRPC = iota + 1
	HealthCheckHttp
)
