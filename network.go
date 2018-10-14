package knockttp

import (
	"encoding/json"
	"io"
)

type Host string

type Services map[Host]Routes

type Network struct {
	Services `json:""`
}

func NewNetworkFromReader(r io.Reader) (*Network, error) {
	var network Network
	dec := json.NewDecoder(r)
	if err := dec.Decode(&network); err != nil {
		return nil, err
	}
	return &network, nil
}
