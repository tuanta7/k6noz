package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type Client struct {
	client *kgo.Client
}
