package handlers

import (
	"testing"
	"time"

	"github.com/bersennaidoo/arcbox/foundation/assert"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 1, 14, 10, 25, 0, 0, time.UTC),
			want: "14 Jan 2024 at 10:25",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "GET",
			tm:   time.Date(2024, 1, 14, 10, 25, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "14 Jan 2024 at 09:25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
