package modules

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Callback func([]byte)

type SerialCommunicate struct {
	Port        string `json:"port"`
	Baud        int
	ReadTimeout time.Duration
	Size        byte
	Parity      serial.Parity
	StopBits    serial.StopBits
	Ser         *serial.Port
}

func (ser *SerialCommunicate) openSerial() {
	config := &serial.Config{Name: ser.Port, Baud: ser.Baud, ReadTimeout: ser.ReadTimeout, Size: ser.Size,
		Parity: ser.Parity, StopBits: ser.StopBits}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	ser.Ser = port
}

func SerialCommunicateInit(nameOrPort string, num, baud int) *SerialCommunicate {
	_portFlag := []string{"COM", "/dev/tty"}
	mode := false
	for _, f := range _portFlag {
		if strings.Contains(nameOrPort, f) {
			mode = true
			break
		}
	}

	var ser *SerialCommunicate
	if mode {
		ser = &SerialCommunicate{Port: nameOrPort, Baud: baud}
	} else {
		if IsLinux() {
			ser = &SerialCommunicate{Port: ser.FindPort(nameOrPort, num), Baud: baud}
		} else {
			panic(any("标识符模式目前仅支持LINUX!"))
		}
	}
	ser.openSerial()

	return ser
}

func (ser *SerialCommunicate) Send(msg string) {
	// 写数据到串口
	_, err := ser.Ser.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func (ser *SerialCommunicate) SendHex(msg []byte) {
	_, err := ser.Ser.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (ser *SerialCommunicate) Read(length ...int) []byte {
	bufLen := 0
	if len(length) == 0 {
		bufLen = 128
	} else {
		bufLen = length[0]
	}
	buf := make([]byte, bufLen)
	n, err := ser.Ser.Read(buf)
	if err != nil {
		log.Fatal(err)
		return buf[:n]
	}
	return buf[:n]
}

func (ser *SerialCommunicate) ReadCallback(callback Callback, wg *sync.WaitGroup) {

	// 创建一个缓冲读取器
	// 设置分割函数为按行读取
	// 读取并打印每一行
	reader := bufio.NewReader(ser.Ser)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	wg.Add(1)
	go func() {
		defer func() {
			ser.Close()
			wg.Done()
		}()
		for scanner.Scan() {
			line := scanner.Bytes()

			if len(line) == 0 {
				continue
			}
			callback(line)
		}
		// 检查是否有读取错误
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				fmt.Println("read success")
			} else {
				log.Fatalf("error: %v", err)
			}
		}
	}()
}

func (ser *SerialCommunicate) Close() {
	err := ser.Ser.Close()
	log.Printf("串口%s已关闭!", ser.Port)
	if err != nil {
		log.Fatal(err)
		return
	}

}
func IsLinux() bool {
	os := runtime.GOOS

	switch os {
	case "windows":
		return false
	case "linux":
		return true
	default:
		return false
	}
}

func (ser *SerialCommunicate) FindPort(name string, num int) string {
	defer func() {
		if e := recover(); e != any(nil) {
			log.Fatal(e)
			return
		}
	}()
	var ttyUSb string
	status, output := getStatusOutput("dmesg | grep ttyUSB*")
	if status != 0 {
		panic(any("该操作系统不支持"))
	}
	var tN int
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, name) && strings.Contains(line, "attached") {
			if tN <= num {
				tN++
				ttyUSb = strings.TrimSpace(strings.Split(line, "ttyUSB")[1])
				continue
			}
		} else if strings.Contains(line, name) &&
			strings.Contains(line, "disconnected") &&
			strings.Contains(line, "tty") {
			tN = 0
		}
	}
	if ttyUSb == "" {
		return ""
	}
	log.Println("/dev/ttyUSB" + ttyUSb)
	return "/dev/ttyUSB" + ttyUSb

}
