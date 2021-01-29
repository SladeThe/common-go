package strings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sortArgs struct {
	a string
	b string
}

type sortWant struct {
	less bool
}

var sortTests = []struct {
	name string
	args sortArgs
	want sortWant
}{{
	name: "empty",
	args: sortArgs{},
	want: sortWant{},
}, {
	name: "empty - non empty",
	args: sortArgs{b: "b"},
	want: sortWant{less: true},
}, {
	name: "non empty - empty",
	args: sortArgs{a: "a"},
	want: sortWant{less: false},
}, {
	name: "alfa less",
	args: sortArgs{a: "a", b: "b"},
	want: sortWant{less: true},
}, {
	name: "alfa equal",
	args: sortArgs{a: "v", b: "v"},
	want: sortWant{less: false},
}, {
	name: "alfa more",
	args: sortArgs{a: "b", b: "a"},
	want: sortWant{less: false},
}, {
	name: "digit - non digit",
	args: sortArgs{a: "1", b: "b"},
	want: sortWant{less: true},
}, {
	name: "non digit - digit",
	args: sortArgs{a: "a", b: "1"},
	want: sortWant{less: false},
}, {
	name: "numbers 1",
	args: sortArgs{a: "v123", b: "v9"},
	want: sortWant{less: false},
}, {
	name: "numbers 2",
	args: sortArgs{a: "8", b: "11"},
	want: sortWant{less: true},
}, {
	name: "numbers 3",
	args: sortArgs{a: "11", b: "11"},
	want: sortWant{less: false},
}, {
	name: "zero numbers 1",
	args: sortArgs{a: "01", b: "01"},
	want: sortWant{less: false},
}, {
	name: "zero numbers 2",
	args: sortArgs{a: "01", b: "1"},
	want: sortWant{less: false},
}, {
	name: "zero numbers 3",
	args: sortArgs{a: "1", b: "01"},
	want: sortWant{less: true},
}, {
	name: "zero numbers 1a",
	args: sortArgs{a: "v01v", b: "v01v"},
	want: sortWant{less: false},
}, {
	name: "zero numbers 2a",
	args: sortArgs{a: "v01v", b: "v1v"},
	want: sortWant{less: false},
}, {
	name: "zero numbers 3a",
	args: sortArgs{a: "v1v", b: "v01v"},
	want: sortWant{less: true},
}, {
	name: "zero - non zero",
	args: sortArgs{a: "v01v", b: "v2v"},
	want: sortWant{less: true},
}, {
	name: "non zero - zero",
	args: sortArgs{a: "v1v", b: "v02v"},
	want: sortWant{less: true},
}, {
	name: "zero a",
	args: sortArgs{a: "0", b: ""},
	want: sortWant{less: false},
}, {
	name: "zero b",
	args: sortArgs{a: "", b: "0"},
	want: sortWant{less: true},
}, {
	name: "zeroes a",
	args: sortArgs{a: "00", b: "0"},
	want: sortWant{less: false},
}, {
	name: "zeroes b",
	args: sortArgs{a: "0", b: "00"},
	want: sortWant{less: true},
}, {
	name: "zeroes",
	args: sortArgs{a: "000", b: "000"},
	want: sortWant{less: false},
}, {
	name: "zeroes and letter a",
	args: sortArgs{a: "000a", b: "000"},
	want: sortWant{less: false},
}, {
	name: "zeroes and letter b",
	args: sortArgs{a: "000", b: "000b"},
	want: sortWant{less: true},
}, {
	name: "zeroes and letters 1",
	args: sortArgs{a: "000a", b: "000b"},
	want: sortWant{less: true},
}, {
	name: "zeroes and letters 2",
	args: sortArgs{a: "00v", b: "000v"},
	want: sortWant{less: true},
}, {
	name: "zeroes and letters 3",
	args: sortArgs{a: "000v", b: "00v"},
	want: sortWant{less: false},
}, {
	name: "zeroes digit - zeroes letter",
	args: sortArgs{a: "0001", b: "000b"},
	want: sortWant{less: false},
}, {
	name: "zeroes letter - zeroes digit",
	args: sortArgs{a: "000a", b: "0001"},
	want: sortWant{less: true},
}, {
	name: "numbers and letter 1",
	args: sortArgs{a: "001v", b: "002v"},
	want: sortWant{less: true},
}, {
	name: "numbers and letter 2",
	args: sortArgs{a: "001v", b: "001v"},
	want: sortWant{less: false},
}, {
	name: "numbers and letter 3",
	args: sortArgs{a: "002v", b: "001v"},
	want: sortWant{less: false},
}, {
	name: "numbers and letter 4",
	args: sortArgs{a: "001a", b: "001b"},
	want: sortWant{less: true},
}, {
	name: "numbers and letter 5",
	args: sortArgs{a: "001b", b: "001a"},
	want: sortWant{less: false},
}, {
	name: "numbers and letter 6",
	args: sortArgs{a: "001", b: "001b"},
	want: sortWant{less: true},
}, {
	name: "numbers and letter 7",
	args: sortArgs{a: "001a", b: "001"},
	want: sortWant{less: false},
}, {
	name: "long numbers 1",
	args: sortArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
	},
	want: sortWant{less: false},
}, {
	name: "long numbers 2",
	args: sortArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "888",
	},
	want: sortWant{less: false},
}, {
	name: "long numbers 3",
	args: sortArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "888",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
	},
	want: sortWant{less: true},
}, {
	name: "long numbers 1a",
	args: sortArgs{
		a: strings.Repeat("1", 1000) + "88",
		b: strings.Repeat("1", 1000) + "999",
	},
	want: sortWant{less: true},
}, {
	name: "long numbers 2a",
	args: sortArgs{
		a: strings.Repeat("1", 1000) + "99",
		b: strings.Repeat("1", 1000) + "888",
	},
	want: sortWant{less: true},
}, {
	name: "long numbers 3a",
	args: sortArgs{
		a: strings.Repeat("1", 1000) + "888",
		b: strings.Repeat("1", 1000) + "99",
	},
	want: sortWant{less: false},
}}

func BenchmarkNumericLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range sortTests {
			if NumericLess(tt.args.a, tt.args.b) != tt.want.less {
				require.Equal(b, tt.want.less, NumericLess(tt.args.a, tt.args.b), "invalid less")
			}
		}
	}
}

func TestNumericLess(t *testing.T) {
	test := func(a sortArgs, w sortWant) func(t *testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, w.less, NumericLess(a.a, a.b), "invalid less")
		}
	}

	for _, tt := range sortTests {
		t.Run(tt.name, test(tt.args, tt.want))
	}
}
