package main

import (
	"bytes"
	"hardware/modules"
	"log"
	"strings"
	"sync"
	"time"
)

var FLAGS = []string{"at", "+", "@gps", "@batt"}

func callback(m []byte) {
	msg := string(m)
	for _, f := range FLAGS {

		if strings.Contains(msg, f) {
			log.Println("string_callback===>", msg)
			break
		}
		if bytes.Contains(m, []byte(f)) {
			log.Println("bytes_callback===>", msg)
			break
		}

	}

}

func main() {
	var wg sync.WaitGroup
	ser := modules.SerialCommunicateInit("COM128", 115200)
	//command := modules.Command{}
	//sta, out := command.RunCmD("ipconfig")
	//fmt.Println(sta, "===", out)
	ser.ReadCallback(callback, &wg)
	time.Sleep(time.Second)
	log.Println("===start send===")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ser.Send("at+csq\r\n")
			ser.Send("at+cpsi\r\n")
			time.Sleep(time.Second)
		}
	}()
	modules.ExitHandle(func() {
		ser.Close()
	})
	wg.Wait()
}
