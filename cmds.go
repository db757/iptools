package main

import (
	"github.com/db757/go-iprange/pkg/iprange"
)

func IPInRange(ipStr, ranges string) bool {
	ipset, err := iprange.ParseRanges(ranges)
	if err != nil {
		return false
	}

	ip := iprange.ParseIP(ipStr)
	return ipset.Contains(ip)
}

func CIDRBoundaries(s string) (from, to string) {
	ipRange := iprange.CIDRToRange(s)
	return ipRange.From().String(), ipRange.To().String()
}
