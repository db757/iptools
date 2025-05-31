package iprange

// Package for handling IP ranges, based on https://pkg.go.dev/github.com/netdata/go.d.plugin/pkg/iprange

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
)

// Family represents IP Range address-family.
type Family uint8

const (
	// V4Family is IPv4 address-family.
	V4Family Family = iota
	// V6Family is IPv6 address-family.
	V6Family
)

// Range represents an IP range.
type Range interface {
	Family() Family
	Contains(ip net.IP) bool
	GetStart() net.IP
	GetEnd() net.IP
	Size() *big.Int
	IsIP6() bool
	IsIP4() bool
	SameNetType(Range) bool
	fmt.Stringer
}

// New returns new IP Range.
// If it is not a valid range (start and end IPs have different address-families, or start > end),
// New returns nil.
func New(start, end net.IP) Range {
	if isV4RangeValid(start, end) {
		return IP4Range{start: start, end: end}
	}
	if isV6RangeValid(start, end) {
		return IP6Range{start: start, end: end}
	}
	return nil
}

type IP4Range struct {
	start net.IP
	end   net.IP
}

// String returns the string form of the range.
func (r IP4Range) String() string {
	return fmt.Sprintf("%s-%s", r.start, r.end)
}

// Family returns the range address family.
func (r IP4Range) Family() Family {
	return V4Family
}

// Contains reports whether the range includes IP.
func (r IP4Range) Contains(ip net.IP) bool {
	return bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) <= 0
}

func (r IP4Range) GetStart() net.IP {
	return r.start
}

func (r IP4Range) GetEnd() net.IP {
	return r.end
}

func (r IP4Range) IsIP6() bool {
	return false
}

func (r IP4Range) IsIP4() bool {
	return true
}

func (r IP4Range) SameNetType(otherRange Range) bool {
	return otherRange.IsIP4()
}

// Size reports the number of IP addresses in the range.
func (r IP4Range) Size() *big.Int {
	return big.NewInt(v4ToInt(r.end) - v4ToInt(r.start) + 1)
}

type IP6Range struct {
	start net.IP
	end   net.IP
}

// String returns the string form of the range.
func (r IP6Range) String() string {
	return fmt.Sprintf("%s-%s", r.start, r.end)
}

// Family returns the range address family.
func (r IP6Range) Family() Family {
	return V6Family
}

// Contains reports whether the range includes IP.
func (r IP6Range) Contains(ip net.IP) bool {
	return bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) <= 0
}

func (r IP6Range) GetStart() net.IP {
	return r.start
}

func (r IP6Range) GetEnd() net.IP {
	return r.end
}

func (r IP6Range) IsIP6() bool {
	return true
}

func (r IP6Range) IsIP4() bool {
	return false
}

func (r IP6Range) SameNetType(otherRange Range) bool {
	return otherRange.IsIP6()
}

// Size reports the number of IP addresses in the range.
func (r IP6Range) Size() *big.Int {
	size := big.NewInt(0)
	size.Add(size, big.NewInt(0).SetBytes(r.end))
	size.Sub(size, big.NewInt(0).SetBytes(r.start))
	size.Add(size, big.NewInt(1))
	return size
}
func v4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}
