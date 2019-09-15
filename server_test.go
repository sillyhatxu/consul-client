package consul

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	address     = "localhost:8500"
	clusterName = "test-http"
	clusterHost = "host.docker.internal"
	//clusterHost = "localhost"
	clusterPort = 8801
)

var consulServer = NewConsulServer(address, clusterName, clusterHost, clusterPort, CheckType(HealthCheckHttp))

func TestPutCache(t *testing.T) {
	err := consulServer.Put("test-cache", []byte("test"))
	assert.Nil(t, err)
}

func TestGetListCache(t *testing.T) {
	err := consulServer.Put("test-cache1", []byte("test1"))
	assert.Nil(t, err)
	err = consulServer.Put("test-cache2", []byte("test2"))
	assert.Nil(t, err)
	err = consulServer.Put("test-cache3", []byte("test3"))
	assert.Nil(t, err)
	err = consulServer.Put("test-cache4", []byte("test4"))
	assert.Nil(t, err)
	err = consulServer.Put("test-cache5", []byte("test5"))
	assert.Nil(t, err)
	pairs, ok := consulServer.List("test-cache")
	assert.EqualValues(t, ok, true)
	for i, pair := range pairs {
		logrus.Infof("i : %d, pair : %s", i, string(pair.Value))
	}

}

func testKey(i int) string {
	return fmt.Sprintf("test-cache%d", i)
}

func TestGetCache(t *testing.T) {
	value, ok := consulServer.Get("test-cache")
	assert.EqualValues(t, ok, true)
	assert.EqualValues(t, value, []byte("test"))

	value, ok = consulServer.Get("test-cache2")
	assert.EqualValues(t, ok, false)
	assert.Nil(t, value)
}

func TestDeleteCache(t *testing.T) {
	err := consulServer.Delete("test-cache")
	assert.Nil(t, err)
}
