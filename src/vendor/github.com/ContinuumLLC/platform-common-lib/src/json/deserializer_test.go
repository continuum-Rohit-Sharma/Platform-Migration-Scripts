package json

import (
	"reflect"
	"strings"
	"testing"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"
)

type Kafkaconfig struct {
	KafkaAddress []string
}

var (
	nestedJSON = `{
				"Cluster" : 
					{
						"IP":"192.168.137.122",
						"CassandraPort" : "1234"
					},"Port":9042, "ProtoVersion":4,"Keyspace":"broker_keyspace"
				}`
	kafkaConfig = `{
					"kafkaAddress": ["localhost:9092","localhost:9093"]
					}`
	kafkaInvalidConfig = `{
						"kafkaAddress": "localhost:9092"
						}`
)

func TestDeserializeToMapNestedConfig(t *testing.T) {
	var mapConf map[string]interface{}
	Deserialize(&mapConf, strings.NewReader(nestedJSON))
	//Reading nested map
	clusterIP := mapConf["Cluster"].(map[string]interface{})["IP"]
	if clusterIP != "192.168.137.122" && mapConf["Keyspace"] != "broker_keyspace" {
		t.Error("Deserializing failed to map[string]interface{}")
	}
}

func TestDeserializeToStructValidConfig(t *testing.T) {
	var conf Kafkaconfig
	Deserialize(&conf, strings.NewReader(kafkaConfig))
	if len(conf.KafkaAddress) == 0 || (conf.KafkaAddress[0] != "localhost:9092" && conf.KafkaAddress[1] != "localhost:9093") {
		t.Error("JSON Config data is not valid to convert into the specified struct")
	}
}

func TestDeserializeToStructInvalidConfig(t *testing.T) {
	var conf Kafkaconfig
	err := Deserialize(&conf, strings.NewReader(kafkaInvalidConfig))
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONFailedToDeserialize {
			t.Errorf("Expected ErrJSONFailedToDeserialize but got %v", exce)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestDeserializeToStructInvalidStream(t *testing.T) {
	var conf Kafkaconfig
	err := Deserialize(&conf, strings.NewReader(""))
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONInvalidStream {
			t.Errorf("Expected ErrJSONInvalidStream but got %v", exce)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestDeserializeToStructInvalidInput(t *testing.T) {
	var conf Kafkaconfig
	//Deserialize is expecting a reference object to be passed for interface{}
	err := Deserialize(conf, strings.NewReader(kafkaConfig))
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONNotAPointerOrNil {
			t.Errorf("Expected ErrJSONNotAPointerOrNil but got %v", exce)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestDeserializeBytes(t *testing.T) {
	tests := []struct {
		name      string
		inputJson []byte
		want      interface{}
		wantErr   bool
	}{
		{
			"TestDeserializeBytes_1:Given vaild json input",
			[]byte(jsonString()),
			jsonObject(),
			false,
		},
		{
			"TestDeserializeBytes_2:Given invalid json input",
			[]byte(`{"partnerID": }`),
			nil,
			true,
		},
		{
			"TestDeserializeBytes_3:Given empty json input",
			[]byte{},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeserializeBytes(tt.inputJson)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func jsonString() string {
	return `{"partnerID": "50001794","endpointID":  "50109539"}`
}

func jsonObject() interface{} {
	return map[string]interface{}{
		"partnerID":  "50001794",
		"endpointID": "50109539",
	}
}

type MockObject struct {
	partnerID string `json:partnerID`
}
