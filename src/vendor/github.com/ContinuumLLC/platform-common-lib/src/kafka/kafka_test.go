package kafka

import (
	"errors"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/kafka/encode"
)

type mockKafka struct {
	NewProducerError error
	NewConsumerError error
	CloseError       error
	PushMessageError error
	PullMessageError error
	IsConnect        bool
	PullHandlerError error
}

func (m *mockKafka) GetProducerService(config ProducerConfig) (ProducerService, error) {
	return newSaramaProducer(&config, m)
}
func (m *mockKafka) GetConfluentProducerService(config ProducerConfig) (ProducerService, error) {
	return newSaramaProducer(&config, m)
}

func (m *mockKafka) GetConsumerService(config ConsumerConfig) (ConsumerService, error) {

	return newSaramaConsumer(&config, m)
}
func (m *mockKafka) GetProducerCommandService() ProducerCommand {
	return m
}
func (m *mockKafka) GetConfluentProducerCommandService() ProducerCommand {
	return m
}
func (m *mockKafka) GetConsumerCommandService() ConsumerCommand {
	return m
}
func (m *mockKafka) NewProducer(brokerAddress []string) error {
	return m.NewProducerError
}
func (m *mockKafka) NewConsumer(brokerAddress []string, GroupID string, Topics []string) error {
	return m.NewConsumerError
}
func (m *mockKafka) NewCustomConsumer(c *ConsumerKafkaInOutParams, brokerAddress []string, GroupID string, Topics []string) error {
	return m.NewConsumerError
}
func (m *mockKafka) Close() error {
	return m.CloseError
}
func (m *mockKafka) PushMessage(topicName string, message string) (int32, int64, error) {
	return 0, 0, m.PushMessageError
}
func (m *mockKafka) PushMessageEncoder(topicName string, message encode.Encoder) (int32, int64, error) {
	return 0, 0, m.PushMessageError
}
func (m *mockKafka) PullMessage(consumerHandler ConsumerHandler) {
}
func (m *mockKafka) LimitedPullMessageNoOffset(consumerHandler ConsumerHandler, limiter Limiter) {
}
func (m *mockKafka) MarkOffset(t string, p int32, o int64) {
}
func (m *mockKafka) IsConnected() bool {
	return m.IsConnect
}

type mockLimiter struct{}

func (m *mockLimiter) IsConsumingAllowed() bool {
	return true
}
func (m *mockLimiter) Wait() {}

func TestNewProducer(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	kafka.NewProducerError = nil
	kafka.IsConnect = true
	config.BrokerAddress = []string{"localhost:9002"}
	_, err := kafka.GetProducerService(*config)
	if err != nil {
		t.Fail()
	}
}
func TestNewProducerConnectError(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	kafka.NewProducerError = errors.New("New Producer Error")
	config.BrokerAddress = []string{"localhost:9002"}
	producerService, err := kafka.GetProducerService(*config)
	producerService.Push("", "")
	if err != nil {
		t.Fail()
	}
}
func TestNewProducerConnectAlreadyConnectedError(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	kafka.IsConnect = true
	config.BrokerAddress = []string{"localhost:9002"}
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.Push("", "")
	if err != nil {
		t.Fail()
	}
}
func TestProducerPush(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	kafka.PushMessageError = nil
	kafka.NewProducerError = nil
	kafka.IsConnect = true
	config.BrokerAddress = []string{"localhost:9002"}

	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.Push("testTopic", "data")
	if err != nil {
		t.Fail()
	}
}

func TestProducerIsConnectedError(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = false
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.Push("testTopic", "data")

	if err == nil || err.Error() != ErrorClientNotConnected {
		t.Fail()
	}
}

func TestProducerPushError(t *testing.T) {
	pushError := "Push Error"
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = true
	kafka.PushMessageError = errors.New(pushError)
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.Push("testTopic", "data")

	if err == nil || err.Error() != pushError {
		t.Fail()
	}
}

func TestProducerClose(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = true
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.CloseConnection()
	if err != nil {
		t.Fail()
	}
}

func TestProducerConnectedClose(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = true
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.CloseConnection()
	if err != nil {
		t.Fail()
	}
}

func TestProducerConnectionClosedClose(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = false
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.CloseConnection()
	if err != nil {
		t.Fail()
	}
}

func TestProducerCloseError(t *testing.T) {
	closingError := "Error Closing"
	kafka := new(mockKafka)
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	kafka.IsConnect = true
	kafka.CloseError = errors.New(closingError)
	producerService, _ := kafka.GetProducerService(*config)
	err := producerService.CloseConnection()
	if err == nil || err.Error() != closingError {
		t.Fail()
	}
}

func TestNewConsumer(t *testing.T) {
	kafka := new(mockKafka)
	kafka.NewConsumerError = nil
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}

	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandler(mockHandle)
	if err == nil || err.Error() != ErrorClientNotConnected {
		t.Fail()
	}
}

func TestNewConsumerError(t *testing.T) {
	newConsumerrError := "NewConsumerError"
	kafka := new(mockKafka)
	kafka.NewConsumerError = errors.New(newConsumerrError)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}

	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandler(mockHandle)
	if err == nil || err.Error() != newConsumerrError {
		t.Fail()
	}
}

func TestConsumerPull(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}
	kafka.NewConsumerError = nil
	kafka.IsConnect = true
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandler(mockHandle)
	if err != nil {
		t.Fail()
	}
}

func TestConsumerConnectionErrorPullError(t *testing.T) {
	pullError := "Pull Error"
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}
	kafka.NewConsumerError = errors.New(pullError)
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandler(mockHandle)
	if err == nil || err.Error() != pullError {
		t.Fail()
	}
}

func TestConsumerConnectionErrorPullWithLimiterError(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"2"}
	config.GroupID = "Grp1"
	config.Topics = []string{"t"}
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandlerWithLimiter(mockHandle, new(mockLimiter))

	if err == nil || err.Error() != ErrorClientNotConnected {
		t.Fail()
	}
	inOut := &ConsumerKafkaInOutParams{
		Errors:              nil,
		Notifications:       nil,
		ReturnErrors:        false,
		ReturnNotifications: false,
		OffsetsInitial:      0,
	}
	err = consumerService.Connect(inOut)

	if err == nil || err.Error() != ErrorClientNotConnected {
		t.Fail()
	}
	consumerService.MarkOffset("", 0, 0)
}

func TestConsumerIsConnectedPullError(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}
	kafka.IsConnect = false
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.PullHandler(mockHandle)
	if err == nil || err.Error() != ErrorClientNotConnected {
		t.Fail()
	}
}

func TestConsumerClose(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}

	kafka.NewConsumerError = nil
	kafka.IsConnect = true

	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.CloseConnection()
	if err != nil {
		t.Fail()
	}
}

func TestConsumerIsConnectClose(t *testing.T) {
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}
	kafka.IsConnect = false
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.CloseConnection()
	if err != nil {
		t.Fail()
	}
}

func TestConsumerCloseError(t *testing.T) {
	closeError := "Consumer Close Error"
	kafka := new(mockKafka)
	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9002"}
	config.GroupID = "Grp1"
	config.Topics = []string{""}
	kafka.IsConnect = true
	kafka.CloseError = errors.New(closeError)
	consumerService, _ := kafka.GetConsumerService(*config)
	err := consumerService.CloseConnection()
	if err == nil || err.Error() != closeError {
		t.Fail()
	}
}

func mockHandle(message ConsumerMessage) {

}
