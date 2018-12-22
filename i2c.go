// +build netbsd

package dev

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

const (
	I2C_OP_READ            = 0
	I2C_OP_READ_WITH_STOP  = 1
	I2C_OP_WRITE           = 2
	I2C_OP_WRITE_WITH_STOP = 3
	I2C_OP_READ_BLOCK      = 5
	I2C_OP_WRITE_BLOCK     = 7
)

type I2c struct {
	dev string
	fd  int
	m   *sync.Mutex
}

func NewI2c(dev string) *I2c {
	i2c := new(I2c)
	i2c.dev = dev
	i2c.fd = -1
	return i2c
}

func NewI2cparam() *I2cparam {
	i2cp := new(I2cparam)
	cmd := new([32]byte)
	i2cp.cmd = cmd
	data := new([32]byte)
	i2cp.buf = data
	return i2cp
}

func (i2c *I2c) Open() (err error) {
	if i2c.fd != -1 {
		return errors.New("Already Open")
	}
	i2c.fd, err = syscall.Open(i2c.dev, syscall.O_RDWR, 0)
	if err != nil {
		return err
	}
	i2c.m = new(sync.Mutex)
	return nil
}

func (i2c *I2c) Close() {
	syscall.Close(i2c.fd)
	i2c.fd = -1
	i2c.m = nil
}

func (i2c *I2c) Exec(i2cp *I2cparam) error {
	i2c.m.Lock()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(i2c.fd),
		i2c_exec, uintptr(unsafe.Pointer(i2cp)))
	i2c.m.Unlock()
	if errno != 0 {
		return errno
	}
	return nil
}

func (i2cp *I2cparam) SetAddr(addr int) {
	i2cp.addr = uint16(addr)
}

func (i2cp *I2cparam) SetCmd(cmd []byte) {
	if len(cmd) > 0 && len(cmd) <= 32 {
		i2cp.cmdlen = len(cmd)
		for i := range cmd {
			i2cp.cmd[i] = cmd[i]
		}
	}
}

func (i2cp *I2cparam) SetOp(op int) {
	i2cp.op = uint32(op)
}

func (i2cp *I2cparam) SetData(data []byte) {
	if len(data) > 0 && len(data) <= 32 {
		i2cp.buflen = len(data)
		for i := range data {
			i2cp.buf[i] = data[i]
		}
	}
}

func (i2cp *I2cparam) GetData() []byte {
	if i2cp.buflen > 0 && i2cp.buflen <= 32 {
		data := make([]byte, i2cp.buflen)
		for i := range data {
			data[i] = i2cp.buf[i]
		}
		return data
	}
	return nil
}
