// +build netbsd

// Example for Raspberry PI with I2C display AQM0802A
//    http://akizukidenshi.com/catalog/g/gP-09422/
// Configuration:
//   (MAKEDEV iic)
// Run with superuser
// 

package main

import (
	"errors"
	"fmt"
	"github.com/oshimaya/dev"
	"time"
)

func lcd_init(i2c *dev.I2c, i2cp *dev.I2cparam) (err error) {
	init1 := []byte{0x38, 0x39, 0x14, 0x70, 0x56, 0x6c}
	init2 := []byte{0x38, 0x0c, 0x01}
	cmd := []byte{0x00}

	err = i2c.Open()
	if err != nil {
		return err
	}
	i2cp.SetCmd(cmd)
	i2cp.SetAddr(0x3E)
	i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	i2cp.SetData(init1)
	err = i2c.Exec(i2cp)
	if err != nil {
		return
	}
	time.Sleep(300 * time.Millisecond)

	i2cp.SetData(init2)
	err = i2c.Exec(i2cp)

	time.Sleep(300 * time.Millisecond)

	return
}

func lcd_pos(i2c *dev.I2c, i2cp *dev.I2cparam, x int, y int) (err error) {
	cmd := []byte{0x00}

	data := []byte{0x80}
	if x < 8 && y < 2 {
		data[0] = byte(0x80 | x | (y << 7))
		i2cp.SetCmd(cmd)
		i2cp.SetAddr(0x3E)
		i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
		i2cp.SetData(data)
		err = i2c.Exec(i2cp)
	} else {
		err = errors.New("Invalid position")
	}
	return err
}

func lcd_display(i2c *dev.I2c, i2cp *dev.I2cparam, str string) (err error) {

	cmd := []byte{0x40}
	len := len(str)
	if len > 8 {
		len = 8
	}
	data := make([]byte, len)
	for i := 0; i < len; i++ {
		data[i] = str[i]
	}
	i2cp.SetCmd(cmd)
	i2cp.SetAddr(0x3E)
	i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	i2cp.SetData(data)
	err = i2c.Exec(i2cp)
	return err

}

func main() {

	i2c := dev.NewI2c("/dev/iic1") // or /dev/iic{0,2}

	lcd_i2cp := dev.NewI2cparam()

	err := lcd_init(i2c, lcd_i2cp)

	if err != nil {
		return
	}
	defer i2c.Close()

	for i := 0; i < 8; i++ {
		for j:=0; j < 8 ; j++ {
			// Error Retry loop (busy timing?))
			err = lcd_pos(i2c, lcd_i2cp, 7-i, 0)
			if err == nil {
				err = lcd_display(i2c, lcd_i2cp, "NetBSD     ")
				if err == nil {
					break
				}
			}
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}
