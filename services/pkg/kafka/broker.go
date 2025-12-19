package kafka

type Broker interface{}

type Publisher interface {
	Publish(topic string, msg []byte) error
}
