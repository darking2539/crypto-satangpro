package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func InitToSendData() (*amqp.Connection, *amqp.Channel, *amqp.Queue, error) {
	
	rabbitMqUrl := os.Getenv("RABBIT_MQ_URL")
	
	conn, err := amqp.Dial(rabbitMqUrl)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel", err)
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		"eth_queue", // name of the queue
		true,          // durable (queue will survive server restarts)
		false,         // delete when unused (queue will be deleted when there are no more subscribers)
		false,         // exclusive (only allow access by the declarer's connection)
		false,         // no-wait (wait for a confirmation from the server)
		nil,           // arguments
	)

	if err != nil {
		log.Println("Failed to declare a queue", err)
		return nil, nil, nil, err
	}


	return conn, ch, &q, nil
}

func SendingData(body interface{}, conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue) {
	
	// Create a message publishing function
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println("Failed to Marshal Json")
		return
	}

	err = ch.Publish(
		"",     // exchange (empty string since we use a default exchange)
		q.Name, // routing key (queue name)
		false,  // mandatory (wait for at least one queue to be bound to the exchange)
		false,  // immediate (do not wait for consumers to process the message)
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         bodyBytes,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Println("Failed to publish a message", err.Error())
		return
	}
}

func InitSubcribeCode(callBack func(bodyResp []byte)) () {
	
	// Connect to RabbitMQ server
	rabbitMqUrl := os.Getenv("RABBIT_MQ_URL")
	conn, err := amqp.Dial(rabbitMqUrl)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	// Declare a queue for email messages
	q, err := ch.QueueDeclare(
		"eth_queue", // name of the queue
		true,          // durable (queue will survive server restarts)
		false,         // delete when unused (queue will be deleted when there are no more subscribers)
		false,         // exclusive (only allow access by the declarer's connection)
		false,         // no-wait (wait for a confirmation from the server)
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // name of the queue
		"",     // consumer tag (empty string will generate a unique tag)
		false,  // auto-ack (do not automatically acknowledge messages)
		false,  // exclusive (only allow access by the declarer's connection)
		false,  // no-local (do not receive messages published by this connection)
		false,  // no-wait (wait for a confirmation from the server)
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}
	

	// Start a goroutine to handle incoming messages
	go func() {
		for d := range msgs {
			callBack(d.Body)
			d.Ack(false)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-make(chan struct{})
}