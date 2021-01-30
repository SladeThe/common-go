package strings

import (
	"unicode/utf8"
)

// NumericLess compares strings with respect to values of positive integer groups.
// For example, 'a9z' is considered less than 'a11z', because 9 < 11.
// If two numbers with leading zeroes have the same value, the shortest of them is considered less, i.e. 12 < 012.
// Digits and non-digits are compared lexicographically, i.e. ' ' (space) < 5 < 'a'.
func NumericLess(a, b string) bool {
	nextA := func() (rune, int) {
		r, size := utf8.DecodeRuneInString(a)
		a = a[size:]
		return r, size
	}

	nextB := func() (rune, int) {
		r, size := utf8.DecodeRuneInString(b)
		b = b[size:]
		return r, size
	}

	for {
		runeA, offsetA := nextA()
		if offsetA == 0 {
			return b != ""
		}

		runeB, offsetB := nextB()
		if offsetB == 0 {
			return false
		}

		if digitA, digitB := isDigit(runeA), isDigit(runeB); digitA != digitB {
			return runeA < runeB
		} else if digitA {
			zeroBalance := 0
			digitCmp := 0

			for runeA == '0' {
				zeroBalance++
				runeA, offsetA = nextA()
			}

			for runeB == '0' {
				zeroBalance--
				runeB, offsetB = nextB()
			}

			if offsetA == 0 {
				return offsetB != 0 || zeroBalance < 0
			}

			if offsetB == 0 {
				return false
			}

			if digitA, digitB = isDigit(runeA), isDigit(runeB); !digitA && !digitB {
				if zeroBalance != 0 {
					return zeroBalance < 0
				} else if runeA != runeB {
					return runeA < runeB
				}
			} else if digitA != digitB {
				return digitB
			} else {
				for {
					if digitCmp == 0 && runeA != runeB {
						if runeA < runeB {
							digitCmp = -1
						} else {
							digitCmp = 1
						}
					}

					runeA, offsetA = nextA()
					runeB, offsetB = nextB()

					if digitA, digitB = isDigit(runeA), isDigit(runeB); digitA != digitB {
						return digitB
					} else if !digitA {
						if digitCmp != 0 {
							return digitCmp < 0
						}

						if zeroBalance != 0 {
							return zeroBalance < 0
						}

						if offsetA == 0 {
							return offsetB != 0
						}

						if offsetB == 0 {
							return false
						}

						if runeA != runeB {
							return runeA < runeB
						}

						break
					}
				}
			}
		} else if runeA != runeB {
			return runeA < runeB
		}
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
