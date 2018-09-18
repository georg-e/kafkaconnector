package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gogo/protobuf/proto"
	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

func decodeMessages(consumer *cluster.Consumer, dst chan *flow.FlowMessage) {
	for {
		msg, ok := <-consumer.Messages()
		if !ok {
			log.Println("Message channel closed.")
			close(dst)
		}
		consumer.MarkOffset(msg, "") // mark message as processed
		flowMsg := new(flow.FlowMessage)
		err := proto.Unmarshal(msg.Value, flowMsg)
		if err != nil {
			log.Printf("Received broken message. Unmarshalling error: %v", err)
			continue
		}
		dst <- flowMsg // TODO: investigate how unbuffered channels affect this program as a whole
	}
}

func encodeMessages(producer sarama.AsyncProducer, topic string, src chan *flow.FlowMessage) {
	for {
		binary, err := proto.Marshal(<-src)
		if err != nil {
			log.Printf("Could not encode message. Marshalling error: %v", err)
			continue
		}
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(binary),
		}
	}
}
