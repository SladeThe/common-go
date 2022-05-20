package strings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type numericLessArgs struct {
	a string
	b string
}

type numericLessWant struct {
	less bool
}

var numericLessTests = []struct {
	name string
	args numericLessArgs
	want numericLessWant
}{{
	name: "empty",
	args: numericLessArgs{},
	want: numericLessWant{},
}, {
	name: "empty - non empty",
	args: numericLessArgs{b: "b"},
	want: numericLessWant{less: true},
}, {
	name: "non empty - empty",
	args: numericLessArgs{a: "a"},
	want: numericLessWant{less: false},
}, {
	name: "alfa less",
	args: numericLessArgs{a: "a", b: "b"},
	want: numericLessWant{less: true},
}, {
	name: "alfa equal",
	args: numericLessArgs{a: "v", b: "v"},
	want: numericLessWant{less: false},
}, {
	name: "alfa more",
	args: numericLessArgs{a: "b", b: "a"},
	want: numericLessWant{less: false},
}, {
	name: "digit - non digit 1",
	args: numericLessArgs{a: "1", b: "b"},
	want: numericLessWant{less: true},
}, {
	name: "digit - non digit 2",
	args: numericLessArgs{
		a: "T2 x",
		b: "T  x",
	},
	want: numericLessWant{less: false},
}, {
	name: "non digit - digit 1",
	args: numericLessArgs{a: "a", b: "1"},
	want: numericLessWant{less: false},
}, {
	name: "non digit - digit 2",
	args: numericLessArgs{
		a: "T  x",
		b: "T2 x",
	},
	want: numericLessWant{less: true},
}, {
	name: "numbers 1",
	args: numericLessArgs{a: "v123", b: "v9"},
	want: numericLessWant{less: false},
}, {
	name: "numbers 2",
	args: numericLessArgs{a: "8", b: "11"},
	want: numericLessWant{less: true},
}, {
	name: "numbers 3",
	args: numericLessArgs{a: "11", b: "11"},
	want: numericLessWant{less: false},
}, {
	name: "zero numbers 1",
	args: numericLessArgs{a: "01", b: "01"},
	want: numericLessWant{less: false},
}, {
	name: "zero numbers 2",
	args: numericLessArgs{a: "01", b: "1"},
	want: numericLessWant{less: false},
}, {
	name: "zero numbers 3",
	args: numericLessArgs{a: "1", b: "01"},
	want: numericLessWant{less: true},
}, {
	name: "zero numbers 1a",
	args: numericLessArgs{a: "v01v", b: "v01v"},
	want: numericLessWant{less: false},
}, {
	name: "zero numbers 2a",
	args: numericLessArgs{a: "v01v", b: "v1v"},
	want: numericLessWant{less: false},
}, {
	name: "zero numbers 3a",
	args: numericLessArgs{a: "v1v", b: "v01v"},
	want: numericLessWant{less: true},
}, {
	name: "zero - non zero",
	args: numericLessArgs{a: "v01v", b: "v2v"},
	want: numericLessWant{less: true},
}, {
	name: "non zero - zero",
	args: numericLessArgs{a: "v1v", b: "v02v"},
	want: numericLessWant{less: true},
}, {
	name: "zero a",
	args: numericLessArgs{a: "0", b: ""},
	want: numericLessWant{less: false},
}, {
	name: "zero b",
	args: numericLessArgs{a: "", b: "0"},
	want: numericLessWant{less: true},
}, {
	name: "zeroes a",
	args: numericLessArgs{a: "00", b: "0"},
	want: numericLessWant{less: false},
}, {
	name: "zeroes b",
	args: numericLessArgs{a: "0", b: "00"},
	want: numericLessWant{less: true},
}, {
	name: "zeroes",
	args: numericLessArgs{a: "000", b: "000"},
	want: numericLessWant{less: false},
}, {
	name: "zeroes and letter a",
	args: numericLessArgs{a: "000a", b: "000"},
	want: numericLessWant{less: false},
}, {
	name: "zeroes and letter b",
	args: numericLessArgs{a: "000", b: "000b"},
	want: numericLessWant{less: true},
}, {
	name: "zeroes and letters 1",
	args: numericLessArgs{a: "000a", b: "000b"},
	want: numericLessWant{less: true},
}, {
	name: "zeroes and letters 2",
	args: numericLessArgs{a: "00v", b: "000v"},
	want: numericLessWant{less: true},
}, {
	name: "zeroes and letters 3",
	args: numericLessArgs{a: "000v", b: "00v"},
	want: numericLessWant{less: false},
}, {
	name: "zeroes digit - zeroes letter",
	args: numericLessArgs{a: "0001", b: "000b"},
	want: numericLessWant{less: false},
}, {
	name: "zeroes letter - zeroes digit",
	args: numericLessArgs{a: "000a", b: "0001"},
	want: numericLessWant{less: true},
}, {
	name: "numbers and letter 1",
	args: numericLessArgs{a: "001v", b: "002v"},
	want: numericLessWant{less: true},
}, {
	name: "numbers and letter 2",
	args: numericLessArgs{a: "001v", b: "001v"},
	want: numericLessWant{less: false},
}, {
	name: "numbers and letter 3",
	args: numericLessArgs{a: "002v", b: "001v"},
	want: numericLessWant{less: false},
}, {
	name: "numbers and letter 4",
	args: numericLessArgs{a: "001a", b: "001b"},
	want: numericLessWant{less: true},
}, {
	name: "numbers and letter 5",
	args: numericLessArgs{a: "001b", b: "001a"},
	want: numericLessWant{less: false},
}, {
	name: "numbers and letter 6",
	args: numericLessArgs{a: "001", b: "001b"},
	want: numericLessWant{less: true},
}, {
	name: "numbers and letter 7",
	args: numericLessArgs{a: "001a", b: "001"},
	want: numericLessWant{less: false},
}, {
	name: "long numbers 1",
	args: numericLessArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
	},
	want: numericLessWant{less: false},
}, {
	name: "long numbers 2",
	args: numericLessArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "888",
	},
	want: numericLessWant{less: false},
}, {
	name: "long numbers 3",
	args: numericLessArgs{
		a: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "888",
		b: strings.Repeat("0", 1000) + strings.Repeat("1", 1000) + "999",
	},
	want: numericLessWant{less: true},
}, {
	name: "long numbers 1a",
	args: numericLessArgs{
		a: strings.Repeat("1", 1000) + "88",
		b: strings.Repeat("1", 1000) + "999",
	},
	want: numericLessWant{less: true},
}, {
	name: "long numbers 2a",
	args: numericLessArgs{
		a: strings.Repeat("1", 1000) + "99",
		b: strings.Repeat("1", 1000) + "888",
	},
	want: numericLessWant{less: true},
}, {
	name: "long numbers 3a",
	args: numericLessArgs{
		a: strings.Repeat("1", 1000) + "888",
		b: strings.Repeat("1", 1000) + "99",
	},
	want: numericLessWant{less: false},
}, {
	name: "number groups 1",
	args: numericLessArgs{
		a: "2021.1.0",
		b: "2020.8.8",
	},
	want: numericLessWant{less: false},
}, {
	name: "number groups 2",
	args: numericLessArgs{
		a: "2021.1.1",
		b: "2021.1.2",
	},
	want: numericLessWant{less: true},
}, {
	name: "number groups 3",
	args: numericLessArgs{
		a: "2021.1.0",
		b: "2021.1.1",
	},
	want: numericLessWant{less: true},
}, {
	name: "number groups 4",
	args: numericLessArgs{
		a: "2020.12.1",
		b: "2020.12.1",
	},
	want: numericLessWant{less: false},
}, {
	name: "custom test 1",
	args: numericLessArgs{
		a: "Terminator 2",
		b: "Terminator Salvation",
	},
	want: numericLessWant{less: true},
}}

func BenchmarkNumericLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range numericLessTests {
			if NumericLess(tt.args.a, tt.args.b) != tt.want.less {
				require.Equal(b, tt.want.less, NumericLess(tt.args.a, tt.args.b), tt.name)
			}
		}
	}
}

func TestNumericLess(t *testing.T) {
	test := func(a numericLessArgs, w numericLessWant) func(t *testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, w.less, NumericLess(a.a, a.b), "invalid less")
		}
	}

	for _, tt := range numericLessTests {
		t.Run(tt.name, test(tt.args, tt.want))
	}
}
