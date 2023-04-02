package device

import (
	"fmt"

	"github.com/davecheney/i2c"
)

type I2cHt16k33 struct {
	i2c    *i2c.I2C
	buffer [10]uint16
}

func NewI2cHt16k33(i2c *i2c.I2C) *I2cHt16k33 {
	return &I2cHt16k33{
		i2c: i2c,
	}
}

func (t *I2cHt16k33) SetBrightness(brightness int) error {
	if err := t.WriteRaw([]byte{0xe0 + byte(brightness)}); err != nil {
		return fmt.Errorf("set brightness: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) DisplayOn() error {
	if err := t.WriteRaw([]byte{0x81}); err != nil {
		return fmt.Errorf("display on: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) OscillatorOn() error {
	if err := t.WriteRaw([]byte{0x21}); err != nil {
		return fmt.Errorf("oscillator on: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) WriteRaw(data []byte) error {
	if _, err := t.i2c.Write(data); err != nil {
		return fmt.Errorf("write bytes: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) WriteData() error {
	data := make([]byte, len(t.buffer))

	i := 0
	for _, item := range t.buffer {
		data[i] = byte(item)
		i++
	}

	if err := t.WriteRaw(data); err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) Clear() error {
	for i := range t.buffer {
		t.buffer[i] = uint16(0)
	}

	if err := t.WriteData(); err != nil {
		return fmt.Errorf("clear: %w", err)
	}

	return nil
}

func (t *I2cHt16k33) SetSegments(humanPos int, on [7]bool) {
	t.buffer[TranslateHumanPosition(humanPos)] = t.ArrayToSegments(on)
}

func (t *I2cHt16k33) SetDigit(position int, digit int) {
	t.buffer[TranslateHumanPosition(position)] = t.IntToSevenSegement(digit)
}

func (t *I2cHt16k33) ArrayToSegments(segmentsOn [7]bool) uint16 {
	tot := uint16(0)
	for index, on := range segmentsOn {
		if on {
			tot += 1 << index
		}
	}
	return tot
}

func (t *I2cHt16k33) IntToSevenSegement(digit int) uint16 {
	values := []uint16{63, 6, 91, 79, 102, 109, 125, 7, 127, 103}
	if digit == -10000 {
		return 0
	} else if digit < 0 {
		return 64
	}
	return values[digit]
}

func (t *I2cHt16k33) AlphaToSegments(letter byte) uint16 {
	var alphamap map[string]uint16 = make(map[string]uint16)

	alphamap["A"] = t.ArrayToSegments([7]bool{true, true, true, false, true, true, true})
	alphamap["b"] = t.ArrayToSegments([7]bool{false, false, true, true, true, true, true})
	alphamap["c"] = t.ArrayToSegments([7]bool{false, false, false, true, true, false, true})
	alphamap["d"] = t.ArrayToSegments([7]bool{false, true, true, true, true, false, true})
	alphamap["e"] = t.ArrayToSegments([7]bool{true, true, false, true, true, true, true})
	alphamap["f"] = t.ArrayToSegments([7]bool{true, false, false, false, true, true, true})
	alphamap["g"] = t.ArrayToSegments([7]bool{true, true, true, true, false, true, true})
	alphamap["h"] = t.ArrayToSegments([7]bool{false, false, true, false, true, true, true})
	alphamap["I"] = t.ArrayToSegments([7]bool{false, false, false, false, true, true, false})
	alphamap["o"] = t.ArrayToSegments([7]bool{false, false, true, true, true, false, true})
	alphamap["u"] = t.ArrayToSegments([7]bool{false, false, true, true, true, false, false})
	alphamap["S"] = t.ArrayToSegments([7]bool{true, false, true, true, false, true, true})

	return alphamap[string(letter)]
}

func (t *I2cHt16k33) SetAlpha(position int, letter byte) {
	t.buffer[TranslateHumanPosition(position)] = t.AlphaToSegments(letter)
}

func (t *I2cHt16k33) SetNumber(number int) error {
	if number > 9999 {
		return fmt.Errorf("number must be between 0 and 9999: %d", number)
	}

	if number < 0 {
		for i := 0; i <= 2; i++ {
			t.SetDigit(i, -1)
		}
	} else {
		num_a := NumToArray(number)
		for i, j := range num_a {
			t.SetDigit(i, j)
		}
	}

	return nil
}

func NumToArray(number int) []int {
	var num_a []int
	if number == 0 {
		num_a = append(num_a, 0)
	} else {
		num_tmp := number
		for num_tmp != 0 {
			num_a = append(num_a, num_tmp%10)
			num_tmp = num_tmp / 10
		}
	}
	for len(num_a) < 4 {
		num_a = append(num_a, -10000)
	}
	return num_a
}

func TranslateHumanPosition(humanPosition int) int {
	positions := [5]int{9, 7, 3, 1, 5}
	return positions[humanPosition]
}
