package common

import (
	"testing"
)

func TestInterval_GetBelowInterval(t *testing.T) {
	tests := []struct {
		name string
		i    Interval
		want Interval
	}{
		{
			name: "already the lowest interval",
			i:    M1,
			want: M1,
		},
		{
			name: "straight case",
			i:    D1,
			want: H12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetBelowInterval(); got != tt.want {
				t.Errorf("Interval.GetBelowInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
