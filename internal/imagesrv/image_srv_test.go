package imagesrv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCacheKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "600/300/https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png",
			expected: "600_300_https___upload_wikimedia_org_wikipedia_commons_4_47_PNG_transparency_demonstration_1_png",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			res := getCacheKey(tc.input)
			fmt.Println(res)
			require.Equal(t, tc.expected, res)
		})
	}
}
