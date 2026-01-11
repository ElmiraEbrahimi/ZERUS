//go:build linux

package memtime

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func osFootprintBytes() (osMemSample, error) {
	file, err := os.Open("/proc/self/smaps_rollup")
	if err != nil {
		return osMemSample{}, err
	}
	defer file.Close()

	var (
		rssKB          uint64
		pssKB          uint64
		privateCleanKB uint64
		privateDirtyKB uint64
		privateHugeKB  uint64

		rssOK       bool
		pssOK       bool
		privateSeen bool
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, valueKB, ok := parseSmapsLineKB(scanner.Text())
		if !ok {
			continue
		}
		switch key {
		case "Rss":
			rssKB = valueKB
			rssOK = true
		case "Pss":
			pssKB = valueKB
			pssOK = true
		case "Private_Clean":
			privateCleanKB = valueKB
			privateSeen = true
		case "Private_Dirty":
			privateDirtyKB = valueKB
			privateSeen = true
		case "Private_Hugetlb":
			privateHugeKB = valueKB
			privateSeen = true
		}
	}
	if err := scanner.Err(); err != nil {
		return osMemSample{}, err
	}
	sample := osMemSample{}
	if rssOK {
		sample.rssBytes = rssKB * 1024
		sample.rssOK = true
	}
	if pssOK {
		sample.pssBytes = pssKB * 1024
		sample.pssOK = true
	}
	if privateSeen {
		sample.ussBytes = (privateCleanKB + privateDirtyKB + privateHugeKB) * 1024
		sample.ussOK = true
	}
	if sample.pssOK {
		sample.footprintBytes = sample.pssBytes
		sample.footprintOK = true
	} else if sample.rssOK {
		sample.footprintBytes = sample.rssBytes
		sample.footprintOK = true
	}

	if !sample.rssOK && !sample.pssOK && !sample.ussOK {
		return sample, fmt.Errorf("smaps_rollup missing expected fields")
	}

	return sample, nil
}

func parseSmapsLineKB(line string) (string, uint64, bool) {
	fields := strings.Fields(line)
	if len(fields) < 3 || fields[2] != "kB" {
		return "", 0, false
	}
	key := strings.TrimSuffix(fields[0], ":")
	value, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0, false
	}
	return key, value, true
}
