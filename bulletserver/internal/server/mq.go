package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

var mqconn *amqp.Connection
var mqchannel *amqp.Channel
var mqqueues []string

const mqex = "bullets_ex"

func initMQ() {
	mqconn, err = amqp.Dial(configs.MQ_Addr)
	pkg.FailOnError(err, "Failed to connect to RabbitMQ")

	mqchannel, err = mqconn.Channel()
	pkg.FailOnError(err, "Failed to create a channel in default exchange")

	err = mqchannel.ExchangeDeclare(
		mqex,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	pkg.FailOnError(err, "Failed to declare the mq exchange")
	for i := 0; i < configs.MQ_QueueNumber; i++ {
		queuename := fmt.Sprintf("bullet_queue_%d", i)
		mqqueues = append(mqqueues, queuename)
		_, err = mqchannel.QueueDeclare(
			queuename,
			true,
			false,
			false,
			false,
			nil,
		)
		pkg.FailOnError(err, fmt.Sprintf("Failed to declare the queue #%d of %d\n", i, configs.MQ_QueueNumber))

		err = mqchannel.QueueBind(
			queuename,
			queuename,
			mqex,
			false,
			nil,
		)
		pkg.FailOnError(err, fmt.Sprintf("Failed to bind the queue #%d of %d to the channel\n", i, configs.MQ_QueueNumber))
	}
}

func StartMQ(ctx *context.Context) {
	var wg sync.WaitGroup
	for _, qn := range mqqueues {
		wg.Add(1)
		go func(queuename string) {
			defer wg.Done()
			msgs, err := mqchannel.Consume(
				queuename,
				"",
				false, // 手动确认
				false,
				false,
				false,
				nil,
			)
			pkg.FailOnError(err, fmt.Sprintf("Failed to consume from queue %s", queuename))

			for msg := range msgs {
				select {
				case <-(*ctx).Done():
					log.Printf("Gracefully shutdown the %s\n", queuename)
					return
				default:
					var bullets []configs.Bullet
					if err := json.Unmarshal(msg.Body, &bullets); err != nil {
						log.Printf("Failed to Unmarshal the message at %s", queuename)
						msg.Nack(false, false)
					} else {
						WriteToDrawingFile(&bullets, queuename)
						msg.Ack(false)
					}
				}
			}
		}(qn)
	}
	log.Println("Start MQ!")
	wg.Wait()
	mqchannel.Close()
	mqconn.Close()
	log.Printf("Gracefully shutdown all consumers\n")
}

func SendToMQwithoutRoutingKey(bullets *[]configs.Bullet) {
	body, err := json.Marshal(bullets)
	pkg.FailOnError(err, "Failed on Marshaling bullets when sending to the mq")
	num := rand.Intn(configs.MQ_QueueNumber)
	routeKey := mqqueues[num]
	SendToMQwithRoutingKey(body, routeKey)
}

func SendToMQwithRoutingKey(body []byte, routeKey string) {
	err = mqchannel.PublishWithContext(
		context.Background(),
		mqex,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	pkg.FailOnError(err, fmt.Sprintf("Failed on sending to the mq %s", routeKey))
}
