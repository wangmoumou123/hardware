package modules

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"runtime"
	"strconv"
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
	config := &serial.Config{Name: ser.Port, Baud: ser.Baud, ReadTimeout: ser.ReadTimeout, Size: serial.DefaultSize,
		Parity: ser.Parity, StopBits: serial.Stop1}
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
		ser = &SerialCommunicate{Port: nameOrPort, Baud: baud, ReadTimeout: 0}
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

func scanBytes(startBytes []byte) bufio.SplitFunc {
	//log.Println("startBytes", startBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//log.Println("startBytes", startBytes)
		//log.Println(bytesToHexString(data))
		if len(startBytes) == 0 {
			startBytes = data[0:2]
		}
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, startBytes); i >= 0 && len(data) > 1 {
			return i + len(startBytes), dropCR(data[i:]), nil
		}
		if atEOF {
			return len(data), dropCR(data), nil
		}
		return 0, nil, nil
	}
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
func bytesToHexString(data []byte) string {
	hexString := ""
	for i, b := range data {
		if i > 0 {
			hexString += " "
		}
		hexString += fmt.Sprintf("%02X", b)
	}
	return hexString
}

func hexStringToBytes(hexString string) ([]byte, error) {
	hexString = strings.ReplaceAll(hexString, " ", "")
	if len(hexString)%2 != 0 {
		return nil, errors.New("invalid hex string length")
	}
	bs := make([]byte, len(hexString)/2)
	for i := 0; i < len(hexString); i += 2 {
		hexPair := hexString[i : i+2]
		b, err := strconv.ParseUint(hexPair, 16, 8)
		if err != nil {
			return nil, err
		}
		bs[i/2] = byte(b)
	}
	return bs, nil
}

func (ser *SerialCommunicate) ReadCallback(callback Callback, wg *sync.WaitGroup, hexStarts ...string) {
	var hex bool
	if len(hexStarts) == 0 {
		hex = false
	} else {
		hex = true
	}
	reader := bufio.NewReader(ser.Ser)
	scanner := bufio.NewScanner(reader)
	if hex {
		bs, err := hexStringToBytes(hexStarts[0])
		if err != nil {
			return
		}
		scanner.Split(scanBytes(bs))
	} else {
		scanner.Split(bufio.ScanLines)
	}
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
