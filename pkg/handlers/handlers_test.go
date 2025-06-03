package handlers

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func captureOutput(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = w
	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	err = fn()
	w.Close()
	os.Stdout = old
	out := <-outC

	return out, err
}

func TestInRangeHandler(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		ranges    string
		short     bool
		wantShort string
		wantLong  string
		wantErr   bool
	}{
		{
			name:      "IPv4 in range",
			ip:        "192.0.2.5",
			ranges:    "192.0.2.0-192.0.2.10",
			wantShort: "true",
			wantLong:  "192.0.2.5 is in 192.0.2.0-192.0.2.10",
		},
		{
			name:      "IPv6 in range",
			ip:        "2001:db8::5",
			ranges:    "2001:db8::1-2001:db8::10",
			wantShort: "true",
			wantLong:  "2001:db8::5 is in 2001:db8::1-2001:db8::10",
		},
		{
			name:      "IP not in range",
			ip:        "192.0.2.20",
			ranges:    "192.0.2.0-192.0.2.10",
			wantShort: "false",
			wantLong:  "192.0.2.20 is NOT in 192.0.2.0-192.0.2.10",
		},
		{
			name:      "IP in CIDR",
			ip:        "192.0.2.128",
			ranges:    "192.0.2.0/24",
			wantShort: "true",
			wantLong:  "192.0.2.128 is in 192.0.2.0/24",
		},
		{
			name:      "IP not in CIDR",
			ip:        "192.0.3.1",
			ranges:    "192.0.2.0/24",
			wantShort: "false",
			wantLong:  "192.0.3.1 is NOT in 192.0.2.0/24",
		},
		{
			name:      "IP in multiple ranges",
			ip:        "192.0.2.5",
			ranges:    "192.0.2.0-192.0.2.10,192.0.2.20-192.0.2.30",
			wantShort: "true",
			wantLong:  "192.0.2.5 is in 192.0.2.0-192.0.2.10,192.0.2.20-192.0.2.30",
		},
		{
			name:      "IP in multiple CIDRs",
			ip:        "192.0.2.1",
			ranges:    "192.0.2.0/24,192.0.3.0/24",
			wantShort: "true",
			wantLong:  "192.0.2.1 is in 192.0.2.0/24,192.0.3.0/24",
		},
		{
			name:      "IP in mixed ranges and CIDRs",
			ip:        "192.0.2.5",
			ranges:    "192.0.2.0-192.0.2.10,192.0.3.0/24",
			wantShort: "true",
			wantLong:  "192.0.2.5 is in 192.0.2.0-192.0.2.10,192.0.3.0/24",
		},
		{
			name:      "Specific IP match",
			ip:        "192.0.2.1",
			ranges:    "192.0.2.1",
			wantShort: "true",
			wantLong:  "192.0.2.1 is in 192.0.2.1",
		},
		{
			name:    "Invalid IP",
			ip:      "256.256.256.256",
			ranges:  "192.0.2.0-192.0.2.10",
			wantErr: true,
		},
		{
			name:    "Invalid range",
			ip:      "192.0.2.5",
			ranges:  "invalid-range",
			wantErr: true,
		},
		{
			name:    "Invalid CIDR",
			ip:      "192.0.2.5",
			ranges:  "192.0.2.0/33",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.Input.Primary = tt.ip
			app.Input.Secondary = tt.ranges
			app.Config.Short = tt.short

			out, err := captureOutput(t, func() error {
				return app.InRangeHandler(context.Background(), nil)
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.short {
				assert.Equal(t, tt.wantShort+"\n", out)
			} else {
				assert.Equal(t, tt.wantLong+"\n", out)
			}
		})
	}
}

func TestCIDRBoundariesHandler(t *testing.T) {
	tests := []struct {
		name      string
		cidr      string
		short     bool
		wantShort string
		wantLong  string
		wantErr   bool
	}{
		{
			name:      "Valid IPv4 CIDR",
			cidr:      "192.0.2.0/24",
			wantShort: "192.0.2.0-192.0.2.255",
			wantLong:  "192.0.2.0/24 (256 addresses):\nfrom: 192.0.2.0\nto: 192.0.2.255",
		},
		{
			name:      "Valid IPv6 CIDR",
			cidr:      "2001:db8::/64",
			wantShort: "2001:db8::-2001:db8::ffff:ffff:ffff:ffff",
			wantLong:  "2001:db8::/64 (more than 18446744073709551615 addresses):\nfrom: 2001:db8::\nto: 2001:db8::ffff:ffff:ffff:ffff",
		},
		{
			name:    "Invalid CIDR",
			cidr:    "invalid-cidr",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.Input.Primary = tt.cidr
			app.Config.Short = tt.short

			out, err := captureOutput(t, func() error {
				return app.CIDRBoundariesHandler(context.Background(), nil)
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.short {
				assert.Equal(t, tt.wantShort+"\n", out)
			} else {
				assert.Equal(t, tt.wantLong+"\n", out)
			}
		})
	}
}

func TestNextHandler(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		short     bool
		wantShort string
		wantLong  string
		wantErr   bool
	}{
		{
			name:      "Valid IPv4",
			ip:        "192.0.2.1",
			wantShort: "192.0.2.2",
			wantLong:  "Next IP: 192.0.2.2",
		},
		{
			name:      "Valid IPv6",
			ip:        "2001:db8::1",
			wantShort: "2001:db8::2",
			wantLong:  "Next IP: 2001:db8::2",
		},
		{
			name:    "Invalid IP",
			ip:      "invalid-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.Input.Primary = tt.ip
			app.Config.Short = tt.short

			out, err := captureOutput(t, func() error {
				return app.NextHandler(context.Background(), nil)
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.short {
				assert.Equal(t, tt.wantShort+"\n", out)
			} else {
				assert.Equal(t, tt.wantLong+"\n", out)
			}
		})
	}
}

func TestPrevHandler(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		short     bool
		wantShort string
		wantLong  string
		wantErr   bool
	}{
		{
			name:      "Valid IPv4",
			ip:        "192.0.2.2",
			wantShort: "192.0.2.1",
			wantLong:  "Prev IP: 192.0.2.1",
		},
		{
			name:      "Valid IPv6",
			ip:        "2001:db8::2",
			wantShort: "2001:db8::1",
			wantLong:  "Prev IP: 2001:db8::1",
		},
		{
			name:    "Invalid IP",
			ip:      "invalid-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.Input.Primary = tt.ip
			app.Config.Short = tt.short

			out, err := captureOutput(t, func() error {
				return app.PrevHandler(context.Background(), nil)
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.short {
				assert.Equal(t, tt.wantShort+"\n", out)
			} else {
				assert.Equal(t, tt.wantLong+"\n", out)
			}
		})
	}
}

func TestGetNHandler(t *testing.T) {
	tests := []struct {
		name      string
		cidr      string
		count     int
		offset    int
		tail      bool
		short     bool
		wantShort string
		wantLong  string
		wantErr   bool
	}{
		{
			name:      "Forward traversal",
			cidr:      "192.0.2.0/24",
			count:     3,
			offset:    0,
			tail:      false,
			wantShort: "192.0.2.1,192.0.2.2,192.0.2.3",
			wantLong:  "3 IPs: 192.0.2.1,192.0.2.2,192.0.2.3",
		},
		{
			name:      "Backward traversal",
			cidr:      "192.0.2.0/24",
			count:     3,
			offset:    0,
			tail:      true,
			wantShort: "192.0.2.254,192.0.2.253,192.0.2.252",
			wantLong:  "3 IPs: 192.0.2.254,192.0.2.253,192.0.2.252",
		},
		{
			name:      "With offset",
			cidr:      "192.0.2.0/24",
			count:     2,
			offset:    2,
			tail:      false,
			wantShort: "192.0.2.3,192.0.2.4",
			wantLong:  "2 IPs: 192.0.2.3,192.0.2.4",
		},
		{
			name:    "Invalid count",
			cidr:    "192.0.2.0/24",
			count:   0,
			wantErr: true,
		},
		{
			name:    "Invalid CIDR",
			cidr:    "invalid-cidr",
			count:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.Input.Primary = tt.cidr
			app.Input.Count = tt.count
			app.Config.Short = tt.short

			cmd := &cli.Command{
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "offset", Value: tt.offset},
					&cli.BoolFlag{Name: "tail", Value: tt.tail},
				},
			}

			out, err := captureOutput(t, func() error {
				return app.GetNHandler(context.Background(), cmd)
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.short {
				assert.Equal(t, tt.wantShort+"\n", out)
			} else {
				assert.Equal(t, tt.wantLong+"\n", out)
			}
		})
	}
}
