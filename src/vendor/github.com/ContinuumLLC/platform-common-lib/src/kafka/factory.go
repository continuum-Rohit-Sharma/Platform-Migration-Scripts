//Package kafka factory methods returns a specific implementation of kafka Client
package kafka

import "github.com/ContinuumLLC/platform-common-lib/src/web/rest"

//ProducerFactoryImpl implments ProducerFactory
type ProducerFactoryImpl struct {
}

//GetProducerService returns a ProducerService
func (ProducerFactoryImpl) GetProducerService(config ProducerConfig) (ProducerService, error) {
	cmdFactory := new(ProducerCommandFactoryImpl)
	return newSaramaProducer(&config, cmdFactory)
}

//GetConfluentProducerService gets confluent Producer
func (ProducerFactoryImpl) GetConfluentProducerService(config ProducerConfig) (ProducerService, error) {
	cmdFactory := new(ProducerCommandFactoryImpl)
	return newSaramaProducer(&config, cmdFactory)
}

//ConsumerFactoryImpl returns a ConsumerConfig
type ConsumerFactoryImpl struct {
}

//GetConsumerService return a ConsumerService
func (ConsumerFactoryImpl) GetConsumerService(config ConsumerConfig) (ConsumerService, error) {
	cmdFactory := new(ConsumerCommandFactoryImpl)
	return newSaramaConsumer(&config, cmdFactory)
}

//ProducerCommandFactoryImpl implements Factory mathod that gets ProducerCommandService
type ProducerCommandFactoryImpl struct {
}

//GetProducerCommandService returns implementation of ProducerCommand
func (ProducerCommandFactoryImpl) GetProducerCommandService() ProducerCommand {
	return new(saramaProducerCommandImpl)
}

//GetConfluentCommandService return confluent implementation of producer command
func (ProducerCommandFactoryImpl) GetConfluentProducerCommandService() ProducerCommand {
	return new(saramaProducerCommandImpl)
}

//ConsumerCommandFactoryImpl implements Factory mathod that gets ProducerCommandService
type ConsumerCommandFactoryImpl struct {
}

//GetConsumerCommandService returns implementation of ProducerCommand
func (ConsumerCommandFactoryImpl) GetConsumerCommandService() ConsumerCommand {
	return new(saramaConsumerCommandImpl)
}

//Health returns a Health state for Kafka
func (ProducerFactoryImpl) Health(kafkaBrokers []string) rest.Statuser {
	cmdFactory := new(ProducerCommandFactoryImpl)
	return status{
		kafkaBrokers: kafkaBrokers,
		factory:      cmdFactory,
	}
}

type status struct {
	kafkaBrokers []string
	factory      ProducerCommandFactory
}

func (k status) Status(conn rest.OutboundConnectionStatus) *rest.OutboundConnectionStatus {
	conn.ConnectionType = "Kafka"
	conn.ConnectionURLs = k.kafkaBrokers
	s := k.factory.GetProducerCommandService()
	err := s.NewProducer(k.kafkaBrokers)
	if err != nil {
		conn.ConnectionStatus = rest.ConnectionStatusUnavailable
	} else {
		conn.ConnectionStatus = rest.ConnectionStatusActive
	}
	return &conn
}
