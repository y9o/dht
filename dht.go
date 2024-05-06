package dht

import (
	"fmt"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type Sleep struct {
	Pre   time.Duration //
	Low   time.Duration // check datasheet
	High  time.Duration //
	Retry time.Duration //
}

// for DHT22
var DefaultSleep = Sleep{
	Pre:   0,
	Low:   11 * time.Millisecond,
	High:  0,
	Retry: time.Second,
}

var Maximum_number_of_reads uint32 = 1000
var ErrinvalidData = fmt.Errorf("invalid received data")
var ErrChecksum = fmt.Errorf("received data checksum error")

type DHTxx struct {
	line  *gpiocdev.Line
	vals  []uint32
	Limit uint32 //Maximum_number_of_reads
	Sleep Sleep  //DefaultSleep
}

func New(chip string, pin int) (*DHTxx, error) {
	line, err := gpiocdev.RequestLine(chip, pin, gpiocdev.AsInput)
	if err != nil {
		return nil, err
	}
	values := make([]uint32, 0, 41)
	return &DHTxx{
		line:  line,
		vals:  values,
		Limit: Maximum_number_of_reads,
		Sleep: DefaultSleep,
	}, nil
}

func (dht *DHTxx) Close() error {
	return dht.line.Close()
}

func (dht *DHTxx) Read(buf DHTdata, retry int) (retries int, err error) {
	var data uint32
	for retries = 0; retries <= retry; retries++ {
		data, err = dht.read()
		if err == nil {
			if buf.Set(data) {
				err = nil
				break
			}
		}
		if retries < retry {
			time.Sleep(dht.Sleep.Retry)
		}
	}
	return
}
func (dht *DHTxx) read() (uint32, error) {
	line := dht.line
	values := dht.vals[:0]

	if dht.Sleep.Pre > 0 {
		time.Sleep(dht.Sleep.Pre)
	}

	if err := line.Reconfigure(gpiocdev.AsOutput(0)); err != nil {
		return 0, fmt.Errorf("AsOutput(0) error:%w", err)
	}

	if dht.Sleep.Low > 0 {
		time.Sleep(dht.Sleep.Low)
	}

	if err := line.SetValue(1); err != nil {
		return 0, fmt.Errorf("SetValue(1) error:%w", err)
	}

	if dht.Sleep.High > 0 {
		time.Sleep(dht.Sleep.High)
	}

	if err := line.Reconfigure(gpiocdev.AsInput); err != nil {
		return 0, fmt.Errorf("AsInput error:%w", err)
	}

	var j uint32 = 0
	last := 1
	for ; j < dht.Limit; j++ {
		v, err := line.Value()
		if err != nil {
			return 0, err
		}
		if last != v && v == 1 {
			i := j
			for j < dht.Limit {
				v, err := line.Value()
				if err != nil {
					return 0, err
				}
				if v == 0 {
					break
				}
				j++
			}
			values = append(values, j-i)
			if len(values) == 41 {
				break
			}
		} else {
			last = v
		}
	}
	if len(values) != 41 {
		if len(values) > 0 {
			return 0, fmt.Errorf("len(values)=%d,values[0]=%d %w", len(values), values[0], ErrinvalidData)
		}
		return 0, fmt.Errorf("j=%d last=%d %w", j, last, ErrinvalidData)
	}
	border := values[0] //80ns
	var recvBits uint64
	for i, v := range values[1:] {
		if v*2 > border {
			recvBits |= (1 << (39 - i))
		}
	}
	sum := uint8(recvBits>>8) + uint8(recvBits>>16) + uint8(recvBits>>24) + uint8(recvBits>>32)
	if sum != uint8(recvBits) {
		// fmt.Print("checksum\n")
		// for i, v := range values {
		// 	fmt.Printf("%03d: val:%03d ", i, v)
		// 	for range v {
		// 		fmt.Print("#")
		// 	}
		// 	fmt.Print("#\n")
		// }
		return 0, ErrChecksum
	}
	return uint32(recvBits >> 8), nil
}
