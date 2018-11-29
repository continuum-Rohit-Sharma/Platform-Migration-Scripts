package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/kafka/encode"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/Shopify/sarama"
	saramacluster "github.com/bsm/sarama-cluster"
)

//ProducerCommandImpl implements a ProducerCommand
type saramaProducerCommandImpl struct {
	syncProducer sarama.SyncProducer
}

//NewProducer create a new SaramaProducer
func (pc *saramaProducerCommandImpl) NewProducer(brokerAddress []string) error {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerAddress, config)
	if err != nil {
		return err
	}
	pc.syncProducer = producer
	return err
}

//Close closes connnection to kafka
func (pc *saramaProducerCommandImpl) Close() error {
	err := pc.syncProducer.Close()
	if err != nil {
		return err
	}
	pc.syncProducer = nil
	return err
}

//PushMessage pushes new Message to kafka
func (pc *saramaProducerCommandImpl) PushMessage(topicName string, message string) (int32, int64, error) {
	// //TODO: Pushed messages may be logged with DEBUG level once logging pattern in implemented in common-lib
	// msg := sarama.ProducerMessage{Topic: topicName, Value: sarama.StringEncoder(message)}
	// p, o, e := pc.syncProducer.SendMessage(&msg)
	// if e != nil {
	// 	fmt.Println("Push Error")
	// 	fmt.Println(e)
	// }
	//pc.syncProducer = nil
	return pc.PushMessageEncoder(topicName, encode.GetStringEncoder(message))
}

//PushMessage pushes new Message to kafka
func (pc *saramaProducerCommandImpl) PushMessageEncoder(topicName string, message encode.Encoder) (int32, int64, error) {
	//TODO: Pushed messages may be logged with DEBUG level once logging pattern in implemented in common-lib
	log := logging.GetLoggerFactory().Get()
	msg := sarama.ProducerMessage{Topic: topicName, Value: message}
	partition, offset, err := pc.syncProducer.SendMessage(&msg)
	if err != nil {
		log.LogWithTransactionf(logging.ERROR, "", "Error Producing Message: Error %v", err)
		return 0, 0, err

	}
	//log.LogWithTransactionf(logging.ERROR, "", " Message Producer, Partition %d, offset %d", partition, offset)
	//pc.syncProducer = nil
	return partition, offset, nil
}

//IsConnected checks if the Producer Connected
func (pc *saramaProducerCommandImpl) IsConnected() bool {
	if pc.syncProducer == nil {
		return false
	}
	return true
}

//saramaConsumerCommandImpl implements ConsumerCommand
type saramaConsumerCommandImpl struct {
	config   *sarama.Config
	consumer *saramacluster.Consumer
}

func (cc *saramaConsumerCommandImpl) NewConsumer(brokerAddress []string, GroupID string, Topics []string) error {
	config := saramacluster.NewConfig()
	config.Consumer.Offsets.Retention = 1
	consumer, err := saramacluster.NewConsumer(brokerAddress, GroupID, Topics, config)
	if err != nil {
		log.Println(err)
		return err
	}
	cc.consumer = consumer
	return err
}

func (cc *saramaConsumerCommandImpl) NewCustomConsumer(
	inOut *ConsumerKafkaInOutParams, brokerAddress []string,
	GroupID string, Topics []string) error {

	config := saramacluster.NewConfig()
	config.Consumer.Return.Errors = inOut.ReturnErrors
	config.Group.Return.Notifications = inOut.ReturnNotifications
	config.Consumer.Offsets.Initial = inOut.OffsetsInitial
	config.Consumer.Offsets.Retention = inOut.Retention
	consumer, err := saramacluster.NewConsumer(brokerAddress, GroupID, Topics, config)
	if err == nil {
		cc.consumer = consumer

		// Process Kafka cluster errors
		go func(consumer *saramacluster.Consumer) {
			for err := range consumer.Errors() {
				inOut.Errors <- err
			}

		}(consumer)

		// Process Kafka cluster notifications
		go func(consumer *saramacluster.Consumer) {
			for ntf := range consumer.Notifications() {
				inOut.Notifications <- fmt.Sprintf("%+v", ntf)
			}

		}(consumer)

	}
	return err
}

func (cc *saramaConsumerCommandImpl) IsConnected() bool {
	if cc.consumer == nil {
		return false
	}
	return true
}

func (cc *saramaConsumerCommandImpl) Close() error {
	err := cc.consumer.Close()
	if err != nil {
		return err
	}
	cc.consumer = nil
	return err
}

func (cc *saramaConsumerCommandImpl) PullMessage(consumerHandler ConsumerHandler) {
	chmsg := cc.consumer.Messages()
	for {
		message := <-chmsg
		cc.consumer.MarkPartitionOffset(message.Topic, message.Partition, message.Offset, "")
		consumerMessage := ConsumerMessage{Message: string(message.Value), Offset: message.Offset, Partition: message.Partition, Topic: message.Topic, ReceivedDateTimeUTC: time.Now().UTC()}
		//TODO: Pulled messages may be logged with DEBUG level once logging pattern in implemented in common-lib
		go consumerHandler(consumerMessage)
	}
}
func (cc *saramaConsumerCommandImpl) LimitedPullMessageNoOffset(consumerHandler ConsumerHandler, limiter Limiter) {

	for {
		if limiter.IsConsumingAllowed() {
			message := <-cc.consumer.Messages()
			consumerMessage := ConsumerMessage{
				Message:             string(message.Value),
				Offset:              message.Offset,
				Partition:           message.Partition,
				Topic:               message.Topic,
				ReceivedDateTimeUTC: time.Now().UTC(),
			}
			consumerHandler(consumerMessage)
		} else {
			limiter.Wait()
		}
	}
}

func (cc *saramaConsumerCommandImpl) MarkOffset(topic string, partition int32, offset int64) {
	cc.consumer.MarkPartitionOffset(topic, partition, offset, "")
}
