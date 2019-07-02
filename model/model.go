package model

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"
)

type XDSData struct {
	sync.Mutex
	resources []byte
}

func New() *XDSData {
	return &XDSData{}
}

func (d *XDSData) SetResources(res []interface{}, logger zerolog.Logger) {
	d.Lock()
	defer d.Unlock()

	if res == nil {
		d.resources = make([]byte, 0)
		return
	}
	resources, err := json.Marshal(res)
	if err != nil {
		logger.Error().AnErr("json.Marshal", err).Msg("Could not marshal into JSON")
		d.resources = make([]byte, 0)
		return
	}
	d.resources = resources
}

func (d *XDSData) GetResources() []byte {
	d.Lock()
	defer d.Unlock()

	return d.resources
}
