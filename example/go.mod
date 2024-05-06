module dht22

go 1.22.0

require (
	github.com/warthog618/go-gpiocdev v0.9.0
	github.com/y9o/dht v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.18.0 // indirect

replace github.com/y9o/dht => ../
