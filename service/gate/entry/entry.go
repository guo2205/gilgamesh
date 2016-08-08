// entry
package entry

import "fractal/fractal"

type GateEntry struct {
	Fractal *fractal.Fractal
}

func (c *GateEntry) Close() error {
	return nil
}

func (c *GateEntry) Run() error {
	return nil
}

func (c *GateEntry) WritePacket(session uint64, d []byte) error {
	return nil
}
