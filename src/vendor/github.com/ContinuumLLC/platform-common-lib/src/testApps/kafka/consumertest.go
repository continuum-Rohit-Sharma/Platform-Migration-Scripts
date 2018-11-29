package main

import (
	"log"

	"github.com/ContinuumLLC/platform-common-lib/src/kafka"
)

func testNewSaramaConsumer() bool {
	config := new(kafka.ConsumerConfig)
	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"

	factory := new(kafka.ConsumerFactoryImpl)
	consumer, _ := factory.GetConsumerService(*config)
	if consumer == nil {
		return false
	}
	return true
}

func testNewSaramaConsumerNilBrokerAddress() bool {
	config := new(kafka.ConsumerConfig)

	config.BrokerAddress = nil
	config.GroupID = "testing1"

	factory := new(kafka.ConsumerFactoryImpl)
	_, err := factory.GetConsumerService(*config)

	if err == nil {
		return false
	}
	log.Println(err)
	return true
}

func testNewSaramaConsumerEmptyBrokerAddress() bool {
	config := new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{}
	config.GroupID = "testing1"

	factory := new(kafka.ConsumerFactoryImpl)
	_, err := factory.GetConsumerService(*config)

	if err == nil {
		return false
	}
	log.Println(err)
	return true
}

func testNewSaramaConsumerEmptyGroupID() bool {
	config := new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = ""

	factory := new(kafka.ConsumerFactoryImpl)
	_, err := factory.GetConsumerService(*config)

	if err == nil {
		return false
	}
	log.Println(err)
	return true
}

func testNewSaramaConsumerNilTopics() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	config.Topics = nil

	factory := new(kafka.ConsumerFactoryImpl)
	_, err := factory.GetConsumerService(*config)

	if err == nil {
		return false
	}
	log.Println(err)
	return true
}

func testNewSaramaConsumerEmptyTopics() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	config.Topics = []string{}

	factory := new(kafka.ConsumerFactoryImpl)
	_, err := factory.GetConsumerService(*config)

	if err == nil {
		return false
	}
	log.Println(err)
	return true
}

/*func TestSaramaConsumerConnect() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	consumer, _ := kafka.GetConsumer(config)
	err := consumer.Connect([]string{"Agentheartbeat"})
	if err != nil {
		return false
	}
	return true
}

func TestSaramaConsumerIsConnected() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	consumer, _ := kafka.GetConsumer(config)
	consumer.Connect([]string{"Agentheartbeat"})

	result := consumer.IsConnected()
	if result == false {
		return false
	}
	return true
}*/

/*func TestSaramaConsumerPull() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	config.Topics = []string{"test201"}

	command := new(kafka.ConsumerCommand)
	factory := new(kafka.ConsumerFactoryImpl)
	consumer, _ := factory.GetConsumerService(*config, *command)

	//consumer.Connect([]string{"test201"})

	//consumer.IsConnected()
	message, err := consumer.Pull()
	fmt.Println("Read message: " + message)
	consumer.CloseConnection()
	if err != nil {
		return false
	}
	return true
}*/

func testSaramaConsumerPullHandler() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	config.GroupID = "testing1"
	config.Topics = []string{"test201"}

	factory := new(kafka.ConsumerFactoryImpl)
	consumer, _ := factory.GetConsumerService(*config)

	_ = consumer.PullHandler(testCallBack)
	return true
}

func testSaramaConsumerPullHandlerLoop() bool {
	var config *kafka.ConsumerConfig
	config = new(kafka.ConsumerConfig)

	config.BrokerAddress = []string{"localhost:9092", "localhost:9093"}
	config.GroupID = "testing1"
	config.Topics = []string{"test301"}

	factory := new(kafka.ConsumerFactoryImpl)
	consumer, _ := factory.GetConsumerService(*config)

	_ = consumer.PullHandler(testCallBack)
	return true
}
func testCallBack(consumerMessage kafka.ConsumerMessage) {
	log.Printf("Received %v \n", consumerMessage)
}
