// +build netbsd,arm netbsd,386

package dev

// IOCTL value

const i2c_exec = 0x80184900 // IOCTL number

type I2cparam struct {
	op     uint32    // 4 +00
	addr   uint16    // 2 +04
	_      uint16    // 2 +06
	cmd    *[32]byte // 4 +08
	cmdlen int       // 4 +0C
	buf    *[32]byte // 4 +10
	buflen int       // 4 +14
}
