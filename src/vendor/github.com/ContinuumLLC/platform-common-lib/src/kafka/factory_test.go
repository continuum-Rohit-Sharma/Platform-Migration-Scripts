package kafka

import "testing"

func TestGetProducer(t *testing.T) {
	config := new(ProducerConfig)
	config.BrokerAddress = []string{"localhost:9092"}
	producerFactory := new(ProducerFactoryImpl)
	producer, _ := producerFactory.GetProducerService(*config)

	if producer == nil {
		t.Log("producer value is nil")
		t.Fail()
	}
}

func TestGetProducer_NilBrokerAddress(t *testing.T) {
	var config *ProducerConfig
	config = new(ProducerConfig)

	config.BrokerAddress = nil
	producerFactory := new(ProducerFactoryImpl)
	_, err := producerFactory.GetProducerService(*config)

	if err == nil || err.Error() != ErrorBrokerAddressNotProvided {
		t.Log("producer created with nil brokerAddress")
		t.Fail()
	}
}

func TestGetProducer_EmptyBrokerAddress(t *testing.T) {
	var config *ProducerConfig
	config = new(ProducerConfig)

	config.BrokerAddress = []string{}

	producerFactory := new(ProducerFactoryImpl)
	_, err := producerFactory.GetProducerService(*config)
	if err == nil || err.Error() != ErrorBrokerAddressNotProvided {
		t.Log("producer created with empty brokerAddress")
		t.Fail()
	}
}

func TestGetConsumer(t *testing.T) {

	config := new(ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9092"}

	consumerFactory := new(ConsumerFactoryImpl)
	_, err := consumerFactory.GetConsumerService(*config)

	if err == nil || err.Error() != ErrorClientGroupIDNotProvided {
		t.Log("consumer value is nil")
		t.Fail()
	}
}

func TestGetConsumer_NilBrokerAddress(t *testing.T) {

	var config *ConsumerConfig
	config = new(ConsumerConfig)

	config.BrokerAddress = nil
	config.GroupID = "test101"
	consumerFactory := new(ConsumerFactoryImpl)
	_, err := consumerFactory.GetConsumerService(*config)
	if err == nil || err.Error() != ErrorBrokerAddressNotProvided {
		t.Log("Consumer created with nil brokerAddress")
		t.Fail()
	}
}
func TestGetConsumer_EmptyBrokerAddress(t *testing.T) {

	var config *ConsumerConfig
	config = new(ConsumerConfig)

	config.BrokerAddress = []string{}
	config.GroupID = "test101"
	consumerFactory := new(ConsumerFactoryImpl)
	_, err := consumerFactory.GetConsumerService(*config)

	if err == nil || err.Error() != ErrorBrokerAddressNotProvided {
		t.Log("Consumer created with empty brokerAddress")
		t.Fail()
	}
}

func TestGetConsumer_EmptyGroupID(t *testing.T) {

	var config *ConsumerConfig
	config = new(ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = ""
	consumerFactory := new(ConsumerFactoryImpl)
	_, err := consumerFactory.GetConsumerService(*config)

	if err == nil || err.Error() != ErrorClientGroupIDNotProvided {
		t.Log("Consumer created with empty GroupID")
		t.Fail()
	}
}

func TestGetConsumer_EmptyTopics(t *testing.T) {

	var config *ConsumerConfig
	config = new(ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "grp1"

	consumerFactory := new(ConsumerFactoryImpl)
	_, err := consumerFactory.GetConsumerService(*config)

	if err == nil || err.Error() != ErrorTopicsNotProvided {
		t.Log("Consumer created with empty GroupID")
		t.Fail()
	}
}
func TestCommandFactory(t *testing.T) {
	cmdFactory := new(ConsumerCommandFactoryImpl)
	cmdService := cmdFactory.GetConsumerCommandService()
	if cmdService == nil {
		t.Fail()
	}
}
