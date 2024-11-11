package main

// 模拟一个发短信的服务
// 就是我们有一个订单系统，他的任务就是给客户发确认短信。我们有一个SMSService 接口，
// 提供 SendSMS 方法来发送短信。在测试中，我们不想真的给客户发送短信，而是通过 mock 来模拟 SMSService。

// 我现在的理解是需要一个interface接口的发送短信，然后又两个方法实现了这个接口，一个是正常的调用和返回，另一个是mock server模拟调用和返回

// SMSService 这里面我看到GPT建议写成大写字母开头的的一个Service，这个Service 字定义的就很好，叫SMSServer,看名字就知道是处理SMS的业务
// 那么这个SMS 业务就可以又很多行为，比如SendMsg
type SMSService interface {
	SendMsg(phoneNumber, message string) error
}

// 实现订单系统
type OrderService struct {
	smsService SMSService // 这是一个接口，
}

// 创建一个新的服务订单, 我们可以看到创建他的时候用的是接口
func NewOrderService(smsService SMSService) *OrderService {
	return &OrderService{smsService: smsService}
}

// 创建一个下订单的函数
func (o *OrderService) PlaceOrder(phoneNumber string) error {
	message := "Your order has been placed successfully!"
	return o.smsService.SendMsg(phoneNumber, message)
}
