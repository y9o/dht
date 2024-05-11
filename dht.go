package dht

import (
	"fmt"
	"time"
	"unsafe"
)

type Cfg struct {
	Pre   time.Duration //
	Low   time.Duration // check datasheet
	High  time.Duration //
	Retry time.Duration //
	Vague int           // The first bit is 0, so the calculation is allowed with 40-$Vague bits of data. Possibility of $Vague bits off measurements 8%(0b1000)=>4%(0b0100)
	// 最初のビットは0なので40-$Vagueビットのデータでの計算を許可する。$Vagueビットずれた計測値になる可能性 8%(0b1000)=>4%(0b0100)
}

// for DHT22 and DHT11
var CfgDHT22 = Cfg{
	Pre:   0,
	Low:   18 * time.Millisecond,
	High:  0,
	Retry: 2 * time.Second,
	Vague: 0,
}

var ErrinvalidData = fmt.Errorf("invalid received data")
var ErrChecksum = fmt.Errorf("received data checksum error")

type DHTxx struct {
	chip      string
	offset    int
	vals      []uint16
	fval      int
	readCount uint32
	Limit     uint32
	Config    Cfg
	pcgo      unsafe.Pointer
}

func New(chip string, pin int) (*DHTxx, error) {
	values := make([]uint16, 0, 90)
	return &DHTxx{
		chip:   chip,
		offset: pin,
		vals:   values,
		Limit:  Maximum_number_of_reads,
		Config: CfgDHT22,
	}, nil
}

func (dht *DHTxx) Close() error {
	return dht.close()
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
			time.Sleep(dht.Config.Retry)
		} else {
			break
		}
	}
	return
}

func value_print(vals []uint16, firstvalue int) {
	hl := []string{"L", "H"}
	hl2 := []string{"-", "="}
	for i, v := range vals {
		fmt.Printf("%2d:%2d:%s:", i, v, hl[firstvalue])
		for range v {
			fmt.Print(hl2[firstvalue])
		}
		fmt.Print("\n")
		firstvalue ^= 1
	}
}
