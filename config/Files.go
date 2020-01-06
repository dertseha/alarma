package config

import (
	"io/ioutil"
)

// FromFile reads given file and returns the corresponding configuration instance, if successful.
func FromFile(filename string) (Instance, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Instance{}, err
	}
	return FromBytes(data)
}

// ToFile writes given instance to the specified file.
func ToFile(filename string, inst Instance) error {
	data, err := inst.ToBytes()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0600)
}
