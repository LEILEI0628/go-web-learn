package Sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest // 设置偏移量为从最旧开始消费
	// 消费者归于消费者组，消费者组可以理解为业务
	consumer, err := sarama.NewConsumerGroup(addrs, "TestGroup", cfg)
	require.NoError(t, err)

	start := time.Now()
	// 带超时的context（10s）
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	// 另一种写法
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*15, cancel)
	err = consumer.Consume(ctx, []string{"TestTopic"}, testConsumerGroupHandler{})
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct {
}

func (h testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("SetUp")
	// TODO 此处似乎有BUG
	// topic => 偏移量
	partitions := session.Claims()["TestTopic"]
	for _, part := range partitions { // 遍历所有分区
		// sarama.OffsetOldest：从最晚的开始消费
		session.ResetOffset("TestTopic", part, sarama.OffsetOldest, "")
	}
	return nil
}

func (h testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

func (h testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// session：与kafka的会话（从建立连接到断开连接的时间）
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value)) // 单个消费模式
		//var bizMsg MyBizMsg
		//err := json.Unmarshal(msg.Value, &bizMsg)
		//if err != nil {
		//	// 消费消息出错：大多数情况下是重试，记录日志
		//	log.Println("消费消息出错")
		//	continue
		//  return
		//}
		log.Printf("消费消息: Topic=%s, Partition=%d, Offset=%d, Value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		session.MarkMessage(msg, "") // 标记为消费成功（提交）
		//session.Commit()
	}
	//single(claim, session) // 单个消费模式
	return nil // msgs被人关闭（退出消费逻辑，一般是服务关闭时）
}

func sync(msgs <-chan *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	// 最简单的异步消费实现：goroutine
	// 单条消费，单条提交（没有控制住goroutine数量，如果消费较慢生产较快会被打爆）
	//for msg := range msgs {
	//	message := msg // for循环中的闭包问题
	//	go func() {
	//		log.Println(string(message.Value)) // 模拟消费msg
	//		session.MarkMessage(message, "")   // 异步消费提交可能是乱序的
	//	}()
	//}

	// 批量消费，批量提交
	const batchSize = 10
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				//goto label1 // 不推荐使用label写法
				done = true // 发生超时
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return // 代表消费者被关闭
				}
				last = msg
				eg.Go(func() error {
					time.Sleep(time.Second) // 模拟消费
					// 可以在这里重试
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		//label1:
		//	log.Println("发生超时")
		err := eg.Wait()
		if err != nil {
			// 可以在这里统一重试
			// 记录日志
			continue
		}
		// 可以只提交最后一条（sarama会批量提交全部消息）
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
}

func single(claim sarama.ConsumerGroupClaim, session sarama.ConsumerGroupSession) {
	msgs := claim.Messages()
	for msg := range msgs {
		//msg := msg
		func(msg *sarama.ConsumerMessage) {
			log.Println(string(msg.Value)) // 单个消费模式
			//var bizMsg MyBizMsg
			//err := json.Unmarshal(msg.Value, &bizMsg)
			//if err != nil {
			//	// 消费消息出错：大多数情况下是重试，记录日志
			//	log.Println("消费消息出错")
			//	continue
			//  return
			//}
			session.MarkMessage(msg, "") // 标记为消费成功（提交）
		}(msg)
	}
}

//type MyBizMsg struct {
//	Name string
//}
