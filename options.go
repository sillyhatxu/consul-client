package consul

import "time"

type Config struct {
	timeout                        time.Duration
	interval                       time.Duration
	deregisterCriticalServiceAfter time.Duration
	checkType                      int
	healthURL                      string
}

type Option func(*Config)

func CheckType(checkType int) Option {
	return func(c *Config) {
		c.checkType = checkType
	}
}

func HealthURL(healthURL string) Option {
	return func(c *Config) {
		c.healthURL = healthURL
	}
}

func Timeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func Interval(interval time.Duration) Option {
	return func(c *Config) {
		c.interval = interval
	}
}

func DeregisterCriticalServiceAfter(deregisterCriticalServiceAfter time.Duration) Option {
	return func(c *Config) {
		c.deregisterCriticalServiceAfter = deregisterCriticalServiceAfter
	}
}
