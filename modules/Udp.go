package modules

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Udp interface {
	Bind() error
	RecvOnce() (string, error)
	RecvAlways(callback func(data []string))
	Send(message string, targetAddr string, targetPort int) error
	SendSeconds(message string, targetAddr string, targetPort int, seconds int) error
	Close()
}

type udp struct {
	Addr         *net.UDPAddr
	Conn         *net.UDPConn
	LastSendTime time.Time
	Mutex        sync.Mutex
}

func UdpInit(addr string, port int) Udp {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil
	}
	u := &udp{
		Addr: udpAddr,
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil
	}
	u.Conn = conn
	return u
}

func (u *udp) Bind() error {
	var err error
	u.Conn, err = net.ListenUDP("udp", u.Addr)
	return err
}

func (u *udp) RecvOnce() (string, error) {
	buffer := make([]byte, 1024)
	n, addr, err := u.Conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}
	data := string(buffer[:n])
	fmt.Printf("Received message from %s: %s\n", addr.String(), data)
	return data, nil
}

func (u *udp) RecvAlways(callback func(data []string)) {
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := u.Conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading from UDP:", err)
				continue
			}
			data := string(buffer[:n])
			fmt.Printf("Received data: %s\n", data)
			if callback != nil {
				callback([]string{data})
			}
		}
	}()
}

func (u *udp) Send(message string, targetAddr string, targetPort int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", targetAddr, targetPort))
	if err != nil {
		return err
	}
	_, err = u.Conn.WriteToUDP([]byte(message), udpAddr)
	return err
}

func (u *udp) SendSeconds(message string, targetAddr string, targetPort int, seconds int) error {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	now := time.Now()
	if now.Sub(u.LastSendTime).Seconds() > float64(seconds) {
		err := u.Send(message, targetAddr, targetPort)
		if err != nil {
			return err
		}
		u.LastSendTime = now
	}
	return nil
}

func (u *udp) Close() {
	err := u.Conn.Close()
	if err != nil {
		return
	}
}
