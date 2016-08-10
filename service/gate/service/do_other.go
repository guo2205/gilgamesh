// do_other
package service

func (c *GateService) doOther(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	switch _type {
	case 0:
		return []byte{}, c.writer(session, data)
	case 1:
	}
	return []byte{}, nil
}
