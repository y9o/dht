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
	var loop int
	flag.StringVar(&chip, "chip", "/dev/gpiochip0", "/dev/gpiochip0")
	flag.IntVar(&line, "line", -1, "4")
	flag.IntVar(&loop, "loop", 10, "loop 10")
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

	//calculation Limit
	//There is no need to calculate every time.
	readtest := 10000
	dur, err := d.GetReadTime(readtest)
	if err != nil {
		log.Fatalln(err)
	}
	once := dur / time.Duration(readtest)
	d.Limit = uint32(6 * time.Millisecond / once)
	log.Printf("GetReadTime(%dloop)=> %v / %d = %v\n", readtest, dur, readtest, once)
	log.Printf("6000ms / %v = %d\n", once, d.Limit)

	time.Sleep(d.Config.Retry)

	var total, retry, success int
	defer func() {
		if total > 0 {
			fmt.Printf("==========\n")
			fmt.Printf("total  :%4d\n", total)
			fmt.Printf("success:%4d  %7.3f\n", success, float64(success)/float64(total)*100)
			fmt.Printf("retry  :%4d  %7.3f\n", retry, float64(retry)/float64(total)*100)
		}
	}()

	var buf dht.DHT22Data
	for range loop {
		fmt.Print("read...\r")
		re, err := d.Read(&buf, 20)
		total += 1 + re
		retry += re
		if err != nil {
			log.Println("err:", err)
			//d.Dump()
		} else {
			success++
			log.Printf("🌡️%2.1f℃ (%d) 🌢%2.1f%% (%d) retry:%d\n", buf.Temp(), buf.Temperature, buf.Hum(), buf.Humidity, re)
		}
		select {
		case <-sigs:
			return
		case <-time.After(10 * time.Second):
			break
		}
	}
	log.Println("last data")
	d.Dump()
}
