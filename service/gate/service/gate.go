// gate
package service

import "fractal/fractal"

type GateService struct {
	Fractal     *fractal.Fractal
	WritePacket func(session uint64, d []byte) error
}

func (c *GateService) OnStart() error {
	return nil
}

func (c *GateService) OnMail(caller string, session uint64, data []byte) ([]byte, error) {
	return nil, nil
}

func (c *GateService) OnClose() error {
	return nil
}
