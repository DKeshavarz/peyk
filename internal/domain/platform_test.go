package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	testCase := []struct{
		tag string
		platform PlatformName
		want bool
	}{
		{
			tag: "telgram",
			platform: Telegram,
			want: true,
		},
		{
			tag: "bale",
			platform: Bale,
			want: true,
		},
		{
			tag: "bad platform",
			platform: "moew",
			want: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.tag, func(t *testing.T) {
			got := tc.platform.Valid()

			assert.Equal(t, tc.want, got)
		})
	}
}
