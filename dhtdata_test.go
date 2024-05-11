package dht

import "testing"

func TestDHT22Data_Set(t *testing.T) {
	type args struct {
		raw uint32
		tmp float32
		hum float32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"40 100", args{0b00000011_11101000_00000001_10010000, 40, 100}, true},
		{"40.0,100.1", args{0b00000011_11101001_00000001_10010000, 40, 100.1}, false},
		{"-40.0,100.0", args{0b00000011_11101000_10000001_10010000, -40, 100}, true},
		{"-40.1,100.0", args{0b00000011_11101000_10000001_10010001, -40.1, 100}, false},
		{"80.0,0.1", args{0b00000000_00000001_00000011_00100000, 80, 0.1}, true},
		{"80.1,100.0", args{0b00000000_00000001_00000011_00100001, 80.1, 0.1}, false},
		{"80.0,0.1", args{0b00000000_00000001_00000011_00100000, 80, 0.1}, true},
		{"80.0,3276.9", args{0b10000000_00000001_00000011_00100000, 80, 3276.9}, false},
		{"-0.1,99.9", args{0b00000011_11100111_10000000_00000001, -0.1, 99.9}, true},
	}
	var buf DHT22Data
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buf.Set(tt.args.raw); got != tt.want {
				t.Errorf("%s:DHT22Data.Set() = %v, want %v", tt.name, got, tt.want)
			}
			if got := buf.Temp(); got != tt.args.tmp {
				t.Errorf("%s:DHT22Data.Temp() = %v, want %v ", tt.name, got, tt.args.tmp)
			}
			if got := buf.Hum(); got != tt.args.hum {
				t.Errorf("%s:DHT22Data.Hum() = %v, want %v ", tt.name, got, tt.args.hum)
			}
		})
	}
}

func TestDHT11Data_Set(t *testing.T) {
	type args struct {
		raw uint32
		hum float32
		tmp float32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"33 26", args{0b00100001_00000000_00011010_00000000, 33, 26}, true},
	}
	var buf DHT11Data
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buf.Set(tt.args.raw); got != tt.want {
				t.Errorf("%s:DHT11Data.Set() = %v, want %v", tt.name, got, tt.want)
			}
			if got := buf.Temp(); got != tt.args.tmp {
				t.Errorf("%s:DHT11Data.Temp() = %v, want %v ", tt.name, got, tt.args.tmp)
			}
			if got := buf.Hum(); got != tt.args.hum {
				t.Errorf("%s:DHT11Data.Hum() = %v, want %v ", tt.name, got, tt.args.hum)
			}
		})
	}
}

func TestDHT12Data_Set(t *testing.T) {
	type args struct {
		raw uint32
		hum float32
		tmp float32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"56.8 26.6", args{0b00111000_00001000_00011010_00000110, 56.8, 26.6}, true},
		{"56.8 -26.6", args{0b00111000_00001000_00011010_10000110, 56.8, -26.6}, false},
	}
	var buf DHT12Data
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buf.Set(tt.args.raw); got != tt.want {
				t.Errorf("%s:DHT12Data.Set() = %v, want %v", tt.name, got, tt.want)
			}
			if got := buf.Temp(); got != tt.args.tmp {
				t.Errorf("%s:DHT12Data.Temp() = %v, want %v ", tt.name, got, tt.args.tmp)
			}
			if got := buf.Hum(); got != tt.args.hum {
				t.Errorf("%s:DHT12Data.Hum() = %v, want %v ", tt.name, got, tt.args.hum)
			}
		})
	}
}
