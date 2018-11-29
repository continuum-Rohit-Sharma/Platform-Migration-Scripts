package pdh

import "golang.org/x/sys/windows"

//FmtCounterValueDouble returns the value in double format
type FmtCounterValueDouble struct {
	CStatus uint32
	Pad     [4]byte
	Value   float64
}

//FmtCounterValueLarge returns the value in large format
type FmtCounterValueLarge struct {
	CStatus uint32
	Pad     [4]byte
	Value   int64
}

//FmtCounterValueItemDouble used for getting multiple instance data in double format
type FmtCounterValueItemDouble struct {
	Name     *uint16 // pointer to a string
	Pad      [4]byte
	FmtValue FmtCounterValueDouble
}

//FmtCounterValueItemLarge used for getting multiple instance data in large format
type FmtCounterValueItemLarge struct {
	Name     *uint16 // pointer to a string
	Pad      [4]byte
	FmtValue FmtCounterValueLarge
}

//CounterPathElements structure contains the components of a counter path.
type CounterPathElements struct {
	MachineName    *uint16
	ObjectName     *uint16
	InstanceName   *uint16
	ParentInstance *uint16
	InstanceIndex  uint32
	CounterName    *uint16
}

//PerfPaths expects perf paths
type PerfPaths struct {
	FullPerfPath string
	Format       int
}

//This is a internal struct used by the API
type perfHCounter struct {
	paths       PerfPaths
	counterInfo CounterPathElements
	hCounter    windows.Handle
	retError    uint32
}

//PerfData provides perf data
type PerfData struct {
	MachineName    string
	ObjectName     string
	InstanceName   string
	ParentInstance string
	InstanceIndex  uint32
	CounterName    string
	DoubleValue    float64
	LargeValue     int64
	Format         int
	RetError       uint32
}
