package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	client *kgo.Client
}

func (c *Client) Poll(ctx context.Context, fn func(rec *kgo.Record)) {
	c.client.PollFetches(ctx).EachRecord(fn)
}
