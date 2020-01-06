package config

import (
	"encoding/json"
)

// Instance is the master configuration for the whole alarma application.
type Instance struct {
	// TimeSpansActive is set to true if time-span related actions are to be performed.
	// If false, it is treated as if all TimeSpan entries were disabled.
	TimeSpansActive bool `json:"time-spans-active"`

	// TimeSpans contains the list of all configured time span entries.
	TimeSpans []TimeSpan `json:"time-spans"`
}

// FromBytes decodes a configuration instance from given byte stream.
func FromBytes(data []byte) (inst Instance, err error) {
	err = json.Unmarshal(data, &inst)
	return
}

// ToBytes encodes the configuration instance to a byte stream.
func (inst Instance) ToBytes() ([]byte, error) {
	return json.MarshalIndent(&inst, "", "  ")
}
