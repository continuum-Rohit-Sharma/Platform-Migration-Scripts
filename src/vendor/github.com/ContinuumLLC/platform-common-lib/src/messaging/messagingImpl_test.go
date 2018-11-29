package messaging

import (
	"errors"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/json"
	jMock "github.com/ContinuumLLC/platform-common-lib/src/json/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/kafka"
	kMock "github.com/ContinuumLLC/platform-common-lib/src/kafka/mock"
	"github.com/golang/mock/gomock"
)

func Test_serviceImpl_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	producerFactory := kMock.NewMockProducerFactory(ctrl)
	consumerFactory := kMock.NewMockConsumerFactory(ctrl)
	serializer := jMock.NewMockSerializerJSON(ctrl)
	deserializer := jMock.NewMockDeserializerJSON(ctrl)

	p1 := kafka.ProducerConfig{
		ClientConfig: kafka.ClientConfig{
			BrokerAddress: []string{"1"},
		},
	}
	producerFactory.EXPECT().GetProducerService(p1).Return(nil, errors.New("Error"))

	p2 := kafka.ProducerConfig{
		ClientConfig: kafka.ClientConfig{
			BrokerAddress: []string{"2"},
		},
	}
	producer := kMock.NewMockProducerService(ctrl)
	producerFactory.EXPECT().GetProducerService(p2).Return(producer, nil)
	serializer.EXPECT().WriteByteStream(&Envelope{Topic: "2"}).Return(nil, errors.New("Error"))
	serializer.EXPECT().WriteByteStream(&Envelope{Topic: "3"}).Return([]byte{}, nil)
	producer.EXPECT().Push("3", gomock.Any()).Return(errors.New("Error"))
	serializer.EXPECT().WriteByteStream(&Envelope{Topic: "4"}).Return([]byte{}, nil)
	producer.EXPECT().Push("4", gomock.Any()).Return(nil)

	type fields struct {
		conf            Config
		producer        kafka.ProducerService
		deserializer    json.DeserializerJSON
		serializer      json.SerializerJSON
		producerFactory kafka.ProducerFactory
		consumerFactory kafka.ConsumerFactory
	}
	type args struct {
		env *Envelope
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			fields:  fields{conf: Config{Address: []string{"1"}, GroupID: "1", Topics: []string{"1"}}, producerFactory: producerFactory, serializer: serializer, deserializer: deserializer, consumerFactory: consumerFactory},
			args:    args{env: &Envelope{}},
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{conf: Config{Address: []string{"2"}, GroupID: "2", Topics: []string{"2"}}, producerFactory: producerFactory, serializer: serializer, deserializer: deserializer, consumerFactory: consumerFactory},
			args:    args{env: &Envelope{Topic: "2"}},
			wantErr: true,
		},
		{
			name:    "3",
			fields:  fields{conf: Config{}, producer: producer, producerFactory: producerFactory, serializer: serializer, deserializer: deserializer, consumerFactory: consumerFactory},
			args:    args{env: &Envelope{Topic: "3"}},
			wantErr: true,
		},
		{
			name:    "4",
			fields:  fields{conf: Config{}, producer: producer, producerFactory: producerFactory, serializer: serializer, deserializer: deserializer, consumerFactory: consumerFactory},
			args:    args{env: &Envelope{Topic: "4"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				conf:            tt.fields.conf,
				producer:        tt.fields.producer,
				deserializer:    tt.fields.deserializer,
				serializer:      tt.fields.serializer,
				producerFactory: tt.fields.producerFactory,
				consumerFactory: tt.fields.consumerFactory,
			}
			if err := s.Publish(tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("serviceImpl.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceImpl_Listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	producerFactory := kMock.NewMockProducerFactory(ctrl)
	consumerFactory := kMock.NewMockConsumerFactory(ctrl)
	serializer := jMock.NewMockSerializerJSON(ctrl)
	deserializer := jMock.NewMockDeserializerJSON(ctrl)

	c1 := kafka.ConsumerConfig{
		ClientConfig: kafka.ClientConfig{BrokerAddress: []string{"1"}},
		GroupID:      "1",
		Topics:       []string{"1"},
	}
	consumerFactory.EXPECT().GetConsumerService(c1).Return(nil, errors.New("Error"))

	c2 := kafka.ConsumerConfig{
		ClientConfig: kafka.ClientConfig{BrokerAddress: []string{"2"}},
		GroupID:      "2",
		Topics:       []string{"2"},
	}
	consumer := kMock.NewMockConsumerService(ctrl)
	consumerFactory.EXPECT().GetConsumerService(c2).Return(consumer, nil)
	consumer.EXPECT().PullHandler(gomock.Any()).Return(nil)

	type fields struct {
		conf            Config
		producer        kafka.ProducerService
		consumer        kafka.ConsumerService
		deserializer    json.DeserializerJSON
		serializer      json.SerializerJSON
		producerFactory kafka.ProducerFactory
		consumerFactory kafka.ConsumerFactory
	}
	type args struct {
		h ListenHandler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			fields:  fields{
				conf: Config{Address: []string{"1"},
					GroupID: "1", Topics: []string{"1"},
				},
				producerFactory: producerFactory,
				serializer: serializer,
				deserializer: deserializer,
				consumerFactory: consumerFactory},
			args:    args{h: func(m *Message) {}},
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{
				conf: Config{Address: []string{"2"}, GroupID: "2",
					Topics: []string{"2"},
				},
				producerFactory: producerFactory,
				serializer: serializer,
				deserializer: deserializer,
				consumerFactory: consumerFactory},
			args:    args{h: func(m *Message) {}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				conf:            tt.fields.conf,
				producer:        tt.fields.producer,
				consumer:		 tt.fields.consumer,
				deserializer:    tt.fields.deserializer,
				serializer:      tt.fields.serializer,
				producerFactory: tt.fields.producerFactory,
				consumerFactory: tt.fields.consumerFactory,
			}
			if err := s.Listen(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("serviceImpl.Listen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceImpl_ListenWithLimiter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	limiter := kMock.NewMockLimiter(ctrl)

	consumer := kMock.NewMockConsumerService(ctrl)
	consumer.EXPECT().PullHandlerWithLimiter(gomock.Any(), limiter).Return(nil)

	type fields struct {
		conf            Config
		consumer        kafka.ConsumerService
	}
	type args struct {
		handler ListenHandler
		limiter kafka.Limiter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "4",
			fields:  fields{
				conf: Config{Address: []string{"4"},
					GroupID: "4",
					Topics: []string{"4"},
				},
				consumer:consumer},
			args:    args{handler: func(m *Message) {}, limiter: limiter},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				conf:            tt.fields.conf,
				consumer:        tt.fields.consumer,
			}
			var err error
			if err = s.ListenWithLimiter(tt.args.handler, tt.args.limiter);(err != nil) != tt.wantErr {
				t.Errorf("serviceImpl.ListenWithLimiter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceImpl_MarkOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := kMock.NewMockConsumerService(ctrl)
	consumer.EXPECT().MarkOffset(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	s := serviceImpl{
		consumer:        consumer,
	}

	s.MarkOffset(PartitionParams{Topic:"topic", Partition:0, Offset:0})

}
func Test_serviceImpl_Connect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumerFactory := kMock.NewMockConsumerFactory(ctrl)

	c1 := kafka.ConsumerConfig{
		ClientConfig: kafka.ClientConfig{BrokerAddress: []string{"3"}},
		GroupID:      "3",
		Topics:       []string{"3"},
	}
	consumerFactory.EXPECT().GetConsumerService(c1).Return(nil, errors.New("Error"))

	c2 := kafka.ConsumerConfig{
		ClientConfig: kafka.ClientConfig{BrokerAddress: []string{"4"}},
		GroupID:      "4",
		Topics:       []string{"4"},
	}
	consumer := kMock.NewMockConsumerService(ctrl)
	consumerFactory.EXPECT().GetConsumerService(c2).Return(consumer, nil)
	consumer.EXPECT().Connect(gomock.Any()).Return(nil)

	type fields struct {
		conf            Config
		consumerFactory kafka.ConsumerFactory
	}
	type args struct {
		a *kafka.ConsumerKafkaInOutParams
	}

	inOut := &kafka.ConsumerKafkaInOutParams {
		Errors: nil,
		Notifications: nil,
		ReturnErrors: false,
		ReturnNotifications: false,
		OffsetsInitial: 0,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "3",
			fields:  fields{
				conf: Config{Address: []string{"3"}, GroupID: "3",
					Topics: []string{"3"}},
				consumerFactory: consumerFactory},
			args:    args{a:inOut},
			wantErr: true,
		},
		{
			name:    "4",
			fields:  fields{
				conf: Config{Address: []string{"4"}, GroupID: "4",
					Topics: []string{"4"}},
				consumerFactory: consumerFactory},
			args:    args{a:inOut},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				conf:            tt.fields.conf,
				consumerFactory: tt.fields.consumerFactory,
			}
			var err error
			if err = s.Connect(tt.args.a);(err != nil) != tt.wantErr {
				t.Errorf("s.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	service := NewService(Config{})
	_, ok := service.(*serviceImpl)
	if !ok {
		t.Error("Invalid serviceImpl")
	}
}
