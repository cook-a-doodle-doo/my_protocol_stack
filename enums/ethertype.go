package enums

type EtherType uint16

const (
	EtherTypeIPv4 EtherType = 0x0800
	EtherTypeIPv6 EtherType = 0x86dd
	EtherTypeARP  EtherType = 0x0806
	EtherTypeRARP EtherType = 0x8635
)

func (e EtherType) Name() string {
	switch e {
	case EtherTypeIPv4:
		return "IPv4"
	case EtherTypeIPv6:
		return "IPv6"
	case EtherTypeARP:
		return "ARP"
	case EtherTypeRARP:
		return "RARP"
	}
	return "UnDef"
}

//https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml
