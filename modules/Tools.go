package modules

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type ExitCallBack func()

func ExitHandle(callback ExitCallBack) chan struct{} {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		defer callback()
		sig := <-sigChan
		fmt.Printf("捕获到信号: %s\n", sig)
		close(done)
	}()
	return done
}

func BytesToHexString(data []byte) string {
	hexString := ""
	for i, b := range data {
		if i > 0 {
			hexString += " "
		}
		hexString += fmt.Sprintf("%02X", b)
	}
	return hexString
}

func HexStringToBytes(hexString string) ([]byte, error) {
	hexString = strings.ReplaceAll(hexString, " ", "")
	if len(hexString)%2 != 0 {
		return nil, errors.New("invalid hex string length")
	}
	bytes := make([]byte, len(hexString)/2)

	for i := 0; i < len(hexString); i += 2 {
		hexPair := hexString[i : i+2]
		b, err := strconv.ParseUint(hexPair, 16, 8)
		if err != nil {
			return nil, err
		}
		bytes[i/2] = byte(b)
	}

	return bytes, nil
}
