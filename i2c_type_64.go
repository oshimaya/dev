// +build netbsd,amd64 netbsd,arm64

package dev

const i2c_exec = 0x80284900 // IOCTL number

// IOCTL argument

type I2cparam struct {
	op     uint32    // 4 +00
	addr   uint16    // 2 +04
	_      uint16    // 2 +06
	cmd    *[32]byte // 8 +08
	cmdlen int       // 8 +10
	buf    *[32]byte // 8 +18
	buflen int       // 8 +20
}
