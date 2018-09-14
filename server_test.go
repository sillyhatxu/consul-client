package consul

import (
	"testing"
	log "github.com/xushikuan/microlog"
	"github.com/stretchr/testify/assert"
)

func TestNewConsulServer(t *testing.T) {
	test := NewConsulServer("dev","test",8080,"")
	log.Info(test,test.Config)
	test.SetConsulAddress("127.0.0.1:8500")
	log.Info(test,test.Config)
	//test.SetDockerEndpoint("unix:///var/run/docker.sock")
	//log.Info(test,test.Config)

	test.SetConfig(&Config{Timeout:"5s",Interval:"30s",DeregisterCriticalServiceAfter:"60s"})
	log.Info(test,test.Config)
}

func TestGetLocalIP(t *testing.T) {
	ip,err := localIP()
	assert.Nil(t,err)
	log.Info(ip)
}

func TestCIDRToIP(t *testing.T) {
	ip,err := cidrToIP("10.255.3.65/16")
	assert.Nil(t,err)
	log.Info(ip)

	ip2,err := cidrToIP("10.0.0.5/24")
	assert.Nil(t,err)
	log.Info(ip2)
}