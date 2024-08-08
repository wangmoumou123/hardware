package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hardware/modules"
	"log"
	"sync"
	"time"
)

var FLAGS = []string{"at", "+", "@gps", "@batt"}

func callback(m []byte) {
	log.Println("hex_ callback===>", modules.BytesToHexString(m), "===", len(m))
	// 从字节数组中提取第 4 到第 11 字节（Go 中的切片是左闭右开）
	data := m[3:11]

	// 创建一个 bytes.Reader 以便从中读取数据
	reader := bytes.NewReader(data)

	// 读取两个小端序的无符号 32 位整数
	var x, y uint32
	err := binary.Read(reader, binary.LittleEndian, &x)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		return
	}
	err = binary.Read(reader, binary.LittleEndian, &y)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		return
	}

	//// 打印读取的值
	//fmt.Printf("X Axis Data: %d\n", x)
	//fmt.Printf("Y Axis Data: %d\n", y)

	// 计算角度
	xAngle := 0.01 * (float64(x) - 9000)
	yAngle := 0.01 * (float64(y) - 9000)

	// 打印角度
	fmt.Printf("X ===> %.2f° Y ===> %.2f°\n", xAngle, yAngle)

	//msg := string(m)
	//for _, f := range FLAGS {
	//
	//	if strings.Contains(msg, f) {
	//		//log.Println("string_callback===>", msg)
	//		break
	//	}
	//	if bytes.Contains(m, []byte(f)) {
	//		//log.Println("bytes_callback===>", msg)
	//		break
	//	}

	//}

}

//func msgCallback(msg string, tp string) {
//	log.Println(tp, "======>", msg)
//}

func main() {
	var wg sync.WaitGroup
	ser := modules.SerialCommunicateInit("COM8", 0, 9600)
	//ser := modules.SerialCommunicateInit("FTDI", 2, 115200)
	ser.ReadCallback(callback, &wg, "")
	//time.Sleep(time.Second)
	//log.Println("===start send===")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// HEX 数据
			hexData := "010300020004E5C9"
			query, _ := hex.DecodeString(hexData)
			//log.Println(modules.BytesToHexString(query))
			ser.SendHex(query)
			//ser.Send("at")
			//read := ser.Read(1024)
			//log.Println(string(read))
			time.Sleep(time.Second)
		}
	}()
	//topics := []string{"ws"}
	//mqttConn := modules.MqttConnInit("ws://192.168.0.56:9001", topics, msgCallback)
	//wg.Add(1)
	//go mqttConn.RunAlways(&wg)
	//
	modules.ExitHandle(func() {
		//mqttConn.Conn.Disconnect(100)
		ser.Close()
	})
	wg.Wait()

	//pin := modules.GpioInit(4, "OUT")
	//
	//for {
	//	time.Sleep(time.Second)
	//	pin.Write(1) // Set pin High
	//	time.Sleep(time.Second)
	//	pin.Write(0) // Set pin Low
	//}

	//u := modules.UdpInit("127.0.0.1", 9900)
	//u.RecvAlways()
}
