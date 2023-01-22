package resizesrv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractParams(t *testing.T) {
	rs := New()

	t.Run("too few params", func(t *testing.T) {
		tests := []string{"", "200", "200/300"}
		for _, tc := range tests {
			tc := tc
			t.Run(tc, func(t *testing.T) {
				_, err := rs.ExtractParams("")
				require.Error(t, err)
				require.ErrorIs(t, err, ErrTooFewParams)
			})
		}
	})

	//nolint:lll
	testURL := "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"

	t.Run("valid params", func(t *testing.T) {
		ip, err := rs.ExtractParams("300/200/" + testURL)
		require.NoError(t, err)
		require.Equal(t, &ImageParams{Width: 300, Height: 200, URL: testURL}, ip)
	})

	t.Run("width validation", func(t *testing.T) {
		tests := []struct {
			input         string
			expectedError string
		}{
			{input: "a/100/" + testURL, expectedError: "width must be int"},
			{input: "34.34/100/" + testURL, expectedError: "width must be int"},
			{input: "-5/100/" + testURL, expectedError: "width must be greater than zero"},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.input, func(t *testing.T) {
				_, err := rs.ExtractParams(tc.input)
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedError)
			})
		}
	})

	t.Run("height validation", func(t *testing.T) {
		tests := []struct {
			input         string
			expectedError string
		}{
			{input: "200/a/" + testURL, expectedError: "height must be int"},
			{input: "200/4.23/" + testURL, expectedError: "height must be int"},
			{input: "200/-100/" + testURL, expectedError: "height must be greater than zero"},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.input, func(t *testing.T) {
				_, err := rs.ExtractParams(tc.input)
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedError)
			})
		}
	})

	t.Run("adds schema to URLs if it is required", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{input: "test.com", expected: "http://test.com"},
			{input: "http://test.com", expected: "http://test.com"},
			{input: "https://test.com", expected: "https://test.com"},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.input, func(t *testing.T) {
				result, err := rs.ExtractParams("200/100/" + tc.input)
				require.NoError(t, err)
				require.Equal(t, tc.expected, result.URL)
			})
		}
	})
}
