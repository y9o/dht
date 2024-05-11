//go:build !cgo

package dht

import (
	"fmt"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

var Maximum_number_of_reads uint32 = 1000

func (dht *DHTxx) close() error {
	return nil
}

func (dht *DHTxx) read() (uint32, error) {
	line, err := gpiocdev.RequestLine(dht.chip, dht.offset, gpiocdev.AsInput)
	if err != nil {
		return 0, err
	}
	defer line.Close()

	dht.vals = dht.vals[:0]

	if dht.Config.Pre > 0 {
		time.Sleep(dht.Config.Pre)
	}

	if err := line.Reconfigure(gpiocdev.AsOutput(0)); err != nil {
		return 0, fmt.Errorf("AsOutput(0) error:%w", err)
	}

	if dht.Config.Low > 0 {
		time.Sleep(dht.Config.Low)
	}

	if err := line.SetValue(1); err != nil {
		return 0, fmt.Errorf("SetValue(1) error:%w", err)
	}

	if dht.Config.High > 0 {
		time.Sleep(dht.Config.High)
	}

	if err := line.Reconfigure(gpiocdev.AsInput); err != nil {
		return 0, fmt.Errorf("AsInput error:%w", err)
	}

	var j uint32 = 0
	lastValue := -1
	firstvalue := -1
	maxNegative160us := uint16((float32(dht.Limit) / 5.0) / (1000.0 / 160.0))
	var count uint16
	for ; j < dht.Limit; j++ {
		v, err := line.Value()
		if err != nil {
			return 0, err
		}
		if lastValue != v {
			lastValue = v
			if cap(dht.vals) == len(dht.vals) {
				break
			}
			if firstvalue == -1 {
				firstvalue = v
			} else {
				dht.vals = append(dht.vals, count)
			}
			count = 1
		} else {
			count++
			if count > maxNegative160us {
				break
			}
		}
	}
	dht.fval = firstvalue
	dht.readCount = j
	for i := 0; i <= dht.Config.Vague; i++ {
		data, err := vals2data(dht.vals, firstvalue, 40-i)
		if err == nil {
			return data, err
		}
	}
	return 0, ErrinvalidData
}

func (dht *DHTxx) GetReadTime(count int) (time.Duration, error) {
	line, err := gpiocdev.RequestLine(dht.chip, dht.offset, gpiocdev.AsInput)
	if err != nil {
		return 0, err
	}
	defer line.Close()

	for j := 0; j < 10; j++ {
		if _, err := line.Value(); err != nil {
			return 0, err
		}
	}
	start := time.Now()
	for j := 0; j < count; j++ {
		if _, err := line.Value(); err != nil {
			return 0, err
		}
	}
	ret := time.Since(start)
	return ret, nil
}

func (dht *DHTxx) Dump() {
	value_print(dht.vals, dht.fval)
}

func vals2data(vals []uint16, firstvalue, bit int) (uint32, error) {
	offset := len(vals) - (bit*2 - 1)
	v := 0
	for ; offset >= 0; offset-- {
		if (offset & 1) == 0 {
			v = firstvalue
		} else {
			v = firstvalue ^ 1
		}
		if v == 1 {
			data, err := _vals2data2(vals[offset:])
			if err != nil {
				data, err = _vals2data(vals[offset:])
			}
			return data, err
		}
	}
	return 0, ErrinvalidData
}
func _vals2data(vals []uint16) (uint32, error) {
	var border int
	for _, v := range vals {
		border = max(border, int(v))
	}
	var _data uint64
	for i, v := range vals {
		if i&1 == 0 {
			_data <<= 1
			if int(v)*2 > border {
				_data |= 1
			}
		}
	}
	if !checksum(_data) {
		_data = 0
		for i, v := range vals {
			if i&1 == 0 {
				_data <<= 1
				if int(v)*2 >= border {
					_data |= 1
				}
			}
		}
		if !checksum(_data) {
			return 0, ErrChecksum
		}
	}
	return uint32(_data >> 8), nil
}

func _vals2data2(vals []uint16) (uint32, error) {
	var border int
	for i, v := range vals {
		if i&1 == 1 {
			border += int(v)
		}
	}
	border /= len(vals) / 2
	var _data uint64
	for i, v := range vals {
		if i&1 == 0 {
			_data <<= 1
			if int(v) >= border {
				_data |= 1
			}
		}
	}
	if !checksum(_data) {
		_data = 0
		for i, v := range vals {
			if i&1 == 0 {
				_data <<= 1
				if int(v) > border {
					_data |= 1
				}
			}
		}
		if !checksum(_data) {
			return 0, ErrChecksum
		}
	}
	return uint32(_data >> 8), nil
}
func checksum(bits uint64) bool {
	return uint8(bits>>8)+uint8(bits>>16)+uint8(bits>>24)+uint8(bits>>32) == uint8(bits)
}
