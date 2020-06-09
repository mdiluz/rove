package maths

// Abs gets the absolute value of an int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// pmod is a mositive modulo
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
func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the lowest int
func Min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}
