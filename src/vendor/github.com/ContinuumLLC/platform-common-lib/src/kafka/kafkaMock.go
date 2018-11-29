package kafka

//Mock is a mock struct of Kafka Consumer
type Mock struct {
	KMockData *MockData
	ConsumerServiceMock
}

//MockData is to hold the Consumer Service objects
type MockData struct {
	ConSvrcObj ConsumerService
	ConSvrcErr error
}

//ConsumerServiceMock is a mock struct of ConsumerService
type ConsumerServiceMock struct {
	KConSrvcData *ConsumerServiceMockData
}

//ConsumerServiceMockData is to hold the properties which is been returned by the methods of ConsumerService
type ConsumerServiceMockData struct {
	CloseConnObj   error
	PullMsgObj     string
	PullMsgErr     error
	PullHandlerErr error
}

//GetConsumerService is a mock method
func (m Mock) GetConsumerService(conf ConsumerConfig) (ConsumerService, error) {
	return m, m.KMockData.ConSvrcErr
}

//CloseConnection is a mock method
func (m ConsumerServiceMock) CloseConnection() error {
	return m.KConSrvcData.CloseConnObj
}

//Pull is a mock method
func (m ConsumerServiceMock) Pull() (string, error) {
	return m.KConSrvcData.PullMsgObj, m.KConSrvcData.PullMsgErr
}

//PullHandler is a mock method
func (m ConsumerServiceMock) PullHandler(ConsumerHandler) error {
	return m.KConSrvcData.PullHandlerErr
}

//PullHandlerWithLimiter
func (m ConsumerServiceMock) PullHandlerWithLimiter(ConsumerHandler, Limiter) error {
	return m.KConSrvcData.PullHandlerErr
}

//PullHandlerWithLimiter
func (m ConsumerServiceMock) Connect(*ConsumerKafkaInOutParams) error {
	return m.KConSrvcData.CloseConnObj
}

//MarkOffset
func (m ConsumerServiceMock) MarkOffset(t string, p int32, o int64)  {}