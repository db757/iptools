package iprange

import (
	"net/netip"
	"regexp"
	"strings"

	"go4.org/netipx"
)

// ParseRanges parses s as a space separated list of IP Ranges, returning the result and an error if any.
// IP Range can be in IPv4 address ("192.0.2.1"), IPv4 range ("192.0.2.0-192.0.2.10")
// IPv4 CIDR ("192.0.2.0/24")
// IPv6 address ("2001:db8::1"), IPv6 range ("2001:db8::-2001:db8::10"),
// or IPv6 CIDR ("2001:db8::/64") form.
// IPv4 CIDR, IPv4 subnet mask and IPv6 CIDR ranges don't include network and broadcast addresses.
func ParseRanges(s string) (*netipx.IPSet, error) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil, nil
	}

	var builder netipx.IPSetBuilder
	for _, v := range parts {
		ipset, err := ParseRange(v)
		if err != nil {
			return nil, err
		}

		if ipset != nil {
			builder.AddSet(ipset)
		}
	}
	return builder.IPSet()
}

var (
	reRange = regexp.MustCompile("^[0-9a-f.:-]+$")           // addr | addr-addr
	reIP    = regexp.MustCompile("^[0-9a-f.:]+$")            // addr
	reCIDR  = regexp.MustCompile("^[0-9a-f.:]+/[0-9]{1,3}$") // addr/prefix_length
)

// ParseRange parses s as an IP Range, returning the result and an error if any.
// The string s can be in IPv4 address ("192.0.2.1"), IPv4 range ("192.0.2.0-192.0.2.10")
// IPv4 CIDR ("192.0.2.0/24")
// IPv6 address ("2001:db8::1"), IPv6 range ("2001:db8::-2001:db8::10"),
// or IPv6 CIDR ("2001:db8::/64") form.
// IPv4 CIDR, IPv4 subnet mask and IPv6 CIDR ranges don't include network and broadcast addresses.
func ParseRange(s string) (*netipx.IPSet, error) {
	s = strings.ToLower(s)
	if s == "" {
		return nil, nil
	}

	var builder netipx.IPSetBuilder
	switch {
	case reIP.MatchString(s):
		builder.Add(netip.MustParseAddr(s))
	case reRange.MatchString(s):
		builder.AddRange(netipx.MustParseIPRange(s))
	case reCIDR.MatchString(s):
		builder.AddPrefix(netip.MustParsePrefix(s))
	}

	return builder.IPSet()
}

func ParseIP(s string) netip.Addr {
	return netip.MustParseAddr(s)
}
