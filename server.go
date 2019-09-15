package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Server struct {
	address     string
	clusterName string
	clusterHost string
	clusterPort int
	client      *api.Client
	config      *Config
}

func NewConsulServer(address string, clusterName string, clusterHost string, clusterPort int, opts ...Option) *Server {
	//default
	config := &Config{
		timeout:                        defaultTimeout,
		interval:                       defaultInterval,
		checkType:                      HealthCheckGRPC,
		healthURL:                      "",
		deregisterCriticalServiceAfter: defaultDeregisterCriticalServiceAfter,
	}
	for _, opt := range opts {
		opt(config)
	}
	if config.healthURL == "" && config.checkType == HealthCheckHttp {
		config.healthURL = fmt.Sprintf("http://%s:%d/health", clusterHost, clusterPort)
	} else if config.healthURL == "" && config.checkType == HealthCheckGRPC {
		config.healthURL = fmt.Sprintf("%s:%d/%s", clusterHost, clusterPort, DefaultHealthCheckGRPCServerName)
	}
	return &Server{
		address:     address,
		clusterName: clusterName,
		clusterHost: clusterHost,
		clusterPort: clusterPort,
		config:      config,
	}
}

func (server *Server) Register() error {
	client, err := server.GetConsulClient()
	if err != nil {
		return err
	}
	healthURLGRPC := ""
	healthURLHTTP := ""
	if server.config.checkType == HealthCheckHttp {
		healthURLHTTP = server.config.healthURL
	} else if server.config.checkType == HealthCheckGRPC {
		healthURLGRPC = server.config.healthURL
	}
	reg := &api.AgentServiceRegistration{
		ID:      server.clusterName,
		Name:    server.clusterName,
		Port:    server.clusterPort,
		Address: server.clusterHost,
		Tags:    []string{server.clusterName},
		Check: &api.AgentServiceCheck{
			HTTP:                           healthURLHTTP,
			GRPC:                           healthURLGRPC, // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			Timeout:                        server.config.timeout.String(),
			Interval:                       server.config.interval.String(),                       //健康检查间隔
			DeregisterCriticalServiceAfter: server.config.deregisterCriticalServiceAfter.String(), //注销时间，相当于过期时间
		},
	}
	if err := client.Agent().ServiceRegister(reg); err != nil {
		logrus.Panic("register server error : ", err)
	}
	return nil
}

func (server *Server) GetConsulClient() (*api.Client, error) {
	if server.client == nil {
		config := api.DefaultConfig()
		config.Address = server.address
		logrus.Infof("consul server : %#v", server)
		logrus.Infof("consul server config : %#v", server.config)
		client, err := api.NewClient(config)
		if err != nil {
			return nil, err
		}
		server.client = client
	}
	return server.client, nil
}

func (server *Server) Get(key string) ([]byte, bool) {
	client, err := server.GetConsulClient()
	if err != nil {
		return nil, false
	}
	kvClient := client.KV()
	pair, _, err := kvClient.Get(key, nil)
	if err != nil {
		logrus.Errorf("kv client get key[%s] error; %v", key, err)
		return nil, false
	}
	if pair == nil {
		return nil, false
	}
	return pair.Value, true
}

func (server *Server) List(key string) ([]*api.KVPair, bool) {
	client, err := server.GetConsulClient()
	if err != nil {
		return nil, false
	}
	kvClient := client.KV()
	pairs, _, err := kvClient.List(key, nil)
	if err != nil {
		logrus.Errorf("kv client get key[%s] error; %v", key, err)
		return nil, false
	}
	if pairs == nil {
		return nil, false
	}
	return pairs, true
}

func (server *Server) Put(key string, value []byte) error {
	client, err := server.GetConsulClient()
	if err != nil {
		return err
	}
	kvClient := client.KV()
	pair := &api.KVPair{Key: key, Flags: 42, Value: value}
	_, err = kvClient.Put(pair, nil)
	return err
}

func (server *Server) Delete(key string) error {
	client, err := server.GetConsulClient()
	if err != nil {
		return err
	}
	kvClient := client.KV()
	_, err = kvClient.Delete(key, nil)
	return err
}

//func (server *ConsulServer) getDockerService() (*dockerSwarm.Service, error) {
//	logrus.Info("Endpoint : ", server.DockerConfig.Endpoint)
//	logrus.Info("Network Name : ", server.DockerConfig.NetworkName)
//	client := docker.NewDockerClient(server.DockerConfig.Endpoint)
//	logrus.Info("Find consul name : ", server.DockerConfig.ConsulName)
//	serviceArray, err := client.ServiceLSFilter(server.DockerConfig.ConsulName)
//	if err != nil {
//		logrus.Error("Docker service ls --filter error.", err)
//		return nil, err
//	}
//	if serviceArray == nil && len(serviceArray) < 1 {
//		return nil, errors.New("Unknow docker service.")
//	}
//	return &serviceArray[rand.Intn(len(serviceArray))], nil
//}
//
//func (server *ConsulServer) getDockerNetwork() (*dockerClient.Network, error) {
//	logrus.Info("Endpoint : ", server.DockerConfig.Endpoint)
//	logrus.Info("Network Name : ", server.DockerConfig.NetworkName)
//	client := docker.NewDockerClient(server.DockerConfig.Endpoint)
//	logrus.Info("Find network name : ", server.DockerConfig.NetworkName)
//	networkArray, err := client.NetworkLSFilter(server.DockerConfig.NetworkName)
//	if err != nil {
//		logrus.Error("Docker network ls --filter error.", err)
//		return nil, err
//	}
//	return &networkArray[rand.Intn(len(networkArray))], nil
//}
//
//func (server *ConsulServer) getConsulIPFromDockerSwarm(virtualIPArray []dockerSwarm.EndpointVirtualIP) (string, error) {
//	dockerNetwork, err := server.getDockerNetwork()
//	if err != nil {
//		logrus.Error("Get docker network error.", err)
//		return "", err
//	}
//	for _, virtualIP := range virtualIPArray {
//		if virtualIP.NetworkID == dockerNetwork.ID {
//			ip, err := cidrToIP(virtualIP.Addr)
//			if err != nil {
//				logrus.Errorf("cidrToIP error.CIDR : %v.", virtualIP.Addr, err)
//				return "", err
//			}
//			logrus.Info("Consul IP : ", ip)
//			return ip, nil
//		}
//	}
//	return "", errors.New("Uknow consul ip.")
//}
//
//func (server *ConsulServer) getConsulAddressFromDockerSwarm() (string, error) {
//	dockerService, err := server.getDockerService()
//	if err != nil {
//		logrus.Error("Get docker service error.", err)
//		return "", err
//	}
//	ip, err := server.getConsulIPFromDockerSwarm(dockerService.Endpoint.VirtualIPs)
//	if err != nil {
//		logrus.Error("Get consul ip from docker swarm error.", err)
//		return "", err
//	}
//	return ip + ":" + default_consul_port, nil
//}
//
//func localIP() (string, error) {
//	conn, err := net.Dial("udp", "8.8.8.8:80")
//	if err != nil {
//		logrus.Error("Get local ip error.", err)
//		return "", err
//	}
//	defer conn.Close()
//	localAddr := conn.LocalAddr().(*net.UDPAddr)
//	localIP := localAddr.IP.String()
//	logrus.Info("Local IP : ", localIP)
//	return localIP, nil
//}
//
////cidr : "10.0.0.0/8"
//func cidrToIP(cidr string) (string, error) {
//	ipAddress, _, err := net.ParseCIDR(cidr)
//	if err != nil {
//		logrus.Error("ParseCIDR error.", err)
//		return "", err
//	}
//	return ipAddress.String(), nil
//}
