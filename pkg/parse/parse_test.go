package parse

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"go4.org/netipx"
)

func TestParseRanges(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string // Expected IPs/Ranges in the set
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single IPv4",
			input: "192.0.2.1",
			want:  []string{"192.0.2.1"},
		},
		{
			name:  "single IPv6",
			input: "2001:db8::1",
			want:  []string{"2001:db8::1"},
		},
		{
			name:  "IPv4 range",
			input: "192.0.2.0-192.0.2.10",
			want:  []string{"192.0.2.0-192.0.2.10"},
		},
		{
			name:  "IPv6 range",
			input: "2001:db8::-2001:db8::10",
			want:  []string{"2001:db8::-2001:db8::10"},
		},
		{
			name:  "IPv4 CIDR",
			input: "192.0.2.0/24",
			want:  []string{"192.0.2.0/24"},
		},
		{
			name:  "IPv6 CIDR",
			input: "2001:db8::/64",
			want:  []string{"2001:db8::/64"},
		},
		{
			name:  "multiple mixed formats",
			input: "192.0.2.1 2001:db8::1 192.0.2.0-192.0.2.10 2001:db8::/64",
			want:  []string{"192.0.2.1", "2001:db8::1", "192.0.2.0-192.0.2.10", "2001:db8::/64"},
		},
		{
			name:  "multiple mixed formats, comma separated",
			input: "192.0.2.1,2001:db8::1,192.0.2.0-192.0.2.10,2001:db8::/64",
			want:  []string{"192.0.2.1", "2001:db8::1", "192.0.2.0-192.0.2.10", "2001:db8::/64"},
		},
		{
			name:    "invalid IP",
			input:   "256.256.256.256",
			wantErr: true,
		},
		{
			name:    "invalid range",
			input:   "192.0.2.10-192.0.2.0", // end before start
			wantErr: true,
		},
		{
			name:    "invalid CIDR",
			input:   "192.0.2.0/33", // invalid prefix length
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ranges(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			for _, wantRange := range tt.want {
				// Try parsing as single IP first
				if addr, err := netip.ParseAddr(wantRange); err == nil {
					assert.True(t, got.Contains(addr), "expected set to contain %s", wantRange)
					continue
				}

				// Try parsing as IP range
				if ipRange, err := netipx.ParseIPRange(wantRange); err == nil {
					assert.True(t, got.ContainsRange(ipRange), "expected set to contain range %s", wantRange)
					continue
				}

				// Try parsing as CIDR
				prefix, err := netip.ParsePrefix(wantRange)
				assert.NoError(t, err, "failed to parse any format for %s", wantRange)
				rng := netipx.RangeOfPrefix(prefix)
				assert.True(t, got.ContainsRange(rng), "expected set to contain CIDR range %s", wantRange)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single IPv4",
			input: "192.0.2.1",
			want:  "192.0.2.1",
		},
		{
			name:  "single IPv6",
			input: "2001:db8::1",
			want:  "2001:db8::1",
		},
		{
			name:  "IPv4 range",
			input: "192.0.2.0-192.0.2.10",
			want:  "192.0.2.0-192.0.2.10",
		},
		{
			name:  "IPv6 range",
			input: "2001:db8::-2001:db8::10",
			want:  "2001:db8::-2001:db8::10",
		},
		{
			name:  "IPv4 CIDR",
			input: "192.0.2.0/24",
			want:  "192.0.2.0/24",
		},
		{
			name:  "IPv6 CIDR",
			input: "2001:db8::/64",
			want:  "2001:db8::/64",
		},
		{
			name:    "invalid IP",
			input:   "256.256.256.256",
			wantErr: true,
		},
		{
			name:    "invalid range",
			input:   "192.0.2.10-192.0.2.0", // end before start
			wantErr: true,
		},
		{
			name:    "invalid CIDR",
			input:   "192.0.2.0/33", // invalid prefix length
			wantErr: true,
		},
		{
			name:    "malformed input",
			input:   "not-an-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Range(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.want == "" {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			// Try parsing as single IP first
			if addr, err := netip.ParseAddr(tt.want); err == nil {
				assert.True(t, got.Contains(addr))
				return
			}

			// Try parsing as IP range
			if ipRange, err := netipx.ParseIPRange(tt.want); err == nil {
				assert.True(t, got.ContainsRange(ipRange))
				return
			}

			// Try parsing as CIDR
			prefix, err := netip.ParsePrefix(tt.want)
			assert.NoError(t, err, "failed to parse %s in any format", tt.want)
			rng := netipx.RangeOfPrefix(prefix)
			assert.True(t, got.ContainsRange(rng))
		})
	}
}
