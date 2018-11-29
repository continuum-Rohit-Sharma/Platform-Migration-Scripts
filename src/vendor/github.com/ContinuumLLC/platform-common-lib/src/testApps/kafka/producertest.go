package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/kafka"
)

func testNewSaramaProducer() bool {
	config := new(kafka.ProducerConfig)

	config.BrokerAddress = []string{"localhost:9092"}
	factory := new(kafka.ProducerFactoryImpl)
	producer, _ := factory.GetProducerService(*config)
	producer.CloseConnection()
	if producer == nil {
		return false
	}
	return true
}

func testNewSaramaProducerNilBrokerAddress() bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = nil

	factory := new(kafka.ProducerFactoryImpl)
	_, err := factory.GetProducerService(*config)

	if err == nil {
		return false
	}
	fmt.Println(err)
	return true
}
func testNewSaramaProducerEmptyBrokerAddress() bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = []string{}

	factory := new(kafka.ProducerFactoryImpl)
	_, err := factory.GetProducerService(*config)

	if err == nil {
		return false
	}
	fmt.Println(err)
	return true
}

/*func TestSaramaProducerConnect() bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = []string{"localhost:9092"}

	producer, _ := kafka.GetProducer(config)
	err := producer.Connect()
	if err != nil {
		return false
	}
	return true
}

func TestSaramaProducerIsConnected() bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = []string{"localhost:9092"}

	producer, _ := kafka.GetProducer(config)
	producer.Connect()

	result := producer.IsConnected()
	if result == false {
		return false
	}
	return true
}*/

func testSaramaProducerPush() bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = []string{"localhost:9092"}

	factory := new(kafka.ProducerFactoryImpl)
	producer, _ := factory.GetProducerService(*config)

	err := producer.Push("test201", "New message:"+time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))
	producer.CloseConnection()
	if err != nil {
		return false
	}
	return true
}

func testSaramaProducerPushLoop(message string) bool {
	var config *kafka.ProducerConfig
	config = new(kafka.ProducerConfig)

	config.BrokerAddress = []string{"localhost:9092", "localhost:9093"}

	factory := new(kafka.ProducerFactoryImpl)
	producer, _ := factory.GetProducerService(*config)

	i := 0

	for {
		i++
		log.Println(fmt.Sprintf("sending "+message+"-%d", i))
		err := producer.Push("test301", fmt.Sprintf(message+"%d", i))
		if err != nil {
			log.Println("-------------------------------Producing Error-------------------------------")
			log.Println(err)
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	producer.CloseConnection()
	return true
}
