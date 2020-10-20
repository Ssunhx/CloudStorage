package mq

var done chan bool

func StartConsumer(qName, cName string, callback func(msg []byte) bool) {
	// 通过 channel.Consume 获取消息信道
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}

	done = make(chan bool)

	go func() {
		// 循环获取队列消息
		for msg := range msgs {
			// 调用 callback 方法处理新消息
			processSuc := callback(msg.Body)
			if processSuc {
				// TODO 将任务写到另一个队列，用于异常重试
			}
		}
	}()

	// 接收 done 信号，没有信号会一直阻塞，避免函数退出
	<-done
	// 关闭通道
	channel.Close()
}

// 停止监听队列
func StopConsume() {
	done <- true
}
