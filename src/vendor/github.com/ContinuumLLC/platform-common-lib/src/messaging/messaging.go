package messaging

import (
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/kafka"
)

const (
	//InvalidMessage An error to be thrown in case of message serlization
	InvalidMessage = "InvalidMessage"

	//ServiceCreationFailed : An error to be thrown in case of Consumer/Producer Service Creation Failed
	ServiceCreationFailed = "ServiceCreationFailed"

	//MessageType : is a key to hold message Type
	MessageType = "continuum.message.type"
	//TransactionID for transaction key in header map
	TransactionID string = "continuum.message.transaction.id"
	//DataCenterTimeStamp for datetimestamp key in header map
	DataCenterTimeStamp string = "continuum.message.dc.timestamp"
	//SiteMigrationMessageType to be used in site migration data pushed on kafka
	SiteMigrationMessageType = "SITE-MIGRATION"
)

//Config is a struct to hold Messaging server configuration
type Config struct {
	Address []string
	GroupID string
	Topics  []string
}

// Envelope is the envelope for Kafka message to be plublished and Produced
type Envelope struct {
	Header  Header
	Topic   string
	Message string
}

//Header is a map for Request Response structures
type Header map[string][]string

//Set : Sets a key to a value (overwriting if it exists)
func (h Header) Set(key string, value string) {
	h[key] = []string{value}
}

//Get : returns values for given key
func (h Header) Get(key string) []string {
	return h[key]
}

//Remove : Removes a Header value for give key
func (h Header) Remove(key string) {
	delete(h, key)
}

//ListenHandler a message Handler after pulling data from kafka
type ListenHandler func(*Message)

//Message represents Message to be comsumed by Listener
type Message struct {
	Envelope            Envelope
	Err                 error
	Topic               string
	ReceivedDateTimeUTC time.Time
	Offset              int64
	Partition           int32
}

type PartitionParams struct {
	Topic     string
	Partition int32
	Offset    int64
}

//Service : is an service to Publish and consume message from Kafka
type Service interface {
	Publish(env *Envelope) error
	Listen(ListenHandler) error

	ListenWithLimiter(ListenHandler, kafka.Limiter) error
	Connect(*kafka.ConsumerKafkaInOutParams) error
	MarkOffset(PartitionParams)
}
