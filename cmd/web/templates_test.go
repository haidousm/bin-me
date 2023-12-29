package main

import (
	"testing"
	"time"

	"binme.haido.us/internal/assert"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 12, 29, 02, 47, 0, 0, time.UTC),
			want: "29 Dec 2023 at 02:47",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "PT",
			tm:   time.Date(2023, 12, 29, 02, 47, 0, 0, time.FixedZone("PDT", 8*60*60)),
			want: "28 Dec 2023 at 18:47",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
