package main

import (
	"fmt"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	serviceMap map[string]interface{}
}

func NewContainer() *Container {
	return &Container{ serviceMap: make(map[string]interface{}) }
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.serviceMap[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	ctor, exists := c.serviceMap[name]
	if !exists {
		return nil, errors.New("service dosn't exists in map")
	}
	ctor_func, ok := ctor.(func() interface{})
	if !ok {
		return nil, errors.New("constructor not a fucntion for service " + name)
	}
	return ctor_func(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
	
	fmt.Printf("%v\n", err)

	container.RegisterType("AnotherService", 5)
	anotherService, err := container.Resolve("AnotherService")
	assert.Error(t, err)
	assert.Nil(t, anotherService)

	fmt.Printf("%v\n", err)
	

}
