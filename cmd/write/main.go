package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/davecheney/i2c"
	device "github.com/toxygene/i2c-ht16k33"
)

func main() {
	number := flag.Int("number", 0, "the number to display")
	bus := flag.Int("bus", 0, "I2C bus number")
	hexAddress := flag.String("address", "", "hex address of the HT16K33 on the I2C bus")

	flag.Parse()

	address, err := strconv.ParseInt(*hexAddress, 16, 8)
	if err != nil {
		println(fmt.Errorf("parse hex address: %w", err).Error())
		os.Exit(1)
	}

	i2c, err := i2c.New(uint8(address), *bus)
	if err != nil {
		println(fmt.Errorf("connect to i2c device: %w", err).Error())
		os.Exit(1)
	}

	defer i2c.Close()

	ht16k33 := device.NewI2cHt16k33(i2c)

	ht16k33.DisplayOn()
	ht16k33.OscillatorOn()
	ht16k33.SetNumber(*number)
	ht16k33.WriteData()
}
