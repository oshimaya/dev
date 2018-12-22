// +build netbsd

package dev

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

const (
	GPIOPIN_INPUT     = 0x00000001
	GPIOPIN_OUTPUT    = 0x00000002
	GPIOPIN_INOUT     = 0x00000004
	GPIOPIN_OPENDRAIN = 0x00000008
	GPIOPIN_PUSHPULL  = 0x00000010
	GPIOPIN_TRISTATE  = 0x00000020
	GPIOPIN_PULLUP    = 0x00000040
	GPIOPIN_PULLDOWN  = 0x00000080
	GPIOPIN_INVIN     = 0x00000100
	GPIOPIN_INVOUT    = 0x00000200
	GPIOPIN_USER      = 0x00000400
	GPIOPIN_PULSATE   = 0x00000800
	GPIOPIN_SET       = 0x00008000
	GPIOPIN_ALT0      = 0x00010000
	GPIOPIN_ALT1      = 0x00020000
	GPIOPIN_ALT2      = 0x00040000
	GPIOPIN_ALT3      = 0x00080000
	GPIOPIN_ALT4      = 0x00100000
	GPIOPIN_ALT5      = 0x00200000
	GPIOPIN_ALT6      = 0x00400000
	GPIOPIN_ALT7      = 0x00800000
	GPIOPIN_EVENTS    = 0x10000000
	GPIOPIN_LEVEL     = 0x20000000
	GPIOPIN_FALLING   = 0x40000000
)

const (
	gpioinfo   = 0x40044700
	gpioset    = 0xC08C4705
	gpiounset  = 0xC08C4706
	gpioread   = 0xC0484707
	gpiowrite  = 0xC0484708
	gpiotoggle = 0xC0484709
	gpioattach = 0xC01C470A
)

type Gpio struct {
	dev string
	fd  int
	m   *sync.Mutex
}

type Req struct {
	name  [64]byte
	pin   int32
	value int32
}

type Conf struct {
	name  [64]byte
	pin   int32
	caps  int32
	flags int32
	name2 [64]byte
}

func NewGpio(dev string) *Gpio {
	gpio := new(Gpio)
	gpio.dev = dev
	gpio.fd = -1
	gpio.m = new(sync.Mutex)
	return gpio
}

func (g *Gpio) Open() (err error) {
	if g.fd != -1 {
		errors.New("Already Open")
	}
	g.fd, err = syscall.Open(g.dev, syscall.O_RDWR, 0)
	if err != nil {
		return err
	}
	// Check pin information ?
	return nil
}

func (g *Gpio) Close() {
	if g.fd != -1 {
		syscall.Close(g.fd)
	}
	g.fd = -1
}

func (g *Gpio) ReadPin(pin int) (int, error) {
	var req Req
	req.pin = int32(pin)
	g.m.Lock()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(g.fd),
		gpioread, uintptr(unsafe.Pointer(&req)))
	g.m.Unlock()
	if errno != 0 {
		return 0, errno
	}
	return int(req.value), nil
}

func (g *Gpio) WritePin(pin int, state int) (int, error) {
	var req Req
	req.name[0] = byte(0)
	req.pin = int32(pin)
	req.value = int32(state)
	g.m.Lock()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(g.fd),
		gpiowrite, uintptr(unsafe.Pointer(&req)))
	g.m.Unlock()
	if errno != 0 {
		return 0, errno
	}
	return int(req.value), nil
}
