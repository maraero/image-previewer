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
	validURL := "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"

	t.Run("valid params", func(t *testing.T) {
		ip, err := rs.ExtractParams("300/200/" + validURL)
		require.NoError(t, err)
		require.Equal(t, &ImageParams{Width: 300, Height: 200, URL: validURL}, ip)
	})

	t.Run("width is not number", func(t *testing.T) {
		_, err := rs.ExtractParams("a/100/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "width must be int")
	})

	t.Run("width is float number", func(t *testing.T) {
		_, err := rs.ExtractParams("34.34/100/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "width must be int")
	})

	t.Run("width is negative", func(t *testing.T) {
		_, err := rs.ExtractParams("-5/100/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "width must be greater than zero")
	})

	t.Run("height is not number", func(t *testing.T) {
		_, err := rs.ExtractParams("200/a/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "height must be int")
	})

	t.Run("height is float number", func(t *testing.T) {
		_, err := rs.ExtractParams("200/4.23/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "height must be int")
	})

	t.Run("height is negative", func(t *testing.T) {
		_, err := rs.ExtractParams("200/-100/" + validURL)
		require.Error(t, err)
		require.ErrorContains(t, err, "height must be greater than zero")
	})
}
