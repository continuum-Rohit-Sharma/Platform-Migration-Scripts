package kafka

// import (
// 	"log"

// 	"github.com/ContinuumLLC/platform-common-lib/src/kafka/encode"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// //ProducerCommandImpl implements a ProducerCommand
// type confluentProducerCommandImpl struct {
// 	producer        *kafka.Producer
// 	deliveryChannel chan kafka.Event
// }

// //NewProducer create a new confluentProducer
// func (pc *confluentProducerCommandImpl) NewProducer(brokerAddress []string) error {

// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokerAddress[0]})
// 	if err != nil {
// 		log.Printf("Failed to create producer: %s\n", err)
// 	}
// 	pc.producer = p
// 	return err

// }

// //Close closes connnection to kafka
// func (pc *confluentProducerCommandImpl) Close() error {
// 	close(pc.deliveryChannel)
// 	pc.producer.Close()
// 	return nil
// }

// //PushMessage pushes new Message to kafka
// func (pc *confluentProducerCommandImpl) PushMessage(topicName string, message string) (int32, int64, error) {
// 	//fmt.Println("confleunt producer")
// 	pc.deliveryChannel = make(chan kafka.Event)
// 	err := pc.producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny}, Value: []byte(message)}, pc.deliveryChannel)
// 	e := <-pc.deliveryChannel
// 	m := e.(*kafka.Message)
// 	if m.TopicPartition.Error != nil {
// 		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
// 		return 0, 0, m.TopicPartition.Error
// 	}
// 	return 0, 0, err
// }

// func (pc *confluentProducerCommandImpl) PushMessageEncoder(topicName string, message encode.Encoder) (int32, int64, error) {
// 	return 0, 0, nil
// }

// //IsConnected checks if the Producer Connected
// func (pc *confluentProducerCommandImpl) IsConnected() bool {
// 	if pc.producer != nil {
// 		return true
// 	}
// 	return false
// }

// //confluentConsumerCommandImpl implements ConsumerCommand
// type confluentConsumerCommandImpl struct {
// 	consumer *kafka.Consumer
// }

// func (cc *confluentConsumerCommandImpl) NewConsumer(brokerAddress []string, GroupID string, Topics []string) error {
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers":    brokerAddress[0],
// 		"group.id":             GroupID,
// 		"session.timeout.ms":   6000,
// 		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}})

// 	if err != nil {
// 		return err
// 	}
// 	err = c.SubscribeTopics(Topics, nil)

// 	if err != nil {
// 		return err
// 	}
// 	cc.consumer = c
// 	return err
// }

// func (cc *confluentConsumerCommandImpl) IsConnected() bool {
// 	if cc.consumer != nil {
// 		return true
// 	}
// 	return false
// }

// func (cc *confluentConsumerCommandImpl) Close() error {
// 	return cc.consumer.Close()
// }

// func (cc *confluentConsumerCommandImpl) PullMessage(consumerHandler ConsumerHandler) {
// 	for {
// 		ev := cc.consumer.Poll(100)
// 		if ev == nil {
// 			continue
// 		}

// 		switch e := ev.(type) {
// 		case *kafka.Message:

// 			log.Printf("Message on %s:\n%s\n",
// 				e.TopicPartition, string(e.Value))
// 			consumerMessage := ConsumerMessage{Message: string(e.Value), Offset: 0, Partition: 0, ReceivedDateTimeUTC: e.Timestamp, Topic: ""}
// 			go consumerHandler(consumerMessage)

// 		case kafka.PartitionEOF:
// 			log.Printf("Reached %v\n", e)
// 		case kafka.Error:
// 			log.Printf("Consumer Error: %v\n", e)
// 		default:
// 			log.Printf("Ignored %v\n", e)
// 		}

// 	}
// }
