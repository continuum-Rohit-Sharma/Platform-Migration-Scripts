//Package kafka Structures to store kafka client configuration details
package kafka

import "errors"

//ClientConfig base properties for configs
type ClientConfig struct {
	BrokerAddress []string
}

//ProducerConfig contains configuration information for kafka consumer
type ProducerConfig struct {
	ClientConfig
	//RequiredAcks  int16
}

//ConsumerConfig contaijns configuration information for kafka consumer
type ConsumerConfig struct {
	ClientConfig
	GroupID string
	Topics  []string
}

//Validates Config for both producer and consumer
func validateConfig(config *ClientConfig) error {
	if config.BrokerAddress == nil {
		err := errors.New(ErrorBrokerAddressNotProvided)
		return err
	} else if len(config.BrokerAddress) == 0 {
		err := errors.New(ErrorBrokerAddressNotProvided)
		return err
	}
	return nil
}
