# dht
golang dht11/dht12/dht22 sensor reader

```go
d, err := dht.New("/dev/gpiochip0", 4)
if err != nil {
	log.Fatal(err)
}
defer d.Close()

var buf dht.DHT22Data
//var buf dht.DHT11Data
//var buf dht.DHT12Data

fmt.Print("read...\r")
if retry, err := d.Read(&buf, 20); err != nil {
	fmt.Print(err)
} else {
	fmt.Printf("🌡️%2.1f℃ (%d) 🌢%2.1f%% (%d) retry:%d\n", buf.Temp(), buf.Temperature, buf.Hum(), buf.Humidity, retry)
}
```


```
🌡️25.1℃ (251) 🌢59.2% (592) retry:2
```


<details>
<summary>example</summary>

build for pi zero on windows
```bat
set GOOS=linux
set GOARCH=arm
set GOARM=6
set CGO_ENABLED=0

go build -ldflags="-s -w"  -trimpath 
move dht22 oridht22
upx --lzma -o dht22 oridht22
del oridht22
```


```bash
$ uname -r
6.6.28+rpt-rpi-v6

$ ./dht22
./dht22 -chip gpiochip0

$ ./dht22 -chip gpiochip0
./dht22 -chip gpiochip0
                -line 0         (ID_SDA)
                -line 1         (ID_SCL)
                -line 2         (GPIO2)
                -line 3         (GPIO3)
                -line 4         (GPIO4)
                -line 5         (GPIO5)
                -line 6         (GPIO6)
                -line 7         (GPIO7)
                -line 8         (GPIO8)
                -line 9         (GPIO9)
                -line 10                (GPIO10)
                -line 11                (GPIO11)
                -line 12                (GPIO12)
                -line 13                (GPIO13)
                -line 14                (GPIO14)
                -line 15                (GPIO15)
                -line 16                (GPIO16)
                -line 17                (GPIO17)
                -line 18                (GPIO18)
                -line 19                (GPIO19)
                -line 20                (GPIO20)
                -line 21                (GPIO21)
                -line 22                (GPIO22)
                -line 23                (GPIO23)
                -line 24                (GPIO24)
                -line 25                (GPIO25)
                -line 26                (GPIO26)
                -line 27                (GPIO27)
                -line 28                (SDA0)
                -line 29                (SCL0)
                -line 30                (CTS0)
                -line 31                (RTS0)
                -line 32                (TXD0)
                -line 33                (RXD0)
                -line 34                (SD1_CLK)
                -line 35                (SD1_CMD)
                -line 36                (SD1_DATA0)
                -line 37                (SD1_DATA1)
                -line 38                (SD1_DATA2)
                -line 39                (SD1_DATA3)
                -line 40                (CAM_GPIO1)
                -line 41                (WL_ON)
                -line 42                (NC)
                -line 43                (WIFI_CLK)
                -line 44                (CAM_GPIO0)
                -line 45                (BT_ON)
                -line 46                (HDMI_HPD_N)
                -line 47                (STATUS_LED_N)
                -line 48                (SD_CLK_R)
                -line 49                (SD_CMD_R)
                -line 50                (SD_DATA0_R)
                -line 51                (SD_DATA1_R)
                -line 52                (SD_DATA2_R)
                -line 53                (SD_DATA3_R)

$ ./dht22 -chip gpiochip0  -line 4
14:30:51 🌡️25.1℃ (251) 🌢59.2% (592) retry:2
14:31:01 🌡️25.2℃ (252) 🌢59.3% (593) retry:0
14:31:12 🌡️25.2℃ (252) 🌢59.2% (592) retry:1
14:31:22 🌡️25.2℃ (252) 🌢59.2% (592) retry:0
14:31:33 🌡️25.2℃ (252) 🌢59.6% (596) retry:1
14:31:45 🌡️25.2℃ (252) 🌢60.6% (606) retry:2
14:31:57 🌡️25.2℃ (252) 🌢59.3% (593) retry:2
14:32:09 🌡️25.2℃ (252) 🌢59.2% (592) retry:2
14:32:20 🌡️25.2℃ (252) 🌢59.2% (592) retry:1
14:32:31 🌡️25.2℃ (252) 🌢59.2% (592) retry:1
```
</details>

## cgo with libgpiod v2 support

This module can be built with `CGO_ENABLED=0` or `CGO_ENABLED=1`.

However, using CGO does not dramatically improve the reading success rate.