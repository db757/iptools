package iprange

import (
	"fmt"
	"strconv"
	"strings"
)

type PortRange interface {
	GetStart() uint16
	GetEnd() uint16
	fmt.Stringer
}

type TCPPortRange struct {
	Start uint16
	End   uint16
}

func (p *TCPPortRange) GetStart() uint16 {
	return p.Start
}

func (p *TCPPortRange) GetEnd() uint16 {
	return p.End
}

func (p *TCPPortRange) String() string {
	return fmt.Sprintf("%d-%d", p.Start, p.End)
}

func ParsePortRange(portStr string) (PortRange, error) {
	portRangeSlice := getPortRangeSlice(portStr)
	var portRange TCPPortRange
	var err error
	portRange.Start, err = StrToUint16(portRangeSlice[0])
	if err != nil {
		return nil, err
	}

	if len(portRangeSlice) < 2 {
		portRange.End = portRange.Start
		return &portRange, nil
	}

	portRange.End, err = StrToUint16(portRangeSlice[1])
	if err != nil {
		return nil, err
	}

	if portRange.Start > portRange.End {
		return nil, fmt.Errorf("start port cannot be greater than end port: %s", portStr)
	}

	return &portRange, nil
}

func getPortRangeSlice(portStr string) []string {
	for _, delim := range []string{":", "-"} {
		if strings.Contains(portStr, delim) {
			return strings.Split(portStr, delim)
		}
	}

	return []string{portStr}
}

func StrToUint16(str string) (uint16, error) {
	strAsInt, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(strAsInt), nil
}
