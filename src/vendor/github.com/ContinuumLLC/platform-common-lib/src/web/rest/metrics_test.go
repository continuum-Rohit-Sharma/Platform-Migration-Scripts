package rest

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	metricKey      = "test"
	metricValInt   = 10
	metricValFloat = 10.0
)

var errSomethingWentWrong = errors.New("something went wrong")

type testMetrics struct {
	Metrics
}

func (m testMetrics) Update() error {
	return errSomethingWentWrong
}

func TestRegistryMetrics(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{})
	if storage.Metrics.(*Metrics).Type != metricsType {
		t.Errorf("expected Type %s, but got %s", metricsType, storage.Metrics.(*Metrics).Type)
	}
}

func TestHandlerMetrics_OK(t *testing.T) {
	oldstorage := storage.Metrics
	storage.Metrics = &Metrics{}
	defer func() { storage.Metrics = oldstorage }()

	mock := &mockResponseWriter{dataHeader: http.Header{}}
	HandlerMetrics(mock, nil)
	if mock.dataWriteHeader != http.StatusOK {
		t.Errorf("expected code %d, but got %d", http.StatusOK, mock.dataWriteHeader)
	}
}

func TestHandlerMetrics_StatusInternalServerError(t *testing.T) {
	oldstorage := storage.Metrics
	storage.Metrics = &testMetrics{}
	defer func() { storage.Metrics = oldstorage }()

	mock := &mockResponseWriter{dataHeader: http.Header{}}
	HandlerMetrics(mock, nil)
	if mock.dataWriteHeader != http.StatusInternalServerError {
		t.Errorf("expected code %d, but got %d", http.StatusInternalServerError, mock.dataWriteHeader)
	}
}

func TestSet(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{
		{Name: metricKey, DataType: DataTypeInteger},
	})

	if err := Add(metricKey, metricValInt); err != nil {
		t.Error(err)
	}

	if err := Set(metricKey, metricValInt); err != nil {
		t.Error(err)
	}

	if storage.Data[0].Value != metricValInt {
		t.Errorf("expected value %v, got %v", storage.Data[0].Value, metricValInt)
	}
}

func TestAdd(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{
		{Name: metricKey, DataType: DataTypeDecimal},
	})
	const count = 3
	for i := 0; i < count; i++ {
		if err := Add(metricKey, metricValFloat); err != nil {
			t.Error(err)
		}
	}

	if storage.Data[0].Value != (metricValFloat * count) {
		t.Errorf("expected value %v, got %v", metricValFloat*count, storage.Data[0].Value)
	}
}

func TestClear(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{
		{Name: metricKey, DataType: DataTypeInteger},
	})

	if err := Set(metricKey, metricValInt); err != nil {
		t.Error(err)
	}

	if storage.Data[0].Value != metricValInt {
		t.Errorf("expected value %v, got %v", storage.Data[0].Value, metricValInt)
	}

	storage.clear()

	if storage.Data[0].Value != nilValue(DataTypeInteger) {
		t.Errorf("expected value %v, got %v", metricValInt, storage.Data[0].Value)
	}
}

func TestNilValue(t *testing.T) {
	if val := nilValue(DataTypeZero); val != nil {
		t.Error("value should be <nil>")
	}
	if val := nilValue(DataTypeInteger); val != 0 {
		t.Error("value should be 0")
	}
	if val := nilValue(DataTypeDecimal); val != 0.0 {
		t.Error("value should be 0.0")
	}
}

func TestToDefault(t *testing.T) {
	cfg := []MetricsConfig{{Name: metricKey}}
	toDefault(cfg)
	if cfg[0].Range != RangeTypeInfinity {
		t.Errorf("Range should be %s", RangeTypeInfinity)
	}
	if cfg[0].Unit != UnitTypeNumbers {
		t.Errorf("Unit should be %s", UnitTypeNumbers)
	}
	if cfg[0].DataType != DataTypeInteger {
		t.Errorf("Unit should be %s", DataTypeInteger)
	}
}

func TestRangeValidete(t *testing.T) {
	testCases := []struct {
		result    bool
		value     interface{}
		rangeType rangeType
	}{
		{result: true, value: 0, rangeType: RangeTypeInfinity},
		{result: true, value: -1, rangeType: RangeTypeInfinity},
		{result: true, value: 1, rangeType: RangeTypeInfinity},
		{result: true, value: -10, rangeType: RangeTypeFromInfinity},
		{result: true, value: 10, rangeType: RangeTypeToInfinity},
		{result: false, value: -10, rangeType: RangeTypeToInfinity},
		{result: false, value: 10, rangeType: RangeTypeFromInfinity},
		{result: true, value: 0.0, rangeType: RangeTypeInfinity},
		{result: true, value: -1.0, rangeType: RangeTypeInfinity},
		{result: true, value: 1.0, rangeType: RangeTypeInfinity},
		{result: true, value: -10.0, rangeType: RangeTypeFromInfinity},
		{result: true, value: 10.0, rangeType: RangeTypeToInfinity},
		{result: false, value: -10.0, rangeType: RangeTypeToInfinity},
		{result: false, value: 10.0, rangeType: RangeTypeFromInfinity},
		{result: false, value: "10.0", rangeType: RangeTypeInfinity},
	}

	for _, test := range testCases {
		result := rangeValidate(test.value, test.rangeType)
		if test.result != result {
			t.Errorf("expected %t,got %t", test.result, result)
		}
	}
}

func TestDoDecimal(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	testCases := []struct {
		result    error
		value     interface{}
		incr      bool
		dataType  dataType
		rangeType rangeType
	}{
		{result: ErrMetricBadFloatType, value: "1.0"},
		{result: errors.New(fmt.Sprintf(valueShouldBeInRange, RangeTypeFromInfinity)), value: 1.0, rangeType: RangeTypeFromInfinity},
		{result: errors.New(fmt.Sprintf(valueShouldBeInRange, RangeTypeToInfinity)), value: -1.0, rangeType: RangeTypeToInfinity},
		{result: nil, value: 1.0, rangeType: RangeTypeToInfinity, incr: true, dataType: DataTypeDecimal},
	}

	for _, test := range testCases {
		RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{
			{Name: metricKey, DataType: test.dataType, Range: test.rangeType},
		})

		result := doDecimal(0, test.value, test.incr)
		if !reflect.DeepEqual(test.result, result) {
			t.Errorf("expected %v, got %v", test.result, result)
		}
	}
}

func TestDoInteger(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	testCases := []struct {
		result    error
		value     interface{}
		incr      bool
		dataType  dataType
		rangeType rangeType
	}{
		{result: ErrMetricBadIntType, value: "1"},
		{result: errors.New(fmt.Sprintf(valueShouldBeInRange, RangeTypeFromInfinity)), value: 1, rangeType: RangeTypeFromInfinity},
		{result: errors.New(fmt.Sprintf(valueShouldBeInRange, RangeTypeToInfinity)), value: -1, rangeType: RangeTypeToInfinity},
		{result: nil, value: 1, rangeType: RangeTypeToInfinity, incr: true, dataType: DataTypeInteger},
	}

	for _, test := range testCases {
		RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{
			{Name: metricKey, DataType: test.dataType, Range: test.rangeType},
		})

		result := doInteger(0, test.value, test.incr)
		if !reflect.DeepEqual(test.result, result) {
			t.Errorf("expected %v, got %v", test.result, result)
		}
	}
}

func TestSave(t *testing.T) {
	oldstorage := storage.Metrics
	defer func() { storage.Metrics = oldstorage }()

	if err := save("", 0, false); err != ErrEmptyKey {
		t.Error(err)
	}

	RegistryMetrics(&Metrics{Type: metricsType}, []MetricsConfig{})

	if err := save(metricKey, 0, false); err != nil {
		t.Error(err)
	}
}
