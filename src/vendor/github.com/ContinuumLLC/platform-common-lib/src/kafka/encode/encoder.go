package encode

import "github.com/Shopify/sarama"

//Encoder base interface with basic functions for encoding kafka message
type Encoder interface {
	Encode() ([]byte, error)
	Length() int
}

//GetBytesEncoder : Encodes a message for Kafka into byte Encoder
func GetBytesEncoder(message []byte) Encoder {
	return sarama.ByteEncoder(message)
}

//GetStringEncoder : Encodes a message for Kafka into byte Encoder
func GetStringEncoder(message string) Encoder {
	return sarama.StringEncoder(message)
}
