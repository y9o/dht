CD = $(shell pwd)
LIBGPIOD=$(CD)/libgpiod
export C_INCLUDE_PATH=$(LIBGPIOD)/include
export CGO_LDFLAGS=-L$(LIBGPIOD)/lib
export GOPATH=/work/_path
export CGO_ENABLED=1
build:
	go build -ldflags="-s -w" -trimpath -buildvcs=false

libgpiod:
	curl -O https://cdn.kernel.org/pub/software/libs/libgpiod/libgpiod-2.1.1.tar.xz
	tar Jxfv libgpiod-2.1.1.tar.xz
	cd libgpiod-2.1.1 && ./configure --host=arm-linux-gnueabi --enable-shared=no --prefix=$(LIBGPIOD) && make && make install

all: libgpiod build
