#!/bin/bash
#
# for Raspberry Pi Zero
#
docker run -it --rm -v $PWD:/work -w /work \
	docker.elastic.co/beats-dev/golang-crossbuild:1.22.2-armel-debian12 \
	-p "linux/armv6" \
	--build-cmd "make all"
