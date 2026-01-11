//go:build darwin

package memtime

/*
#include <errno.h>
#include <libproc.h>
#include <sys/resource.h>
#include <unistd.h>

static int memtime_proc_pid_rusage(struct rusage_info_v2 *info, int *err) {
	int rc = proc_pid_rusage(getpid(), RUSAGE_INFO_V2, (rusage_info_t *)info);
	if (rc != 0 && err != NULL) {
		*err = errno;
	}
	return rc;
}
*/
import "C"

import (
	"fmt"
	"syscall"
)

func osFootprintBytes() (osMemSample, error) {
	var info C.struct_rusage_info_v2
	var errnum C.int
	ret := C.memtime_proc_pid_rusage(&info, &errnum)
	if ret != 0 {
		return osMemSample{}, fmt.Errorf("proc_pid_rusage: %v", syscall.Errno(errnum))
	}

	return osMemSample{
		rssBytes:       uint64(info.ri_resident_size),
		rssOK:          true,
		footprintBytes: uint64(info.ri_phys_footprint),
		footprintOK:    true,
	}, nil
}
