package iprange

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePortRange(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		portStr string
		want    PortRange
		wantErr bool
	}{
		{
			name:    "single valid port",
			portStr: "443",
			want:    &TCPPortRange{Start: 443, End: 443},
		},
		{
			name:    "valid port range",
			portStr: "443-5000",
			want:    &TCPPortRange{Start: 443, End: 5000},
		},
		{
			name:    "valid port range",
			portStr: "443:5000",
			want:    &TCPPortRange{Start: 443, End: 5000},
		},
		{
			name:    "empty string",
			portStr: "",
			wantErr: true,
		},
		{
			name:    "not numeric",
			portStr: "aa",
			wantErr: true,
		},
		{
			name:    "too large",
			portStr: "65536",
			wantErr: true,
		},
		{
			name:    "negative number",
			portStr: "-1",
			wantErr: true,
		},
		{
			name:    "negative number in range",
			portStr: "1--1",
			wantErr: true,
		},
		{
			name:    "negative number in range, different delim",
			portStr: "1:-1",
			wantErr: true,
		},
		{
			name:    "double delim",
			portStr: "1-:1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePortRange(tt.portStr)
			if tt.wantErr {
				assert.NotNil(err, "ParsePortRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.True(reflect.DeepEqual(got, tt.want), "ParsePortRange() = %v, want %v", got, tt.want)
		})
	}
}
