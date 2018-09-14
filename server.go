package consul

import (
	log "github.com/xushikuan/microlog"
	consulapi "github.com/hashicorp/consul/api"
	dockerClient "github.com/fsouza/go-dockerclient"
	dockerSwarm "github.com/docker/docker/api/types/swarm"
	"github.com/xushikuan/docker-client"
	"net"
	"strconv"
	"math/rand"
	"errors"
)

const (
	default_consul_name = "consul"
	timeout                           = "3s"
	interval                          = "10s"
	deregister_critical_service_after = "30s"
	default_endpoint = "unix:///var/run/docker.sock"
	default_consul_port = "8500"
)


type ConsulServer struct {
	EnviromentName string

	Name string

	Port int

	HealthURL string

	ConsulAddress string

	DockerConfig *DockerConfig

	Config *Config
}

type DockerConfig struct {

	IsDocker bool

	ConsulName string

	Endpoint string

	StackName string

	NetworkName string
}

type Config struct {
	Timeout string

	Interval string

	DeregisterCriticalServiceAfter string
}

func NewConsulServer(enviromentName string,name string,port int,healthURL string) *ConsulServer {
	return &ConsulServer{
		EnviromentName:enviromentName,
		Name:name,
		Port:port,
		HealthURL:healthURL,
		DockerConfig:&DockerConfig{IsDocker:true,ConsulName:default_consul_name,Endpoint:default_endpoint,StackName:enviromentName,NetworkName:enviromentName+"_default"},
		Config:&Config{Timeout:timeout,Interval:interval,DeregisterCriticalServiceAfter:deregister_critical_service_after},
	}
}

func (server *ConsulServer) SetConsulAddress(consulAddress string){
	server.ConsulAddress = consulAddress
	server.DockerConfig.IsDocker = false
}

func (server *ConsulServer) SetDockerConfig(config *DockerConfig){
	server.DockerConfig = config
}

func (server *ConsulServer) SetConfig(config *Config){
	server.Config = config
}

func (server *ConsulServer) RegisterServer() {
	config := consulapi.DefaultConfig()
	if server.DockerConfig.IsDocker{
		consulAddress,err := server.getConsulAddressFromDockerSwarm()
		if err != nil {
			log.Panic("Get consul address form docker swarm error : ", err)
		}
		config.Address = consulAddress
	}else{
		config.Address = server.ConsulAddress
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Panic("New consul client error : ", err)
	}
	ip,err := localIP()
	if err != nil{
		log.Panic("Register server error.",err)
	}
	//http : http://172.28.2.106:18002/health
	//GRPC : 172.28.2.106:18002/payment
	healthCheck := "http://" + ip + ":" + strconv.Itoa(server.Port) + server.HealthURL

	reg := &consulapi.AgentServiceRegistration{
		ID:      ip + ":" + strconv.Itoa(server.Port),
		Name:    server.Name,
		Port:    server.Port,
		Address: ip,
		Tags:    []string{server.Name},
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           healthCheck,
			Timeout:                        server.Config.Timeout,
			Interval:                       server.Config.Interval,                       //健康检查间隔
			GRPC:                           ip + "/" + server.Name,                       // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: server.Config.DeregisterCriticalServiceAfter, //注销时间，相当于过期时间,check失败后30秒删除本服务
		},
	}
	if err := client.Agent().ServiceRegister(reg); err != nil {
		log.Panic("register server error : ", err)
	}
}

func (server *ConsulServer) getDockerService() (*dockerSwarm.Service,error) {
	client := docker.NewDockerClient(server.DockerConfig.Endpoint)
	serviceArray,err := client.ServiceLSFilter(server.DockerConfig.ConsulName)
	if err != nil{
		log.Error("Docker service ls --filter error.",err)
		return nil,err
	}
	if serviceArray == nil && len(serviceArray) < 1 {
		return nil, errors.New("Unknow docker service.")
	}
	return &serviceArray[rand.Intn(len(serviceArray))],nil
}


func (server *ConsulServer) getDockerNetwork() (*dockerClient.Network,error) {
	client := docker.NewDockerClient(server.DockerConfig.Endpoint)
	networkArray,err := client.NetworkLSFilter(server.DockerConfig.NetworkName)
	if err != nil{
		log.Error("Docker network ls --filter error.",err)
		return nil,err
	}
	return &networkArray[rand.Intn(len(networkArray))],nil
}



func (server *ConsulServer) getConsulIPFromDockerSwarm(virtualIPArray []dockerSwarm.EndpointVirtualIP) (string,error) {
	dockerNetwork,err := server.getDockerNetwork()
	if err != nil{
		log.Error("Get docker network error.",err)
		return "",err
	}
	for _,virtualIP := range virtualIPArray{
		if virtualIP.NetworkID == dockerNetwork.ID{
			log.Info("Consul IP : ",virtualIP.Addr)
			ip,err := cidrToIP(virtualIP.Addr)
			if err != nil{
				return "",err
			}
			return ip,nil
		}
	}
	return "",errors.New("Uknow consul ip.")
}

func (server *ConsulServer) getConsulAddressFromDockerSwarm() (string,error) {
	dockerService,err := server.getDockerService()
	if err != nil{
		log.Error("Get docker service error.",err)
		return "",err
	}
	ip,err := server.getConsulIPFromDockerSwarm(dockerService.Endpoint.VirtualIPs)
	return ip + ":" + default_consul_port,errors.New("Uknow consul ip.")
}

func localIP() (string,error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Error("Get local ip error.",err)
		return "",err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP := localAddr.IP.String()
	log.Info("Local IP : ",localIP)
	return localIP,nil
}

//cidr : "10.0.0.0/8"
func cidrToIP(cidr string) (string,error) {
	ipAddress, _, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Error("ParseCIDR error.",err)
		return "",err
	}
	return ipAddress.String(),nil
}
