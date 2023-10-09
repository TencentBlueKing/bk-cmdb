package params

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpecialCharChange(t *testing.T) {
	type testUnit struct {
		src string
		dst string
	}

	testUnits := []testUnit{
		{
			src: ".",
			dst: `\.`,
		},
		{
			src: "(",
			dst: `\(`,
		},
		{
			src: ")",
			dst: `\)`,
		},
		{
			src: "\\",
			dst: `\\`,
		},
		{
			src: "|",
			dst: `\|`,
		},
		{
			src: "[",
			dst: `\[`,
		},
		{
			src: "]",
			dst: `\]`,
		},
		{
			src: "-",
			dst: `\-`,
		},
		{
			src: "*",
			dst: `\*`,
		},
		{
			src: "{",
			dst: `\{`,
		},
		{
			src: "}",
			dst: `\}`,
		},
		{
			src: "^",
			dst: `\^`,
		},
		{
			src: "$",
			dst: `\$`,
		},
		{
			src: "?",
			dst: `\?`,
		},

		{
			src: "aaa",
			dst: `aaa`,
		},
		{
			src: "12345",
			dst: `12345`,
		},
		{
			src: "12345",
			dst: `12345`,
		},
		{
			src: "!@#%&_+=,<>/`~",
			dst: "!@#%&_+=,<>/`~",
		},
		{
			src: "12345676890qwertyuiopasdfghjklzxcvbnm QWERTYUIOPASDFGHJKLZXCVBNM",
			dst: "12345676890qwertyuiopasdfghjklzxcvbnm QWERTYUIOPASDFGHJKLZXCVBNM",
		},
	}

	for _, item := range testUnits {
		dst := SpecialCharChange(item.src)
		// blog.Infof("src:%s, expact dst:%s, dst:%s", item.src, item.dst, dst)
		require.Equal(t, item.dst, dst)
	}

}
