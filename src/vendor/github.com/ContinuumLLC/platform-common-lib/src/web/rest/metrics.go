package rest

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rangeType string
type dataType string
type unitType string

// All types for metrics
const (
	RangeTypeZero         rangeType = ""
	RangeTypeToInfinity   rangeType = "(0, infinity)"
	RangeTypeFromInfinity rangeType = "(-infinity, 0)"
	RangeTypeInfinity     rangeType = "(-infinity, infinity)"
	DataTypeZero          dataType  = ""
	DataTypeDecimal       dataType  = "decimal"
	DataTypeInteger       dataType  = "integer"
	UnitTypeZero          unitType  = ""
	UnitTypeNumbers       unitType  = "numbers"
	UnitTypeNumbersSec    unitType  = "numbers/sec"
	UnitTypeMicroseconds  unitType  = "microseconds"
	metricsType                     = "metrics"
)

// All errors
var (
	ErrMetricBadFloatType = errors.New("metric must be a float64 value")
	ErrMetricBadIntType   = errors.New("metric must be a int value")
	ErrEmptyKey           = errors.New("key cannot be empty")
	valueShouldBeInRange  = "value should be in range %s"
)

// Metrics implements Metricer
type Metrics struct {
	GeneralInfo
	Type string `json:"type"`
}

// Metricser is an interface for metrics
type Metricser interface {
	metrics()
	Update() error
}

// Update for custom fields update
func (m *Metrics) Update() error {
	return nil
}

// metrics initialize of Metrics structure
func (m *Metrics) metrics() {
	m.TimeStampUTC = time.Now().UTC()
	m.Type = metricsType
}

// Response struct for response
type Response struct {
	Metrics Metricser       `json:"generalInfo"`
	Data    []MetricsConfig `json:"metrics"`
}

// HandlerMetrics used for returning health status of service
func HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	storage.Metrics.metrics()
	if err := storage.Metrics.Update(); err != nil {
		SendInternalServerError(w, "Update error", fmt.Errorf("rest.HandlerMetrics err: %s", err))
		return
	}
	RenderJSON(w, GetResponse())
}

// GetResponse create new response
func GetResponse() Response {
	storage.Lock()
	defer storage.Unlock()

	res := Response{Metrics: storage.Metrics}
	res.Data = make([]MetricsConfig, len(storage.Data))
	copy(res.Data, storage.Data)
	storage.clear()
	return res
}

// MetricsConfig is exportable type for initialization
type MetricsConfig struct {
	Name     string      `json:"name"`
	DataType dataType    `json:"datatype"`
	Unit     unitType    `json:"unit"`
	Range    rangeType   `json:"range"`
	Value    interface{} `json:"value"`
}

// Set adds value to metrics by key
func Set(key string, value interface{}) error {
	return save(key, value, false)
}

func doInteger(position int, value interface{}, incr bool) error {
	val, ok := value.(int)
	if !ok {
		return ErrMetricBadIntType
	}

	if !rangeValidate(val, storage.Data[position].Range) {
		return errors.New(fmt.Sprintf(valueShouldBeInRange, storage.Data[position].Range))
	}

	if incr {
		val = storage.Data[position].Value.(int) + val
	}

	storage.Data[position].Value = val

	return nil
}

func rangeValidate(v interface{}, rt rangeType) bool {
	switch tv := v.(type) {
	case int:
		if rt == RangeTypeInfinity {
			return true
		}
		if rt == RangeTypeFromInfinity {
			return tv <= 0
		}
		if rt == RangeTypeToInfinity {
			return tv >= 0
		}
	case float64:
		if rt == RangeTypeInfinity {
			return true
		}
		if rt == RangeTypeFromInfinity {
			return tv <= 0
		}
		if rt == RangeTypeToInfinity {
			return tv >= 0
		}
	}
	return false
}

func doDecimal(position int, value interface{}, incr bool) error {
	val, ok := value.(float64)
	if !ok {
		return ErrMetricBadFloatType
	}

	if !rangeValidate(val, storage.Data[position].Range) {
		return errors.New(fmt.Sprintf(valueShouldBeInRange, storage.Data[position].Range))
	}

	if incr {
		val = storage.Data[position].Value.(float64) + val
	}

	storage.Data[position].Value = val

	return nil
}

func save(key string, value interface{}, incr bool) error {
	if key == "" {
		return ErrEmptyKey
	}

	storage.Lock()
	defer storage.Unlock()

	for i := range storage.Data {
		if storage.Data[i].Name == key {
			switch storage.Data[i].DataType {
			case DataTypeInteger:
				return doInteger(i, value, incr)
			case DataTypeDecimal:
				return doDecimal(i, value, incr)
			}
		}
	}
	return nil
}

// Add adds value to metrics by key
func Add(key string, value interface{}) error {
	return save(key, value, true)
}

// RegistryMetrics fo initial metrics configuration
func RegistryMetrics(m Metricser, configs []MetricsConfig) {
	storage.Lock()
	defer storage.Unlock()

	storage.Metrics = m
	storage.Data = toDefault(configs)
	storage.clear()
}

func toDefault(configs []MetricsConfig) []MetricsConfig {
	for i := range configs {
		if configs[i].Range == RangeTypeZero {
			configs[i].Range = RangeTypeInfinity
		}
		if configs[i].Unit == UnitTypeZero {
			configs[i].Unit = UnitTypeNumbers
		}
		if configs[i].DataType == DataTypeZero {
			configs[i].DataType = DataTypeInteger
		}
	}
	return configs
}

var storage metricsStorage

// Safe storage for metrics
type metricsStorage struct {
	Metrics Metricser
	sync.Mutex
	Data []MetricsConfig
}

func (ms *metricsStorage) clear() {
	for i := range storage.Data {
		storage.Data[i].Value = nilValue(storage.Data[i].DataType)
	}
}

func nilValue(t dataType) interface{} {
	switch t {
	case DataTypeInteger:
		return 0
	case DataTypeDecimal:
		return 0.0
	}
	return nil
}
