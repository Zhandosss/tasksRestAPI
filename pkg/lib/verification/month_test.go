package verification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMonth(t *testing.T) {

	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{
			"one char correct month",
			"2",
			true,
		}, {
			"zero input",
			"0",
			false,
		}, {
			"two char correct month 1",
			"12",
			true,
		}, {
			"two char correct month 2",
			"10",
			true,
		}, {
			"two char incorrect month",
			"13",
			false,
		}, {
			"empty input",
			"",
			false,
		}, {
			"incorrect char",
			"1a",
			false,
		}, {
			"three number month",
			"123",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, monthVer(test.input))
		})
	}
}
