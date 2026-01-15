package dos

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Time struct {
	Hour   uint16
	Minute uint16
	Second uint16
}

func (t *Time) Get() uint16 {
	return 0 |
		((t.Second >> 1) & 0x1f) |
		((t.Minute & 0x3f) << 5) |
		((t.Hour & 0x1f) << 11)
}

func (t *Time) Set(v uint16) {
	t.Second = (v & 0x1f) << 1
	t.Minute = (v >> 5) & 0x3f
	t.Hour = (v >> 11) & 0x1f
}

func (t *Time) Marshal(w io.WriteSeeker) error {
	return binary.Write(w, binary.LittleEndian, t.Get())
}

func (t *Time) Unmarshal(r io.ReadSeeker) error {
	var i uint16
	if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
		return err
	}
	t.Set(i)
	return nil
}

func (t *Time) SizeOf() uint32 {
	return 2
}

func (t *Time) Stringify() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

func NewTime(hour, minute, second uint16) *Time {
	return &Time{Hour: hour, Minute: minute, Second: second}
}

type Date struct {
	Year  uint16
	Month uint16
	Day   uint16
}

func (d *Date) Get() uint16 {
	return 0 |
		(d.Day & 0x1f) |
		((d.Month & 0xf) << 5) |
		(((d.Year - 1980) & 0x7f) << 9)
}

func (d *Date) Set(v uint16) {
	d.Day = v & 0x1f
	d.Month = (v >> 5) & 0xf
	d.Year = ((v >> 9) & 0x7f) + 1980
}

func (d *Date) Marshal(w io.WriteSeeker) error {
	return binary.Write(w, binary.LittleEndian, d.Get())
}

func (d *Date) Unmarshal(r io.ReadSeeker) error {
	var i uint16
	if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
		return err
	}
	d.Set(i)
	return nil
}

func (d *Date) SizeOf() uint32 {
	return 2
}

func (d *Date) Stringify() string {
	return fmt.Sprintf("%04d/%02d/%02d", d.Year, d.Month, d.Day)
}

func NewDate(year, month, day uint16) *Date {
	return &Date{Year: year, Month: month, Day: day}
}
