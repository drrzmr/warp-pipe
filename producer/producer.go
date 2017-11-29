package producer

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// Run kafka producer
func Run() (err error) {

	config := sarama.NewConfig()
	// successfully delivered messages will be returned on the success channel
	config.Producer.Return.Successes = true
	// messages that failed to deliver will be returned on the errors channel
	config.Producer.Return.Errors = true
	// wait for all in-sync replicas to commit before responding
	config.Producer.RequiredAcks = sarama.WaitForAll

	brokers := []string{"kafka:29092"}

	producer, err := sarama.NewAsyncProducer(brokers, config)

	if err != nil {
		return err
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	var (
		enqueued int
		errors   int
	)

	doneCh := make(chan struct{})

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			msg := &sarama.ProducerMessage{
				Topic: "test",
				Value: sarama.StringEncoder("surprise!"),
			}

			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Println("Message sent")
				fmt.Printf("Enqueued: %d; Erros: %d\n", enqueued, errors)
			case <-producer.Successes():
				enqueued--
				fmt.Println("Message delivered")
				fmt.Printf("Enqueued: %d; Erros: %d\n", enqueued, errors)
				doneCh <- struct{}{}
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			}
		}
	}()

	<-doneCh
	return nil
}
