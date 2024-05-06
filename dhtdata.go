package dht

type DHTdata interface {
	Set(raw uint32) bool
}

type DHTxxData struct {
	Temperature uint16
	Humidity    uint16
}

func (data *DHTxxData) Set(raw uint32) bool {
	data.Temperature = uint16(raw)
	data.Humidity = uint16(raw >> 16)
	return true
}

type DHT22Data DHTxxData

func (data *DHT22Data) Set(raw uint32) bool {
	data.Temperature = uint16(raw)
	data.Humidity = uint16(raw >> 16)
	temp := data.Temp()
	if temp < -40 || temp > 80 {
		return false
	}
	if data.Humidity > 1000 {
		return false
	}
	return true
}
func (data DHT22Data) Temp() float32 {
	if data.Temperature&0x8000 != 0 {
		return float32(-int16(data.Temperature&0x7fff)) / 10
	} else {
		return float32(data.Temperature) / 10
	}
}
func (data DHT22Data) Hum() float32 {
	return float32(data.Humidity) / 10
}
