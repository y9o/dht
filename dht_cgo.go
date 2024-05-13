//go:build cgo

package dht

/*
#cgo CFLAGS: -I${SRCDIR}/libgpiod/include
#cgo LDFLAGS: -lgpiod
#include <stdlib.h>
#include "dht.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

var Maximum_number_of_reads uint32 = 1500

func (dht *DHTxx) close() error {
	if dht.pcgo != nil {
		C.freeDHT((*C.dht_info)(dht.pcgo))
		dht.pcgo = nil
	}
	return nil
}

func (dht *DHTxx) read() (uint32, error) {
	cfg := C.dht_config{}
	cfg.pre_sleep = C.int(dht.Config.Pre / time.Microsecond)
	cfg.high_sleep = C.int(dht.Config.High / time.Microsecond)
	cfg.low_sleep = C.int(dht.Config.Low / time.Microsecond)
	cfg.limit = C.int(dht.Limit)
	cfg.vague = C.int(dht.Config.Vague)
	if dht.pcgo == nil {
		chip := C.CString(dht.chip)
		p := C.newDHT(chip, C.uint(dht.offset))
		C.free(unsafe.Pointer(chip))
		if p == nil {
			return 0, errors.New("failed newDHT")
		}
		dht.pcgo = unsafe.Pointer(p)
	}
	data := C.uint(0)
	ret := C.readDHT((*C.dht_info)(dht.pcgo), &cfg, &data)
	if ret > 0 {
		return uint32(data), nil
	} else {
		return 0, fmt.Errorf("faild readDHT(%d)", ret)
	}
}

func (dht *DHTxx) GetReadTime(count int) (time.Duration, error) {
	if dht.pcgo == nil {
		chip := C.CString(dht.chip)
		p := C.newDHT(chip, C.uint(dht.offset))
		C.free(unsafe.Pointer(chip))
		if p == nil {
			return 0, errors.New("failed newDHT")
		}
		dht.pcgo = unsafe.Pointer(p)
	}
	ret := C.getReadTimeDHT((*C.dht_info)(dht.pcgo), C.int(count))
	if ret < 0 {
		return 0, fmt.Errorf("faild getReadTimeDHT")
	}
	return time.Duration(ret), nil
}

func (dht *DHTxx) Dump() {
	if dht.pcgo == nil {
		return
	}
	dht.vals = dht.vals[:cap(dht.vals)]
	ret := C.copyBufDHT((*C.dht_info)(dht.pcgo), (*C.ushort)(unsafe.Pointer(&dht.vals[0])), C.int(len(dht.vals)))
	if ret != -1 {
		for i := range cap(dht.vals) {
			if dht.vals[i] == 0 {
				dht.vals = dht.vals[:i]
				break
			}
		}
	}
	value_print(dht.vals, int(ret))
}
