// bosh bme280 environmental sensor I2C device
// datasheet:http://akizukidenshi.com/download/ds/bosch/BST-BME280_DS001-10.pdf
//
// +build netbsd

package bme280

import (
	"github.com/oshimaya/dev"
	"log"
	"time"
)

type Bme280 struct {
	i2c   *dev.I2c // i2c device
	i2cp  *dev.I2cparam
	Temp  float64
	Press float64
	Hum   float64
	Time  time.Time
	tfine int32
	calib calibdata
	sense [8]byte // raw sense data (not use yet...)
}

// calibration data
type calibdata struct {
	t1 uint16
	t2 int16
	t3 int16
	p1 uint16
	p2 int16
	p3 int16
	p4 int16
	p5 int16
	p6 int16
	p7 int16
	p8 int16
	p9 int16
	h1 uint8
	h2 int16
	h3 uint8
	h4 int16
	h5 int16
	h6 int8
}

func NewBme280(i2c *dev.I2c) *Bme280 {
	bme := new(Bme280)
	bme.i2c = i2c
	bme.i2cp = dev.NewI2cparam()
	return bme
}

// reset
func (bme *Bme280) reset() error {
	data := []byte{0xB6}
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetCmdOne(0xE0)
	bme.i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	bme.i2cp.SetData(data)
	return bme.i2c.Exec(bme.i2cp)
}

func (bme *Bme280) setCtrlMeas() error {
	data := []byte{0x27} // temp:x1, pressure:x1
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetCmdOne(0xF4)
	bme.i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	bme.i2cp.SetData(data)
	return bme.i2c.Exec(bme.i2cp)
}

func (bme *Bme280) setCtrlHum() error {
	data := []byte{0x01}
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetCmdOne(0xF2)
	bme.i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	bme.i2cp.SetData(data)
	return bme.i2c.Exec(bme.i2cp)
}

func (bme *Bme280) setConfig() error {
	data := []byte{0xA0}
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetCmdOne(0xF5)
	bme.i2cp.SetOp(dev.I2C_OP_WRITE_WITH_STOP)
	bme.i2cp.SetData(data)
	return bme.i2c.Exec(bme.i2cp)
}

func (bme *Bme280) getCarribData() error {
	data := make([]byte, 26)
	data1 := make([]byte, 7)
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetOp(dev.I2C_OP_READ_WITH_STOP)

	bme.i2cp.SetCmdOne(0x88)
	bme.i2cp.SetData(data)
	err := bme.i2c.Exec(bme.i2cp)
	if err != nil {
		return err
	}
	data = bme.i2cp.GetData()

	bme.i2cp.SetCmdOne(0xE1)
	bme.i2cp.SetData(data1)
	err = bme.i2c.Exec(bme.i2cp)
	if err != nil {
		return err
	}
	data1 = bme.i2cp.GetData()

	bme.calib.t1 = uint16(data[0]) | uint16(data[1])<<8
	bme.calib.t2 = int16(data[2]) | int16(data[3])<<8
	bme.calib.t3 = int16(data[4]) | int16(data[5])<<8
	bme.calib.p1 = uint16(data[6]) | uint16(data[7])<<8
	bme.calib.p2 = int16(data[8]) | int16(data[9])<<8
	bme.calib.p3 = int16(data[10]) | int16(data[11])<<8
	bme.calib.p4 = int16(data[12]) | int16(data[13])<<8
	bme.calib.p5 = int16(data[14]) | int16(data[15])<<8
	bme.calib.p6 = int16(data[16]) | int16(data[17])<<8
	bme.calib.p7 = int16(data[18]) | int16(data[19])<<8
	bme.calib.p8 = int16(data[20]) | int16(data[21])<<8
	bme.calib.p9 = int16(data[22]) | int16(data[23])<<8

	bme.calib.h1 = uint8(data[25])
	bme.calib.h2 = int16(data1[0]) | int16(data1[1])<<8
	bme.calib.h3 = uint8(data1[2])
	bme.calib.h4 = int16(data1[3])<<4 | int16(data1[4])&0x0F
	bme.calib.h4 = int16(data1[4])>>4 | int16(data1[5])<<4
	bme.calib.h6 = int8(data1[6])

	return nil
}

func (bme *Bme280) getTempRaw() (int32, error) {
	data := make([]byte, 3)
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetOp(dev.I2C_OP_READ_WITH_STOP)
	bme.i2cp.SetCmdOne(0xFA)
	bme.i2cp.SetData(data)
	err := bme.i2c.Exec(bme.i2cp)
	if err != nil {
		return 0, err
	}
	data = bme.i2cp.GetData()
	x := int32(data[0])<<12 | int32(data[1])<<4 | int32(data[2])>>4
	return x, nil
}

func (bme *Bme280) getPressRaw() (int32, error) {
	data := make([]byte, 3)
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetOp(dev.I2C_OP_READ_WITH_STOP)
	bme.i2cp.SetCmdOne(0xF7)
	bme.i2cp.SetData(data)
	err := bme.i2c.Exec(bme.i2cp)
	if err != nil {
		return 0, err
	}
	data = bme.i2cp.GetData()
	x := int32(data[0])<<12 | int32(data[1])<<4 | int32(data[2])>>4
	return x, nil
}

func (bme *Bme280) getHumRaw() (int32, error) {
	data := make([]byte, 2)
	bme.i2cp.SetAddr(0x76)
	bme.i2cp.SetOp(dev.I2C_OP_READ_WITH_STOP)
	bme.i2cp.SetCmdOne(0xFD)
	bme.i2cp.SetData(data)
	err := bme.i2c.Exec(bme.i2cp)
	if err != nil {
		return 0, err
	}
	data = bme.i2cp.GetData()
	x := int32(data[0])<<8 | int32(data[1])
	return x, nil
}

func (bme *Bme280) Init() {
	bme.reset()
	time.Sleep(300 * time.Millisecond)
	bme.setCtrlMeas()
	time.Sleep(100 * time.Millisecond)
	bme.setCtrlHum()
	time.Sleep(100 * time.Millisecond)
	bme.setConfig()
	time.Sleep(100 * time.Millisecond)
	bme.getCarribData()
}

func (bme *Bme280) getTemp() (float64, error) {
	traw, err := bme.getTempRaw()
	if err != nil {
		return 0, err
	}
	n1 := ((traw>>3 - int32(bme.calib.t1)<<1) * int32(bme.calib.t2)) >> 11
	n2 := ((((traw>>4 - int32(bme.calib.t1)) *
		(traw>>4 - int32(bme.calib.t1))) >> 12) *
		int32(bme.calib.t3)) >> 14
	bme.tfine = n1 + n2
	temp := float64((bme.tfine*5+128)>>8) / 100
	return temp, nil
}

func (bme *Bme280) getPress() (float64, error) {
	praw, err := bme.getPressRaw()
	if err != nil {
		return 0, err
	}
	p1 := int64(bme.tfine) - 128000
	p2 := p1 * p1 * int64(bme.calib.p6)
	p2 += (p1 * int64(bme.calib.p5)) << 17
	p2 += int64(bme.calib.p4) << 35
	p1 = (p1*p1*int64(bme.calib.p3))>>8 +
		((p1 * int64(bme.calib.p2)) << 12)
	p1 = ((int64(1)<<47 + p1) * int64(bme.calib.p1)) >> 33
	if p1 == 0 {
		return 0, nil
	}
	p := ((((1048576 - int64(praw)) << 31) - p2) * 3125) / p1
	p1 = (int64(bme.calib.p9) * (p >> 13) * (p >> 13)) >> 25
	p2 = (int64(bme.calib.p8) * p) >> 19
	x := uint32(((p + p1 + p2) >> 8) + (int64(bme.calib.p7) << 4))
	return float64(x) / 256 / 100, nil

}

// XXX: not work yet properly (???)
//
func (bme *Bme280) getHum() (float64, error) {
	hraw, err := bme.getHumRaw()
	if err != nil {
		return 0, err
	}
	h1 := bme.tfine - 76800
	h2 := hraw<<14 - int32(bme.calib.h4)<<20 - int32(bme.calib.h5)*h1
	h3 := (h2 + 16384) >> 15
	h4 := (h1 * int32(bme.calib.h6)) >> 10
	h5 := (h1 * int32(bme.calib.h3)) >> 11
	h6 := (h4 * (h5 + 32768)) >> 10
	h7 := ((h6+2097152)*int32(bme.calib.h2) + 8192) >> 14
	h1 = h3 * h7
	h1 -= ((((h1 >> 15) * (h1 >> 15)) >> 7) * int32(bme.calib.h1)) >> 4
	if h1 < 0 {
		return 0, nil
	}
	if h1 > 419430400 {
		h1 = 419430400
	}
	return float64(h1>>12) / 1024.0, nil
}

func (bme *Bme280) SenseNow() {
	var err error
	bme.Time = time.Now()
	bme.Temp, err = bme.getTemp()
	if err != nil {
		log.Println("Error: Temp: ", err)
		time.Sleep(1 * time.Second)
		bme.Init()
		return
	}
	bme.Press, err = bme.getPress()
	if err != nil {
		log.Println("Error: Press: ", err)
		time.Sleep(1 * time.Second)
		bme.Init()
		return
	}
	bme.Hum, err = bme.getHum()
	if err != nil {
		log.Println("Error: Hum: ", err)
		time.Sleep(1 * time.Second)
		bme.Init()
		return
	}
}
