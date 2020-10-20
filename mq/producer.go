package mq

import (
	"cloudstorage/config"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func initChannel() bool {
	// 判断 channel 是否创建
	if channel != nil {
		return true
	}

	// 获取 rabbitmq 连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		return false
	}
	// 打开 channel， 用于发布接受消息
	channel, err = conn.Channel()
	if err != nil {
		return false
	}
	return true
}

func Publish(exchange, routingKey string, msg []byte) bool {
	// 判断 channel 是否正常
	if !initChannel() {
		return false
	}
	// 消息发布
	if nil == channel.Publish(
		exchange,
		routingKey,
		false, // 如果没有对应的queue，会丢弃这条消息
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		}) {
		return true
	}
	return false
}
