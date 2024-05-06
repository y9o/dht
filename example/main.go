package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/y9o/dht"

	"github.com/warthog618/go-gpiocdev"
)

func main() {
	log.Default().SetFlags(log.Ltime)
	var chip string
	var line int
	flag.StringVar(&chip, "chip", "", "gpiochip0")
	flag.IntVar(&line, "line", -1, "4")
	flag.Parse()
	if chip == "" {
		for _, chipname := range gpiocdev.Chips() {
			fmt.Printf("%s -chip %s\n", os.Args[0], chipname)
		}
		return
	}
	if line < 0 {
		c, err := gpiocdev.NewChip(chip)
		if err != nil {
			log.Fatalln(err)
		}
		defer c.Close()
		fmt.Printf("%s -chip %s\n", os.Args[0], chip)
		for i := range c.Lines() {
			info, err := c.LineInfo(i)
			if err != nil {
				log.Fatalf("%d:%s\n", i, err)
			}
			fmt.Printf("\t\t-line %d\t\t(%s)\n", i, info.Name)
		}
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	d, err := dht.New(chip, line)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	var buf dht.DHT22Data

	for range 10 {
		fmt.Print("read...\r")
		if retry, err := d.Read(&buf, 20); err != nil {
			log.Print(err)
		} else {
			log.Printf("🌡️%2.1f℃ (%d) 🌢%2.1f%% (%d) retry:%d\n", buf.Temp(), buf.Temperature, buf.Hum(), buf.Humidity, retry)
		}
		select {
		case <-sigs:
			return
		case <-time.After(10 * time.Second):
			break
		}
	}
}
