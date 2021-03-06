package maths

// Abs gets the absolute value of an int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Pmod is a mositive modulo
// golang's % is a "remainder" function si misbehaves for negative modulus inputs
func Pmod(x, d int) int {
	if x == 0 || d == 0 {
		return 0
	}
	x = x % d
	if x >= 0 {
		return x
	} else if d < 0 {
		return x - d
	} else {
		return x + d
	}
}

// Max returns the highest int
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the lowest int
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// RoundUp rounds a value up to the nearest multiple
func RoundUp(toRound int, multiple int) int {
	remainder := Pmod(toRound, multiple)
	if remainder == 0 {
		return toRound
	}

	return (multiple - remainder) + toRound
}

// RoundDown rounds a value down to the nearest multiple
func RoundDown(toRound int, multiple int) int {
	remainder := Pmod(toRound, multiple)
	return toRound - remainder
}
