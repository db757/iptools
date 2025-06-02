package main

import (
	"fmt"
	"log"

	"github.com/db757/go-iprange/pkg/iprange"
	"go4.org/netipx"
)

type Result interface {
	Result() string
	Short() string
	Error() error
}

type ipInRangeResult struct {
	ip       string
	ranges   string
	err      error
	contains bool
}

func (r *ipInRangeResult) Result() string {
	if r.contains {
		return fmt.Sprintf("%s is in %s", r.ip, r.ranges)
	}
	return fmt.Sprintf("%s is NOT in %s", r.ip, r.ranges)
}

func (r *ipInRangeResult) Short() string {
	return fmt.Sprintf("%t", r.contains)
}

func (r *ipInRangeResult) Error() error {
	return r.err
}

func IPInRange(ipStr, ranges string) ipInRangeResult {
	result := ipInRangeResult{
		ip:     ipStr,
		ranges: ranges,
	}

	ipset, err := iprange.ParseRanges(ranges)
	if err != nil {
		log.Printf("failed to parse ranges %q: %v", ranges, err)
		result.err = err
		return result
	}

	ip, err := iprange.ParseIP(ipStr)
	if err != nil {
		log.Printf("failed to parse IP %q: %v", ipStr, err)
		result.err = err
		return result
	}

	result.contains = ipset.Contains(ip)
	return result
}

type cidrBoundariesResult struct {
	cidr    string
	ipRange netipx.IPRange
	err     error
}

func (r *cidrBoundariesResult) Result() string {
	return fmt.Sprintf("from: %s\nto: %s", r.ipRange.From(), r.ipRange.To())
}

func (r *cidrBoundariesResult) Short() string {
	return fmt.Sprintf("%s-%s", r.ipRange.From(), r.ipRange.To())
}

func (r *cidrBoundariesResult) Error() error {
	return r.err
}

func CIDRBoundaries(s string) cidrBoundariesResult {
	result := cidrBoundariesResult{
		cidr: s,
	}
	ipRange, err := iprange.CIDRToRange(s)
	if err != nil {
		log.Printf("failed to parse CIDR %q: %v", s, err)
		result.err = err
		return result
	}

	result.ipRange = ipRange
	return result
}
