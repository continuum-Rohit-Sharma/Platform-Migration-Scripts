package pdh

import (
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modpdhdll                       = syscall.NewLazyDLL("pdh.dll")
	procPdhOpenQuery                = modpdhdll.NewProc("PdhOpenQueryW")
	procPdhCloseQuery               = modpdhdll.NewProc("PdhCloseQuery")
	procPdhAddCounter               = modpdhdll.NewProc("PdhAddCounterW")
	procPdhCollectQueryData         = modpdhdll.NewProc("PdhCollectQueryData")
	procPdhGetFormattedCounterValue = modpdhdll.NewProc("PdhGetFormattedCounterValue")
	procPdhRemoveCounter            = modpdhdll.NewProc("PdhRemoveCounter")
	procPdhGetFormattedCounterArray = modpdhdll.NewProc("PdhGetFormattedCounterArrayW")
	procPdhParseCounterPath         = modpdhdll.NewProc("PdhParseCounterPathW")
)

//OpenQuery Creates a new query that is used to manage the collection of performance data.
func OpenQuery(datasrc string, userData uint32, phQuery *windows.Handle) (ret uint32) {
	var pDataSrc *uint16
	if len(datasrc) > 0 {
		pDataSrc, _ = windows.UTF16PtrFromString(datasrc)
	}

	p0, _, _ := procPdhOpenQuery.Call(uintptr(unsafe.Pointer(pDataSrc)),
		uintptr(userData), uintptr(unsafe.Pointer(phQuery)))
	return uint32(p0)
}

//CloseQuery Closes all counters contained in the specified query, closes all handles related to the query, and frees all memory associated with the query.
func CloseQuery(phQuery windows.Handle) (ret uint32) {
	p0, _, _ := procPdhCloseQuery.Call(uintptr(phQuery))
	return uint32(p0)
}

//AddCounter Adds the specified counter to the query.
func AddCounter(phQuery windows.Handle, szFullCounterPath string, userData uint32, phCounter *windows.Handle) (ret uint32) {
	pCtrPath, _ := windows.UTF16PtrFromString(szFullCounterPath)
	p0, _, _ := procPdhAddCounter.Call(uintptr(phQuery),
		uintptr(unsafe.Pointer(pCtrPath)),
		uintptr(userData),
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(p0)
}

//CollectQueryData Collects the current raw data value for all counters in the specified query and updates the status code of each counter.
func CollectQueryData(phQuery windows.Handle) (ret uint32) {
	p0, _, _ := procPdhCollectQueryData.Call(uintptr(phQuery))
	return uint32(p0)
}

//GetFormattedCounterValueDouble get the value
func GetFormattedCounterValueDouble(phCounter windows.Handle, counterType *uint32, pValue *FmtCounterValueDouble) (ret uint32) {
	p0, _, _ := procPdhGetFormattedCounterValue.Call(uintptr(phCounter),
		uintptr(FmtDouble),
		uintptr(unsafe.Pointer(counterType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(p0)
}

//GetFormattedCounterArrayDouble gets all the instance values in double
func GetFormattedCounterArrayDouble(phCounter windows.Handle, bufferSize *uint32, bufferCount *uint32, pValue *FmtCounterValueItemDouble) (ret uint32) {
	p0, _, _ := procPdhGetFormattedCounterArray.Call(uintptr(phCounter),
		uintptr(FmtDouble),
		uintptr(unsafe.Pointer(bufferSize)),
		uintptr(unsafe.Pointer(bufferCount)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(p0)
}

//GetFormattedCounterValueLarge get the value in Large
func GetFormattedCounterValueLarge(phCounter windows.Handle, counterType *uint32, pValue *FmtCounterValueLarge) (ret uint32) {
	p0, _, _ := procPdhGetFormattedCounterValue.Call(uintptr(phCounter),
		uintptr(FmtLarge),
		uintptr(unsafe.Pointer(counterType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(p0)
}

//GetFormattedCounterArrayLarge gets all the instance values in large
func GetFormattedCounterArrayLarge(phCounter windows.Handle, bufferSize *uint32, bufferCount *uint32, pValue *FmtCounterValueItemLarge) (ret uint32) {
	p0, _, _ := procPdhGetFormattedCounterArray.Call(uintptr(phCounter),
		uintptr(FmtLarge),
		uintptr(unsafe.Pointer(bufferSize)),
		uintptr(unsafe.Pointer(bufferCount)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(p0)
}

//RemoveCounter removes the counter
func RemoveCounter(phCounter windows.Handle) (ret uint32) {
	p0, _, _ := procPdhRemoveCounter.Call(uintptr(phCounter))
	return uint32(p0)
}

//ParseCounterPath parses the given perf path
func ParseCounterPath(FullPathBuffer string, pCounterPathElements *CounterPathElements, bufferSize *uint32, dwFlags uint32) (ret uint32) {
	pCtrPath, _ := windows.UTF16PtrFromString(FullPathBuffer)
	p0, _, _ := procPdhParseCounterPath.Call(
		uintptr(unsafe.Pointer(pCtrPath)),
		uintptr(unsafe.Pointer(pCounterPathElements)),
		uintptr(unsafe.Pointer(bufferSize)),
		uintptr(dwFlags))
	return uint32(p0)
}

//UTF16PtrToString converts utf16 pointer string to string
func UTF16PtrToString(p *uint16) string {
	return LpOleStrToString(p)
}

// LpOleStrToString converts COM Unicode to Go string.
func LpOleStrToString(p *uint16) string {
	if p == nil {
		return ""
	}

	length := LpOleStrLen(p)
	a := make([]uint16, length)

	ptr := unsafe.Pointer(p)

	for i := 0; i < int(length); i++ {
		a[i] = *(*uint16)(ptr)
		ptr = unsafe.Pointer(uintptr(ptr) + 2)
	}

	return string(utf16.Decode(a))
}

// LpOleStrLen returns the length of Unicode string.
func LpOleStrLen(p *uint16) (length int64) {
	if p == nil {
		return 0
	}

	ptr := unsafe.Pointer(p)

	for i := 0; ; i++ {
		if 0 == *(*uint16)(ptr) {
			length = int64(i)
			break
		}
		ptr = unsafe.Pointer(uintptr(ptr) + 2)
	}
	return
}

func removePerfCounters(perfCounter []perfHCounter) {
	for _, perfCtr := range perfCounter {
		if 0 != perfCtr.hCounter {
			_ = RemoveCounter(perfCtr.hCounter)
		}
	}
}

func addPerfCounters(hQuery windows.Handle, perfPath []PerfPaths) (perfCounter []perfHCounter) {
	ilen := len(perfPath)
	perfCounter = make([]perfHCounter, ilen)

	for iIndex, path := range perfPath {
		perfCounter[iIndex].paths = path

		var size = uint32(unsafe.Sizeof(CounterPathElements{}))
		var buffSize = (20 * size)

		parsedCounterBuf := make([]CounterPathElements, buffSize)
		perfCounter[iIndex].retError = ParseCounterPath(path.FullPerfPath, &parsedCounterBuf[0], &buffSize, 0)
		if MoreData == perfCounter[iIndex].retError {
			parsedCounterBuf := make([]CounterPathElements, buffSize)
			perfCounter[iIndex].retError = ParseCounterPath(path.FullPerfPath, &parsedCounterBuf[0], &buffSize, 0)
		}

		if ErrorSuccess != perfCounter[iIndex].retError {
			continue
		}

		perfCounter[iIndex].counterInfo = parsedCounterBuf[0]
		perfCounter[iIndex].retError = AddCounter(hQuery, path.FullPerfPath, 0, &perfCounter[iIndex].hCounter)
	}

	return
}

func getPerfLargeValue(perfCounter perfHCounter) (perfData []PerfData) {
	var bufSize uint32
	var bufCount uint32
	var size = uint32(unsafe.Sizeof(FmtCounterValueItemLarge{}))

	var filledBuf []FmtCounterValueItemLarge

	bufCount = 500
	bufSize = (bufCount * size)

	filledBuf = make([]FmtCounterValueItemLarge, bufSize)
	ret := GetFormattedCounterArrayLarge(perfCounter.hCounter, &bufSize, &bufCount, &filledBuf[0])

	if MoreData == ret {
		bufCount *= 2
		bufSize = (bufCount * size)

		filledBuf = make([]FmtCounterValueItemLarge, bufSize)
		ret = GetFormattedCounterArrayLarge(perfCounter.hCounter, &bufSize, &bufCount, &filledBuf[0])
	}
	if ErrorSuccess != ret {
		perfData = make([]PerfData, 1)
		perfData[0].RetError = ret
		perfData[0].Format = FmtLarge
		perfData[0].MachineName = UTF16PtrToString(perfCounter.counterInfo.MachineName)
		perfData[0].ObjectName = UTF16PtrToString(perfCounter.counterInfo.ObjectName)
		perfData[0].CounterName = UTF16PtrToString(perfCounter.counterInfo.CounterName)
		perfData[0].InstanceName = UTF16PtrToString(perfCounter.counterInfo.InstanceName)
		perfData[0].ParentInstance = UTF16PtrToString(perfCounter.counterInfo.ParentInstance)
		perfData[0].InstanceIndex = perfCounter.counterInfo.InstanceIndex
		return
	}

	perfData = make([]PerfData, bufCount)

	for i := 0; i < int(bufCount); i++ {
		perfData[i].Format = FmtLarge
		perfData[i].MachineName = UTF16PtrToString(perfCounter.counterInfo.MachineName)
		perfData[i].ObjectName = UTF16PtrToString(perfCounter.counterInfo.ObjectName)
		perfData[i].CounterName = UTF16PtrToString(perfCounter.counterInfo.CounterName)
		perfData[i].InstanceName = UTF16PtrToString(filledBuf[i].Name)
		perfData[i].ParentInstance = UTF16PtrToString(perfCounter.counterInfo.ParentInstance)
		perfData[i].InstanceIndex = perfCounter.counterInfo.InstanceIndex
		perfData[i].LargeValue = filledBuf[i].FmtValue.Value
	}
	return
}

func getPerfDoubleValue(perfCounter perfHCounter) (perfData []PerfData) {
	var bufSize uint32
	var bufCount uint32
	var size = uint32(unsafe.Sizeof(FmtCounterValueItemDouble{}))

	var filledBuf []FmtCounterValueItemDouble

	bufCount = 500
	bufSize = (bufCount * size)

	filledBuf = make([]FmtCounterValueItemDouble, bufSize)
	ret := GetFormattedCounterArrayDouble(perfCounter.hCounter, &bufSize, &bufCount, &filledBuf[0])

	if MoreData == ret {
		bufCount *= 2
		bufSize = (bufCount * size)

		filledBuf = make([]FmtCounterValueItemDouble, bufSize)
		ret = GetFormattedCounterArrayDouble(perfCounter.hCounter, &bufSize, &bufCount, &filledBuf[0])
	}
	if ErrorSuccess != ret {
		perfData = make([]PerfData, 1)
		perfData[0].RetError = ret
		perfData[0].Format = FmtDouble
		perfData[0].MachineName = UTF16PtrToString(perfCounter.counterInfo.MachineName)
		perfData[0].ObjectName = UTF16PtrToString(perfCounter.counterInfo.ObjectName)
		perfData[0].CounterName = UTF16PtrToString(perfCounter.counterInfo.CounterName)
		perfData[0].InstanceName = UTF16PtrToString(perfCounter.counterInfo.InstanceName)
		perfData[0].ParentInstance = UTF16PtrToString(perfCounter.counterInfo.ParentInstance)
		perfData[0].InstanceIndex = perfCounter.counterInfo.InstanceIndex
		return
	}

	perfData = make([]PerfData, bufCount)

	for i := 0; i < int(bufCount); i++ {
		perfData[i].Format = FmtDouble
		perfData[i].MachineName = UTF16PtrToString(perfCounter.counterInfo.MachineName)
		perfData[i].ObjectName = UTF16PtrToString(perfCounter.counterInfo.ObjectName)
		perfData[i].CounterName = UTF16PtrToString(perfCounter.counterInfo.CounterName)
		perfData[i].InstanceName = UTF16PtrToString(filledBuf[i].Name)
		perfData[i].ParentInstance = UTF16PtrToString(perfCounter.counterInfo.ParentInstance)
		perfData[i].InstanceIndex = perfCounter.counterInfo.InstanceIndex
		perfData[i].DoubleValue = filledBuf[i].FmtValue.Value
	}
	return
}

//QueryPerformanceData gets the given performance data
func QueryPerformanceData(perfPath []PerfPaths) (perfData []PerfData, ret uint32) {

	ilen := len(perfPath)
	if ilen <= 0 {
		return
	}

	//Open the PDH Query
	var hQuery windows.Handle
	var strSrc string
	ret = OpenQuery(strSrc, 0, &hQuery)
	if ErrorSuccess != ret {
		return
	}

	perfCounter := addPerfCounters(hQuery, perfPath)

	//Now its time to collect the data
	ret = CollectQueryData(hQuery)
	if ErrorSuccess != ret {
		return
	}
	time.Sleep(1 * time.Second) //This is needed since some counters needs two data sets after 1 second
	ret = CollectQueryData(hQuery)
	if ErrorSuccess != ret {
		return
	}

	for _, perfCtr := range perfCounter {
		if FmtLarge == perfCtr.paths.Format && ErrorSuccess == perfCtr.retError {
			perfData = append(perfData, getPerfLargeValue(perfCtr)...)
		} else if FmtDouble == perfCtr.paths.Format && ErrorSuccess == perfCtr.retError {
			perfData = append(perfData, getPerfDoubleValue(perfCtr)...)
		}
	}

	removePerfCounters(perfCounter)
	CloseQuery(hQuery)
	return
}
