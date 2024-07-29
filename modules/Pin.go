package modules

import (
	"github.com/stianeikeland/go-rpio/v4"
	"golang.org/x/exp/slices"
	"log"
)

type Gpio interface {
	Read() int
	Write(value rpio.State)
}
type gpio struct {
	Pin  int
	Gpio rpio.Pin
	Mode string
}

func GpioInit(pin int, mode string) Gpio {
	defer func() {
		err := rpio.Close()
		if err != nil {
			return
		}
		if err := recover(); err != nil {
			return
		}
	}()
	if slices.Contains([]string{"IN", "OUT"}, mode) {
		panic("params error")
	}
	g := &gpio{
		Pin:  pin,
		Mode: mode,
	}

	err := rpio.Open()
	if err != nil {
		return nil
	}
	g.Gpio = rpio.Pin(g.Pin)

	if mode == "OUT" {
		g.Gpio.Output()
	} else if mode == "IN" {
		g.Gpio.Input()
	}
	return g
}

func (g *gpio) Read() int {
	return int(g.Gpio.Read())
}
func (g *gpio) Write(value rpio.State) {
	g.Gpio.Write(value)
}

type PWM interface {
	SetFrq(frq int)
	SetDutyCycle(dutyLen, cycleLen int)
	StopAll()
}
type pwm struct {
	Pin      int
	Pwm      rpio.Pin
	Frq      int
	DutyLen  int
	CycleLen int
}

func PwmInit(pin, frq int) PWM {
	defer func() {
		//err := rpio.Close()
		//if err != nil {
		//	return
		//}
		if err := recover(); err != nil {
			return
		}
	}()
	err := rpio.Open()
	if err != nil {
		panic("error")
	}
	p := &pwm{
		Pin:      pin,
		Frq:      frq,
		DutyLen:  0,
		CycleLen: 100,
	}
	p.Pwm = rpio.Pin(p.Pin)
	p.Pwm.Pwm()
	p.SetDutyCycle(p.DutyLen, p.CycleLen)
	log.Println("here")
	p.SetFrq(p.Frq)
	log.Println("here2")
	return p
}

func (p *pwm) SetFrq(frq int) {
	p.Pwm.Freq(frq)
}
func (p *pwm) SetDutyCycle(dutyLen, cycleLen int) {
	p.Pwm.DutyCycle(uint32(dutyLen), uint32(cycleLen))
}

func (p *pwm) StopAll() {
	err := rpio.Close()
	if err != nil {
		return
	}
}
