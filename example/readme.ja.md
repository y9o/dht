
# for Raspberry Pi Zero

pi zero 32bit を対象にしたクロスビルドです。

その他の場合はパラメータを置き換えてください。

```bash
$ cat /etc/os-release
PRETTY_NAME="Raspbian GNU/Linux 12 (bookworm)"
NAME="Raspbian GNU/Linux"
VERSION_ID="12"
VERSION="12 (bookworm)"
VERSION_CODENAME=bookworm
ID=raspbian
ID_LIKE=debian
HOME_URL="http://www.raspbian.org/"
SUPPORT_URL="http://www.raspbian.org/RaspbianForums"
BUG_REPORT_URL="http://www.raspbian.org/RaspbianBugs"
$ uname -a
Linux pi0 6.6.28+rpt-rpi-v6 #1 Raspbian 1:6.6.28-1+rpt1 (2024-04-22) armv6l GNU/Linux
```

## CGO use libgpiod 2.x.x

cgoではlibgpiod Version 2が必要です。(linux 5.10)

クロスビルド用にdockerの中で実行するMakefileがあります。

Makefileでは現時点の最新版libgpiod-2.1.1.tar.xzをダウンロードしてstatic buildします。

```bash
$ docker run -it --rm -v $PWD:/work -w /work \
	docker.elastic.co/beats-dev/golang-crossbuild:1.22.2-armel-debian12 \
	-p "linux/armv6" \
	--build-cmd "make all"
```

## non CGO

```bat
set GOOS=linux
set GOARCH=arm
set GOARM=6
set CGO_ENABLED=0

go build
```
```bash 
$ GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build
```

## どっち

GPIOの読み取りはCGOのほうが2倍早いので正確なデータが取れるはずだけど、今のコードで私の環境(DHT22/AM2302)ではリトライ回数に大きな違いはなさそう。

### CGO_ENABLED=0
```
GetReadTime(10000loop)=> 63.204044ms / 10000 = 6.32µs
6000ms / 6.32µs = 949
^C==========
total  :  27
success:  21   77.778
retry  :   6   22.222
```
### CGO_ENABLED=1
```
GetReadTime(10000loop)=> 27.846608ms / 10000 = 2.784µs
6000ms / 2.784µs = 2155
^C==========
total  :  39
success:  27   69.231
retry  :  12   30.769
```