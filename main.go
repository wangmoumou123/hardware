package main

import (
	"bytes"
	"log"
	"strings"
)

var FLAGS = []string{"at", "+", "@gps", "@batt"}

func callback(m []byte) {
	msg := string(m)
	for _, f := range FLAGS {

		if strings.Contains(msg, f) {
			//log.Println("string_callback===>", msg)
			break
		}
		if bytes.Contains(m, []byte(f)) {
			//log.Println("bytes_callback===>", msg)
			break
		}

	}

}

func msgCallback(msg string, tp string) {
	log.Println(tp, "======>", msg)
}

func main() {
	//var wg sync.WaitGroup
	////ser := modules.SerialCommunicateInit("COM128", 0, 115200)
	//ser := modules.SerialCommunicateInit("FTDI", 2, 115200)
	//ser.ReadCallback(callback, &wg)
	//time.Sleep(time.Second)
	//log.Println("===start send===")
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for {
	//		ser.Send("at+csq\r\n")
	//		ser.Send("at+cpsi\r\n")
	//		time.Sleep(time.Second)
	//	}
	//}()
	//topics := []string{"ws"}
	//mqttConn := modules.MqttConnInit("ws://192.168.0.56:9001", topics, msgCallback)
	//wg.Add(1)
	//go mqttConn.RunAlways(&wg)
	//
	//modules.ExitHandle(func() {
	//	mqttConn.Conn.Disconnect(100)
	//	ser.Close()
	//})
	//wg.Wait()

	//pin := modules.GpioInit(4, "OUT")
	//
	//for {
	//	time.Sleep(time.Second)
	//	pin.Write(1) // Set pin High
	//	time.Sleep(time.Second)
	//	pin.Write(0) // Set pin Low
	//}

	//u := modules.UdpInit("127.0.0.1", 9900)
	//u.Send("dasda")

}
