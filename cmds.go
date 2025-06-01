package main

import (
    "log"

    "github.com/db757/go-iprange/pkg/iprange"
)

func IPInRange(ipStr, ranges string) bool {
    ipset, err := iprange.ParseRanges(ranges)
    if err != nil {
        log.Printf("failed to parse ranges %q: %v", ranges, err)
        return false
    }

    ip, err := iprange.ParseIP(ipStr)
    if err != nil {
        log.Printf("failed to parse IP %q: %v", ipStr, err)
        return false
    }
    return ipset.Contains(ip)
}

func CIDRBoundaries(s string) (from, to string) {
    ipRange, err := iprange.CIDRToRange(s)
    if err != nil {
        log.Printf("failed to parse CIDR %q: %v", s, err)
        return "", ""
    }
    return ipRange.From().String(), ipRange.To().String()
}
