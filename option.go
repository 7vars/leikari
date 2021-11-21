package leikari

import (
	"fmt"
)

type Option struct {
	Name string
	Value interface{}
}

func (opt Option) String() string {
	return fmt.Sprintf("%v", opt.Value)
}

func (opt Option) Int() (int, bool) {
	v, ok := opt.Value.(int)
	return v, ok
}

func (opt Option) Int8() (int8, bool) {
	v, ok := opt.Value.(int8)
	return v, ok
}

func (opt Option) Int16() (int16, bool) {
	v, ok := opt.Value.(int16)
	return v, ok
}

func (opt Option) Int32() (int32, bool) {
	v, ok := opt.Value.(int32)
	return v, ok
}

func (opt Option) Int64() (int64, bool) {
	v, ok := opt.Value.(int64)
	return v, ok
}

func (opt Option) Uint() (uint, bool) {
	v, ok := opt.Value.(uint)
	return v, ok
}

func (opt Option) Uint8() (uint8, bool) {
	v, ok := opt.Value.(uint8)
	return v, ok
}

func (opt Option) Uint16() (uint16, bool) {
	v, ok := opt.Value.(uint16)
	return v, ok
}

func (opt Option) Uint32() (uint32, bool) {
	v, ok := opt.Value.(uint32)
	return v, ok
}

func (opt Option) Uint64() (uint64, bool) {
	v, ok := opt.Value.(uint64)
	return v, ok
}

func (opt Option) Float32() (float32, bool) {
	v, ok := opt.Value.(float32)
	return v, ok
}

func (opt Option) Float64() (float64, bool) {
	v, ok := opt.Value.(float64)
	return v, ok
}

func (opt Option) Bool() bool {
	return opt.Value.(bool)
}