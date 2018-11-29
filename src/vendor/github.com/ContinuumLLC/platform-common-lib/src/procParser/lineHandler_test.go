package procParser

import "testing"
import "reflect"

func TestModeHandlerFactoryTabular(t *testing.T) {
	factory := new(modeHandlerFactoryImpl)
	handler := factory.GetModeHandler(ModeTabular)
	if handler == nil {
		t.Error("Unable to create handler")
	}
}

func TestModeHandlerFactoryKeyValue(t *testing.T) {
	factory := new(modeHandlerFactoryImpl)
	handler := factory.GetModeHandler(ModeKeyValue)
	if handler == nil {
		t.Error("Unable to create handler")
	}
}

func TestModeHandlerFactoryInvalidHandler(t *testing.T) {
	factory := new(modeHandlerFactoryImpl)
	handler := factory.GetModeHandler(3)
	if handler != nil {
		t.Error("Created handler for incorrect handler mode")
	}
}

func TestModeHandlerFactoryWithData(t *testing.T) {
	var testData = []struct {
		testName  string
		strString string
		keyIndex  int
		expKey    string
		fileMode  Mode
		cfg       Config
		feildLen  int
		expLine   []string
	}{
		{
			"Test1",
			"key value1 value2 value3",
			0,
			"key",
			ModeTabular,
			Config{ParserMode: ModeTabular},
			4,
			[]string{"key", "value1", "value2", "value3"},
		},
		{
			"Test2",
			"key     value1        value2              value3",
			0,
			"key",
			ModeTabular,
			Config{ParserMode: ModeTabular},
			4,
			[]string{"key", "value1", "value2", "value3"},
		},
		{
			"Test3",
			"key value1 value2 value3",
			5,
			"",
			ModeTabular,
			Config{ParserMode: ModeTabular},
			4,
			[]string{"key", "value1", "value2", "value3"},
		},
		{
			"Test4",
			"key:value1 value2 value3",
			0,
			"key",
			ModeKeyValue,
			Config{ParserMode: ModeKeyValue},
			4,
			[]string{"key", "value1", "value2", "value3"},
		},
		{
			"Test5",
			" key : value1 value2 value3 ",
			0,
			"key",
			ModeSeparator,
			Config{ParserMode: ModeSeparator, Separator: ":"},
			2,
			[]string{"key", "value1 value2 value3"},
		},
		{
			"Test6",
			" key : value1 value2 value3 ",
			0,
			"key",
			ModeSeparator,
			Config{ParserMode: ModeSeparator, Separator: ":"},
			2,
			[]string{"key", "value1 value2 value3"},
		},
		{
			"Test7",
			" key :		value1 value2		value3 		",
			0,
			"key",
			ModeSeparator,
			Config{ParserMode: ModeSeparator, Separator: ":"},
			2,
			[]string{"key", "value1 value2		value3"},
		},
		{
			"Test8",
			"key   :  value1        value2              value3",
			0,
			"key",
			ModeKeyValue,
			Config{ParserMode: ModeKeyValue},
			4,
			[]string{"key", "value1", "value2", "value3"},
		},
	}

	for _, tdata := range testData {
		factory := new(modeHandlerFactoryImpl)
		handler := factory.GetModeHandler(tdata.fileMode)
		line := handler.HandleLine(tdata.strString, tdata.cfg)
		key, _ := getKeyValue(line.Values, tdata.keyIndex)
		if key != tdata.expKey {
			t.Errorf("Test = %s, Expected key value is |%s|, but received |%s|", tdata.testName, tdata.expKey, key)
		}
		if len(line.Values) != tdata.feildLen {
			t.Errorf("Test = %s, Expected field len is %d, but received %d", tdata.testName, len(line.Values), tdata.feildLen)
		}
		if !reflect.DeepEqual(line.Values, tdata.expLine) {
			t.Errorf("Test = %s, Expected line is %v but received %v", tdata.testName, tdata.expLine, line.Values)
		}
	}
}

func TestSplitLines(t *testing.T) {
	var testData = []struct {
		strTest string
		iExpLen int
	}{
		{
			"",
			0,
		},
		{
			"abc def ghi",
			3,
		},
	}

	for _, tdata := range testData {
		values := splitLines(tdata.strTest, " ")
		if tdata.iExpLen != len(values) {
			t.Errorf("Expected length is %d, but received length is %d", tdata.iExpLen, len(values))
		}
	}
}
