package enums

type HardwareType uint16

const (
	HardwareTypeEthernet HardwareType = 1
	HardwareTypeLoopback HardwareType = 2
)

func (h HardwareType) Name() string {
	switch h {
	case HardwareTypeEthernet:
		return "ethernet"
	case HardwareTypeLoopback:
		return "loopback"
	}
	return "undifined"
}
