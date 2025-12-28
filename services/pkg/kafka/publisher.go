package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Publisher struct {
	kafkaClient *kgo.Client
}

func NewPublisher(seeds []string) (*Publisher, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		kafkaClient: client,
	}, nil
}

func (p *Publisher) Close() {
	p.kafkaClient.Close()
}

func (p *Publisher) PublishSync(ctx context.Context, topic string, key, value []byte) error {
	results := p.kafkaClient.ProduceSync(ctx, &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	})

	if err := results.FirstErr(); err != nil {
		return err
	}

	return nil
}
