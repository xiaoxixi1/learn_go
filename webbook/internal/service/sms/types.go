package sms

import "context"

/*
*

	接口抽象：
	  1 对于初始化客户端，主要是鉴权，这本身算是某种特定实现在初始化的时候需要解决的内容
	    所以不抽象到短信发送的接口中
	  2 发送一次信息，需要的参数是：目标手机号码，appid,签名，模板，参数
	  3 appid 和签名，一般是固定的，同样可以在初始化的时候指定
*/
// 发送短信的抽象，目的是为了适配不同的短信供应商的抽象
type Service interface {
	/**
	  cxt : 上下文
	  tplId : 模板ID
	  args  : 模板参数
	  numbers : 电话号码，可能有多个
	*/
	Send(cxt context.Context, tplId string, args []string, number ...string) error
}
