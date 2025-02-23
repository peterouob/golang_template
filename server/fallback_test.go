package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterUserService_NotFound(t *testing.T) {
	serviceName := "nonexistent_service"

	_, exists := userService[serviceName]
	assert.False(t, exists, "測試應該在 serviceName 不存在時執行")

	service := RegisterUserService(serviceName)

	assert.NotNil(t, service)
	assert.Equal(t, service.ServiceName, serviceName)
}
