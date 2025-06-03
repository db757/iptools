package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/netip"
	"strings"

	"github.com/db757/iptools/pkg/parse"
	"github.com/urfave/cli/v3"
	"go4.org/netipx"
)

type AppInput struct {
	Primary   string
	Secondary string
	Count     int
}

type AppConfig struct {
	Short bool
}

type App struct {
	Config AppConfig
	Input  AppInput
}

func NewApp() App {
	return App{}
}

func (a *App) InRangeHandler(context.Context, *cli.Command) error {
	result := ipInRange(a.Input.Primary, a.Input.Secondary)
	return a.handleResult(&result)
}

func (a *App) CIDRBoundariesHandler(context.Context, *cli.Command) error {
	result := cidrBoundaries(a.Input.Primary)
	return a.handleResult(&result)
}

func (a *App) NextHandler(context.Context, *cli.Command) error {
	result := next(a.Input.Primary)
	return a.handleResult(&result)
}

func (a *App) PrevHandler(context.Context, *cli.Command) error {
	result := prev(a.Input.Primary)
	return a.handleResult(&result)
}

func (a *App) GetNHandler(_ context.Context, cmd *cli.Command) error {
	result := getN(a.Input.Primary, a.Input.Count, cmd.Int("offset"), cmd.Bool("tail"))
	return a.handleResult(&result)
}

func (a *App) handleResult(result Result) error {
	if result.Error() != nil {
		return result.Error()
	}

	if a.Config.Short {
		fmt.Println(result.Short())
		return nil
	}

	fmt.Println(result.Result())
	return nil
}

type Result interface {
	Result() string
	Short() string
	Error() error
}

type ResultError struct {
	err error
}

func (r *ResultError) Error() error {
	return r.err
}

type ipInRangeResult struct {
	ip       string
	ranges   string
	contains bool
	ResultError
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

func ipInRange(ipStr, ranges string) ipInRangeResult {
	result := ipInRangeResult{
		ip:     ipStr,
		ranges: ranges,
	}

	ipset, err := parse.ParseRanges(ranges)
	if err != nil {
		log.Printf("failed to parse ranges %q: %v", ranges, err)
		result.err = err
		return result
	}

	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		log.Printf("failed to parse IP %q: %v", ipStr, err)
		result.err = err
		return result
	}

	result.contains = ipset.Contains(ip)
	return result
}

type cidrBoundariesResult struct {
	from, to string
	cidr     netip.Prefix
	cidrLen  string
	ResultError
}

func (r *cidrBoundariesResult) Result() string {
	return fmt.Sprintf("%s (%s addresses):\nfrom: %s\nto: %s", r.cidr, cidrLen(r.cidr), r.from, r.to)
}

func cidrLen(cidr netip.Prefix) string {
	bitlen := cidr.Addr().BitLen()
	if bitlen == 0 {
		return "unknown number of"
	}

	if cidr.Addr().Is6() && cidr.Bits() < 65 {
		// Too large to calculate directly
		return fmt.Sprintf("more than %d", uint64(math.MaxUint64))
	}

	cidrLen := uint64(1) << (bitlen - cidr.Bits())
	return fmt.Sprintf("%d", cidrLen)
}

func (r *cidrBoundariesResult) Short() string {
	return fmt.Sprintf("%s-%s", r.from, r.to)
}

func (r *cidrBoundariesResult) Error() error {
	return r.err
}

func cidrBoundaries(s string) cidrBoundariesResult {
	result := cidrBoundariesResult{}

	cidr, err := netip.ParsePrefix(s)
	if err != nil {
		result.err = fmt.Errorf("failed to parse CIDR %q: %w", s, err)
		return result
	}

	result.cidr = cidr
	result.cidrLen = cidrLen(cidr)
	result.from = cidr.Addr().String()
	result.to = netipx.PrefixLastIP(cidr).String()
	return result
}

type NextPrevResult struct {
	outputPrefix string
	addr         netip.Addr
	ResultError
}

func (r *NextPrevResult) Result() string {
	return fmt.Sprintf("%s IP: %s", r.outputPrefix, r.addr)
}

func (r *NextPrevResult) Short() string {
	return r.addr.String()
}

func next(s string) NextPrevResult {
	result := NextPrevResult{
		outputPrefix: "Next",
	}
	ip, err := netip.ParseAddr(s)
	if err != nil {
		result.err = fmt.Errorf("failed to parse IP %q: %w", s, err)
	}

	result.addr = ip.Next()
	return result
}

func prev(s string) NextPrevResult {
	result := NextPrevResult{
		outputPrefix: "Prev",
	}
	ip, err := netip.ParseAddr(s)
	if err != nil {
		result.err = fmt.Errorf("failed to parse IP %q: %w", s, err)
	}

	result.addr = ip.Prev()
	return result
}

type GetNResult struct {
	ips   []string
	count int
	ResultError
}

func (r *GetNResult) Result() string {
	return fmt.Sprintf("%d IPs: %s", r.count, strings.Join(r.ips, ","))
}

func (r *GetNResult) Short() string {
	return fmt.Sprint(strings.Join(r.ips, ","))
}

func getN(s string, count int, offset int, tail bool) GetNResult {
	result := GetNResult{
		count: count,
	}

	cidr, err := netip.ParsePrefix(s)
	if err != nil {
		result.err = fmt.Errorf("failed to parse prefix %s: %w", s, err)
		return result
	}

	ipRange := netipx.RangeOfPrefix(cidr)

	if count < 1 {
		result.err = fmt.Errorf("count must be greater than 0")
		return result
	}

	// Skip network address
	ip := ipRange.From().Next()
	last := ipRange.To()

	next := func(ip netip.Addr) netip.Addr {
		return ip.Next()
	}

	if tail {
		ip = ipRange.To().Prev()
		last = ipRange.From()
		next = func(ip netip.Addr) netip.Addr {
			return ip.Prev()
		}
	}

	for range offset {
		ip = next(ip)
	}

	for range count {
		if !ip.IsValid() || ip == last {
			break
		}
		result.ips = append(result.ips, ip.String())
		ip = next(ip)
	}

	return result
}
