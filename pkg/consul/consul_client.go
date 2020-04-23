package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

//Client provides an interface for getting data out of Consul
type ClientI interface {
	// Get a Service from consul
	Service(string, string) ([]string, error)
	// Register a service with local agent
	Register(string, int) error
	// Deregister a service with local agent
	DeRegister(string) error
}

type Client struct {
	consul *consul.Client
}

//NewConsul returns a Client interface for given consul address
func NewConsulClient(addr string) (*Client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{consul: c}, nil
}

// Register a service with consul local agent
func (c *Client) Register(name string, port int) error {
	reg := &consul.AgentServiceRegistration{
		ID:   name,
		Name: name,
		Port: port,
	}
	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *Client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

//todo not work
func (c *Client) Service(service, tag string) ([]string, error) {
	addrs, _, err := c.consul.Health().Service(service, tag, true, nil)
	if len(addrs) == 0 && err == nil {
		return nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}
