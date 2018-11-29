package messaging

import (
	"sync"

	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/json"
	"github.com/ContinuumLLC/platform-common-lib/src/kafka"
)

var service Service
var once sync.Once

type serviceImpl struct {
	conf            Config
	producer        kafka.ProducerService
	consumer        kafka.ConsumerService
	deserializer    json.DeserializerJSON
	serializer      json.SerializerJSON
	producerFactory kafka.ProducerFactory
	consumerFactory kafka.ConsumerFactory
}

//NewService returns new instance of NewService with the provided configuration
func NewService(conf Config) Service {
	once.Do(func() {
		service = &serviceImpl{
			conf:            conf,
			deserializer:    json.FactoryJSONImpl{}.GetDeserializerJSON(),
			serializer:      json.FactoryJSONImpl{}.GetSerializerJSON(),
			producerFactory: kafka.ProducerFactoryImpl{},
			consumerFactory: kafka.ConsumerFactoryImpl{},
		}
	})
	return service
}

func (s *serviceImpl) Publish(env *Envelope) error {
	if s.producer == nil {
		service, err := s.producerFactory.GetProducerService(s.producerConfig())
		if err != nil {
			return exception.New(ServiceCreationFailed, err)
		}
		s.producer = service
	}
	message, err := s.serializer.WriteByteStream(env)
	if err != nil {
		return exception.New(InvalidMessage, err)
	}
	return s.producer.Push(env.Topic, string(message))
}

func (s *serviceImpl) Listen(h ListenHandler) error {
	consumer, err := s.consumerFactory.GetConsumerService(s.consumerConfig())
	if err != nil {
		return exception.New(ServiceCreationFailed, err)
	}

	listener := handler{
		listener:     h,
		deserializer: s.deserializer,
	}
	s.consumer = consumer
	return consumer.PullHandler(listener.handle)
}

//TODO : Remove this after Proper implementation of Kafka Consumer and Producer
type handler struct {
	listener     ListenHandler
	deserializer json.DeserializerJSON
}

func (h handler) handle(m kafka.ConsumerMessage) {
	env := Envelope{}
	err := h.deserializer.ReadString(&env, m.Message)
	msg := Message{
		Envelope:            env,
		Err:                 err,
		Topic:               m.Topic,
		ReceivedDateTimeUTC: m.ReceivedDateTimeUTC,
		Offset:              m.Offset,
		Partition:           m.Partition,
	}
	h.listener(&msg)
}

func (s *serviceImpl) producerConfig() kafka.ProducerConfig {
	return kafka.ProducerConfig{
		ClientConfig: kafka.ClientConfig{
			BrokerAddress: s.conf.Address,
		},
	}
}

func (s *serviceImpl) consumerConfig() kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		ClientConfig: kafka.ClientConfig{BrokerAddress: s.conf.Address},
		GroupID:      s.conf.GroupID,
		Topics:       s.conf.Topics,
	}
}

func (s *serviceImpl) ListenWithLimiter(h ListenHandler, limiter kafka.Limiter) error {
	listener := handler{
		listener:     h,
		deserializer: s.deserializer,
	}
	return s.consumer.PullHandlerWithLimiter(listener.handle, limiter)
}

func (s *serviceImpl) Connect(inOut *kafka.ConsumerKafkaInOutParams) error {
	var err error
	s.consumer, err = s.consumerFactory.GetConsumerService(s.consumerConfig())
	if err != nil {
		return exception.New(ServiceCreationFailed, err)
	}

	return s.consumer.Connect(inOut)
}

func (s *serviceImpl) MarkOffset(pp PartitionParams) {
	s.consumer.MarkOffset(pp.Topic, pp.Partition, pp.Offset)
}
