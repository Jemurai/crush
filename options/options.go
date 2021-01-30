package options

import (
	"encoding/json"
)

// Options is the information we need about a particular run
type Options struct {
	Directory string  `json:"directory"`
	Tag       string  `json:"tag"`
	Ext       string  `json:"ext"`
	Compare   string  `json:"compare"`
	Threshold float64 `json:"threshold"`
}

func (o Options) String() string {
	s, _ := json.Marshal(o)
	return string(s)
}
