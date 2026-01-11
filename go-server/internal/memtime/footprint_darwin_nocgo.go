//go:build darwin && !cgo

package memtime

import "fmt"

func osFootprintBytes() (osMemSample, error) {
	return osMemSample{}, fmt.Errorf("memtime: phys_footprint requires cgo")
}
