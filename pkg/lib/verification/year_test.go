package verification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYear(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "correct year",
			input: "2023",
			want:  true,
		}, {
			name:  "XX centry year",
			input: "1999",
			want:  false,
		}, {
			name:  "3 number year",
			input: "345",
			want:  false,
		}, {
			name:  "5 number year",
			input: "6789",
			want:  false,
		}, {
			name:  "empty year",
			input: "",
			want:  false,
		}, {
			name:  "incorect char",
			input: "2ab3",
			want:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, yearVer(test.input))
		})
	}
}
