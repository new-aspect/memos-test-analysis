package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockSMSService struct {
	mock.Mock
}

func (m *MockSMSService) SendMsg(phoneNumber string, message string) error {
	args := m.Called(phoneNumber, message)
	return args.Error(1)
}

func TestOrderService_PlaceOrder(t *testing.T) {
	// 创建一个 MockSMSService
	mockSMS := new(MockSMSService)

	// 设置mock的期望值和返回值
	mockSMS.On("SendMsg", "1234567890", "Your order has been placed successfully!").
		Return(nil)

	// 创建 OrderService, 传入 mock 的 SMSService
	orderService := NewOrderService(mockSMS)

	// 调用 PlaceOrder 方法
	err := orderService.PlaceOrder("1234567890")

	// 验证结果
	assert.NoError(t, err)
	mockSMS.AssertExpectations(t)

}
