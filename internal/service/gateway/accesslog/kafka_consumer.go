package accesslog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type accessLogProcessor func(accessLog AccessLog)

type KafkaConsumerService struct {
	kafkaClient *kgo.Client
	isStarted   bool
}

func New(brokers []string, topic string) (*KafkaConsumerService, error) {
	log.Printf("Validating configurations...")

	if len(brokers) == 0 {
		return nil, fmt.Errorf("missing brokers")
	}

	if topic == "" {
		return nil, fmt.Errorf("missing topic")
	}

	log.Printf("Initializing kafka consumer...")

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("api-gateway"),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		return nil, err
	}

	kcs := &KafkaConsumerService{
		kafkaClient: cl,
		isStarted:   false,
	}

	log.Printf("Initialize kafka consumer completed")

	return kcs, nil
}

func (kcs *KafkaConsumerService) Start(processor accessLogProcessor) error {
	log.Printf("Starting kafka consumer...")

	kcs.isStarted = true

	go func() {
		ctx := context.Background()

		for {
			fetches := kcs.kafkaClient.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				kcs.Stop()
				log.Fatalf(fmt.Sprint(errs))
			}

			log.Printf("A new group of access logs arrived")
			// We can iterate through a callback function.
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				// We can even use a second callback!
				p.EachRecord(func(record *kgo.Record) {
					var accessLogs []AccessLog
					err := json.Unmarshal(record.Value, &accessLogs)
					if err != nil {
						log.Printf("Faield to parse message into json: %s", string(record.Value))
					}

					for _, accessLog := range accessLogs {
						// Process each log
						processor(accessLog)
					}
				})
			})
		}
	}()

	log.Printf("Kafka consumer started")

	return nil
}

func (kcs *KafkaConsumerService) Stop() {
	log.Printf("Stoping kafka consumer...")

	kcs.isStarted = false

	kcs.kafkaClient.Close()

	log.Printf("Kafka consumer stopped")
}
