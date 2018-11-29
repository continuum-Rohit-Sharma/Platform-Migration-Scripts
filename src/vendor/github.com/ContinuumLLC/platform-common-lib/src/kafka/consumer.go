//Package kafka saramaConsumer uses the sarama kafka Consumer Client for Go
package kafka

import "errors"

type saramaConsumerServiceImpl struct {
	config  *ConsumerConfig
	command ConsumerCommand
}

//newSaramaConsumer returns new instance of saramaConsumer with the provided configuration
func newSaramaConsumer(config *ConsumerConfig, cmdFactory ConsumerCommandFactory) (*saramaConsumerServiceImpl, error) {
	err := validateClientConfig(config)
	if err != nil {
		return nil, err
	}

	consumer := new(saramaConsumerServiceImpl)
	consumer.config = config
	consumer.command = cmdFactory.GetConsumerCommandService()
	return consumer, nil
}

func validateClientConfig(config *ConsumerConfig) error {
	err := validateConfig(&config.ClientConfig)
	if err != nil {
		return err
	}
	if config.GroupID == "" {
		err := errors.New(ErrorClientGroupIDNotProvided)
		return err
	} else if len(config.Topics) == 0 {
		err := errors.New(ErrorTopicsNotProvided)
		return err
	}
	return nil
}

//Connect connects the consumer kafka broker
func (sc *saramaConsumerServiceImpl) connect() error {
	var result bool

	result = sc.isConnected()
	if result == true {
		return nil
	}
	err := sc.command.NewConsumer(sc.config.BrokerAddress, sc.config.GroupID, sc.config.Topics)

	if err != nil {
		return err
	}
	return nil
}

//IsConnected checks is kafka server is connected
func (sc saramaConsumerServiceImpl) isConnected() bool {
	return sc.command.IsConnected()
}

//PullHandler pulls message and  one message from kafka queue and returns it
func (sc saramaConsumerServiceImpl) PullHandler(consumerHandler ConsumerHandler) error {
	var err error
	err = sc.connect()

	if err != nil {
		return err
	}

	result := sc.isConnected()
	if result == false {
		return errors.New(ErrorClientNotConnected)
	}

	sc.command.PullMessage(consumerHandler)
	return nil
}

//Connect method, to have a possibility of reconnecting at client
func (sc saramaConsumerServiceImpl) Connect(config *ConsumerKafkaInOutParams) error {
	if sc.isConnected() {
		return nil
	}

	err := sc.command.NewCustomConsumer(config, sc.config.BrokerAddress, sc.config.GroupID, sc.config.Topics)
	if err != nil {
		return err
	}

	if !sc.isConnected() {
		return errors.New(ErrorClientNotConnected)
	}

	return nil
}

//PullHandlerWithLimiter pulls message and  one message from kafka queue and returns it, speed is reduced with the limiter
func (sc saramaConsumerServiceImpl) PullHandlerWithLimiter(consumerHandler ConsumerHandler,
	limiter Limiter) error {

	if !sc.isConnected() {
		return errors.New(ErrorClientNotConnected)
	}
	sc.command.LimitedPullMessageNoOffset(consumerHandler, limiter)
	return nil
}

//MarkOffset method to perform mark offset
func (sc saramaConsumerServiceImpl) MarkOffset(topic string, partition int32, offset int64) {
	sc.command.MarkOffset(topic, partition, offset)
}
//CloseConnection closes connection to kafka server
func (sc *saramaConsumerServiceImpl) CloseConnection() error {
	var result bool

	result = sc.isConnected()
	if result == false {
		return nil
	}

	err := sc.command.Close()
	if err != nil {
		return err
	}
	return nil
}
