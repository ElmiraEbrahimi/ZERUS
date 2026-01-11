//go:build windows

package memtime

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	psapiDLL                 = windows.NewLazySystemDLL("psapi.dll")
	getProcessMemoryInfoProc = psapiDLL.NewProc("GetProcessMemoryInfo")
)

type processMemoryCountersEx struct {
	cb                         uint32
	pageFaultCount             uint32
	peakWorkingSetSize         uintptr
	workingSetSize             uintptr
	quotaPeakPagedPoolUsage    uintptr
	quotaPagedPoolUsage        uintptr
	quotaPeakNonPagedPoolUsage uintptr
	quotaNonPagedPoolUsage     uintptr
	pagefileUsage              uintptr
	peakPagefileUsage          uintptr
	privateUsage               uintptr
}

func osFootprintBytes() (osMemSample, error) {
	var mem processMemoryCountersEx
	mem.cb = uint32(unsafe.Sizeof(mem))
	r1, _, e1 := getProcessMemoryInfoProc.Call(
		uintptr(windows.CurrentProcess()),
		uintptr(unsafe.Pointer(&mem)),
		uintptr(mem.cb),
	)
	if r1 == 0 {
		if e1 != 0 {
			return osMemSample{}, error(e1)
		}
		return osMemSample{}, syscall.EINVAL
	}

	return osMemSample{
		rssBytes:       uint64(mem.workingSetSize),
		rssOK:          true,
		footprintBytes: uint64(mem.privateUsage),
		footprintOK:    true,
	}, nil
}
