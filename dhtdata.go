package dht

type DHTdata interface {
	Set(raw uint32) bool
}

type DHTxxData struct {
	Humidity    uint16
	Temperature uint16
}

func (data *DHTxxData) Set(raw uint32) bool {
	data.Humidity = uint16(raw >> 16)
	data.Temperature = uint16(raw)
	return true
}

type DHT22Data struct {
	Humidity    uint16
	Temperature int16
}

func (data *DHT22Data) Set(raw uint32) bool {
	data.Humidity = uint16(raw >> 16)
	data.Temperature = int16(raw & 0x7FFF)
	if (raw & 0x8000) != 0 {
		data.Temperature *= -1
	}
	if data.Temperature < -400 || data.Temperature > 800 {
		return false
	}
	if data.Humidity > 1000 {
		return false
	}
	return true
}
func (data DHT22Data) Hum() float32 {
	return float32(data.Humidity) / 10
}
func (data DHT22Data) Temp() float32 {
	return float32(data.Temperature) / 10
}

type DHT11Data DHT22Data

func (data *DHT11Data) Set(raw uint32) bool {
	data.Humidity = uint16(raw >> 24)
	data.Temperature = int16(raw>>8) & 0xFF
	if data.Humidity > 90 || data.Humidity < 20 {
		return false
	}
	if data.Temperature > 50 {
		return false
	}
	return true
}
func (data DHT11Data) Hum() float32 {
	return float32(data.Humidity)
}
func (data DHT11Data) Temp() float32 {
	return float32(data.Temperature)
}

type DHT12Data DHT22Data

func (data *DHT12Data) Set(raw uint32) bool {
	data.Humidity = uint16(raw>>24)*10 + (uint16(raw>>16) & 0xFF)
	data.Temperature = (int16(raw>>8)&0xFF)*10 + int16(raw&0x7F)
	if (raw & 0x80) != 0 {
		data.Temperature *= -1
	}
	if data.Humidity > 950 || data.Humidity < 200 {
		return false
	}
	if data.Temperature < -200 || data.Temperature > 600 {
		return false
	}
	return true
}
func (data DHT12Data) Hum() float32 {
	return float32(data.Humidity) / 10
}
func (data DHT12Data) Temp() float32 {
	return float32(data.Temperature) / 10
}
