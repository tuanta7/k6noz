package kafka

type Broker interface {
	Publish(topic string, msg []byte) error
}
