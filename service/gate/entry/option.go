// option
package entry

type Option struct {
	PerSecondMaxPacket uint32
}

var DefaultOption = &Option{
	PerSecondMaxPacket: 20,
}
