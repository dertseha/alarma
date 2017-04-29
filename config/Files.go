package config

import (
	"io/ioutil"
)

// FromFile reads given file and returns the corresponding configuration instance, if successful.
func FromFile(filename string) (inst Instance, err error) {
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err == nil {
		inst, err = FromBytes(data)
	}
	return
}

// ToFile writes given instance to the specified file.
func ToFile(filename string, inst Instance) (err error) {
	var data []byte
	data, err = inst.ToBytes()
	if err == nil {
		err = ioutil.WriteFile(filename, data, 0600)
	}
	return
}
