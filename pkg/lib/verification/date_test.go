package verification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDate(t *testing.T) {

	var tests = []struct {
		name  string
		year  string
		month string
		day   string
		want  bool
	}{
		{
			"leap year",
			"2020",
			"2",
			"29",
			true,
		}, {
			"not leap year",
			"2021",
			"2",
			"29",
			false,
		}, {
			"30 day month 1",
			"2012",
			"4",
			"31",
			false,
		}, {
			"30 day month 2",
			"2012",
			"9",
			"31",
			false,
		}, {
			"february",
			"2013",
			"2",
			"30",
			false,
		}, {
			"normal",
			"2014",
			"10",
			"31",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Date(test.day, test.month, test.year))
		})
	}
}
